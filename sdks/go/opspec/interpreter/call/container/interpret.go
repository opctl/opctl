package container

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/runcontext"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/caches"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/cmd"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/dirs"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/envvars"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/files"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/image"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/sockets"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// Interpret a container
func Interpret(
	ctx context.Context,
	scope map[string]*model.Value,
	containerCallSpec *model.ContainerCallSpec,
	containerID string,
	opPath string,
	dataDirPath string,
) (*model.ContainerCall, error) {

	containerCall := &model.ContainerCall{
		BaseCall: model.BaseCall{
			OpPath: opPath,
		},
		Dirs:        map[string]string{},
		EnvVars:     map[string]string{},
		Files:       map[string]string{},
		Sockets:     map[string]string{},
		WorkDir:     containerCallSpec.WorkDir,
		ContainerID: containerID,
		Ports:       containerCallSpec.Ports,
	}

	// construct dcg container path
	scratchDirPath := filepath.Join(
		dataDirPath,
		"dcg",
		containerID,
		"fs",
	)

	if err := unsudo.CreateDir(scratchDirPath); err != nil {
		return nil, err
	}

	// interpret cmd
	var err error
	containerCall.Cmd, err = cmd.Interpret(
		scope,
		containerCallSpec.Cmd,
	)
	if err != nil {
		return nil, err
	}

	dataCachePath := filepath.Join(dataDirPath, "ops")

	// interpret dirs
	containerCall.Dirs, err = dirs.Interpret(
		scope,
		containerCallSpec.Dirs,
		scratchDirPath,
		dataCachePath,
	)
	if err != nil {
		return nil, err
	}

	// interpret envVars
	containerCall.EnvVars, err = envvars.Interpret(
		scope,
		containerCallSpec.EnvVars,
	)
	if err != nil {
		return nil, err
	}

	// interpret files
	containerCall.Files, err = files.Interpret(
		scope,
		containerCallSpec.Files,
		scratchDirPath,
		dataCachePath,
	)
	if err != nil {
		return nil, err
	}

	// interpret image
	containerCall.Image, err = image.Interpret(
		scope,
		containerCallSpec.Image,
		scratchDirPath,
	)
	if err != nil {
		return nil, err
	}

	// interpret sockets
	containerCall.Sockets, err = sockets.Interpret(
		scope,
		containerCallSpec.Sockets,
		scratchDirPath,
	)
	if err != nil {
		return nil, err
	}

	// interpret caches
	cacheDirs, err := caches.Interpret(
		scope,
		containerCallSpec.Caches,
		containerCall.Image.Ref,
		dataDirPath,
		runcontext.NoCacheRoot(ctx),
	)
	if err != nil {
		return nil, err
	}
	declaredMounts := map[string]struct{}{}
	for containerPath := range containerCall.Dirs {
		declaredMounts[path.Clean(containerPath)] = struct{}{}
	}
	for containerPath := range containerCall.Files {
		declaredMounts[path.Clean(containerPath)] = struct{}{}
	}
	for containerPath, hostAddress := range containerCall.Sockets {
		// Match the unix socket bind mounts in docker.constructHostConfig.
		if strings.ContainsAny(hostAddress, `/\`) {
			declaredMounts[path.Clean(containerPath)] = struct{}{}
		}
	}
	for containerPath, hostPath := range cacheDirs {
		if _, exists := declaredMounts[containerPath]; exists {
			return nil, fmt.Errorf("unable to bind cache dir %q: container path is also declared in dirs, files or sockets", containerPath)
		}
		containerCall.Dirs[containerPath] = hostPath
	}

	// interpret name as string
	if containerCallSpec.Name != nil {
		containerCallName, err := str.Interpret(
			scope,
			*containerCallSpec.Name,
		)
		if err != nil {
			return nil, err
		}
		containerCall.Name = containerCallName.String
	}

	// interpret workDir
	if containerCallSpec.WorkDir != "" {
		containerCallWorkDir, err := str.Interpret(
			scope,
			containerCallSpec.WorkDir,
		)
		if err != nil {
			return nil, err
		}

		containerCall.WorkDir = *containerCallWorkDir.String
	}

	return containerCall, nil

}
