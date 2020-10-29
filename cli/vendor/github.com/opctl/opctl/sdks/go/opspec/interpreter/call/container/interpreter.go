package container

import (
	"path/filepath"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/cmd"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/dirs"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/envvars"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/files"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/image"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/sockets"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	// Interpret interprets an ContainerCallSpec into a ContainerCall
	Interpret(
		scope map[string]*model.Value,
		containerCallSpec *model.ContainerCallSpec,
		containerID string,
		rootCallID string,
		opPath string,
	) (*model.ContainerCall, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter(
	dataDirPath string,
) Interpreter {
	return _interpreter{
		cmdInterpreter:     cmd.NewInterpreter(),
		dirsInterpreter:    dirs.NewInterpreter(dataDirPath),
		envVarsInterpreter: envvars.NewInterpreter(),
		filesInterpreter:   files.NewInterpreter(dataDirPath),
		imageInterpreter:   image.NewInterpreter(),
		os:                 ios.New(),
		dataDirPath:        dataDirPath,
		stringInterpreter:  str.NewInterpreter(),
		socketsInterpreter: sockets.NewInterpreter(),
	}
}

type _interpreter struct {
	cmdInterpreter     cmd.Interpreter
	dirsInterpreter    dirs.Interpreter
	envVarsInterpreter envvars.Interpreter
	filesInterpreter   files.Interpreter
	imageInterpreter   image.Interpreter
	os                 ios.IOS
	dataDirPath        string
	stringInterpreter  str.Interpreter
	socketsInterpreter sockets.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	containerCallSpec *model.ContainerCallSpec,
	containerID string,
	rootCallID string,
	opPath string,
) (*model.ContainerCall, error) {

	containerCall := &model.ContainerCall{
		BaseCall: model.BaseCall{
			RootCallID: rootCallID,
			OpPath:     opPath,
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
		itp.dataDirPath,
		"dcg",
		rootCallID,
		"containers",
		containerID,
		"fs",
	)
	if err := itp.os.MkdirAll(scratchDirPath, 0700); nil != err {
		return nil, err
	}

	// interpret cmd
	var err error
	containerCall.Cmd, err = itp.cmdInterpreter.Interpret(scope, containerCallSpec.Cmd)
	if nil != err {
		return nil, err
	}

	// interpret dirs
	containerCall.Dirs, err = itp.dirsInterpreter.Interpret(scope, containerCallSpec.Dirs, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret envVars
	containerCall.EnvVars, err = itp.envVarsInterpreter.Interpret(scope, containerCallSpec.EnvVars)
	if nil != err {
		return nil, err
	}

	// interpret files
	containerCall.Files, err = itp.filesInterpreter.Interpret(scope, containerCallSpec.Files, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret image
	containerCall.Image, err = itp.imageInterpreter.Interpret(scope, containerCallSpec.Image, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret name as string
	if nil != containerCallSpec.Name {
		containerCallName, err := itp.stringInterpreter.Interpret(scope, *containerCallSpec.Name)
		if nil != err {
			return nil, err
		}
		containerCall.Name = containerCallName.String
	}

	// interpret workDir
	if "" != containerCallSpec.WorkDir {
		containerCallWorkDir, err := itp.stringInterpreter.Interpret(scope, containerCallSpec.WorkDir)
		if nil != err {
			return nil, err
		}

		containerCall.WorkDir = *containerCallWorkDir.String
	}

	// interpret sockets
	containerCall.Sockets, err = itp.socketsInterpreter.Interpret(scope, containerCallSpec.Sockets, scratchDirPath)

	return containerCall, err

}
