package containercall

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func (df _ContainerCall) Interpret(
	scope map[string]*model.Value,
	scgContainerCall *model.SCGContainerCall,
	containerId string,
	rootOpId string,
	pkgPath string,
) (*model.DCGContainerCall, error) {

	dcgContainerCall := &model.DCGContainerCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId: rootOpId,
			PkgRef:   pkgPath,
		},
		Dirs:        map[string]string{},
		EnvVars:     map[string]string{},
		Files:       map[string]string{},
		Sockets:     map[string]string{},
		WorkDir:     scgContainerCall.WorkDir,
		ContainerId: containerId,
		Name:        scgContainerCall.Name,
		Ports:       scgContainerCall.Ports,
	}

	// construct dcg container path
	scratchDirPath := filepath.Join(
		df.rootFSPath,
		"dcg",
		rootOpId,
		"containers",
		containerId,
		"fs",
	)
	if err := df.os.MkdirAll(scratchDirPath, 0700); nil != err {
		return nil, err
	}

	// construct cmd
	for _, cmdEntry := range scgContainerCall.Cmd {
		// interpolate each entry
		dcgContainerCall.Cmd = append(dcgContainerCall.Cmd, df.interpolater.Interpolate(cmdEntry, scope))
	}

	// construct dirs
	for scgContainerDirPath, scgContainerDirBind := range scgContainerCall.Dirs {

		if "" == scgContainerDirBind {
			// bound implicitly
			scgContainerDirBind = scgContainerDirPath
		}

		isBoundToPkg := strings.HasPrefix(scgContainerDirBind, "/")
		value, isBoundToScope := scope[scgContainerDirBind]

		switch {
		case isBoundToPkg:
			// bound to pkg dir
			dcgContainerCall.Dirs[scgContainerDirPath] = filepath.Join(scratchDirPath, scgContainerDirBind)

			// pkg dirs must be passed by value
			if err := df.dirCopier.OS(
				filepath.Join(pkgPath, scgContainerDirBind),
				dcgContainerCall.Dirs[scgContainerDirPath],
			); nil != err {
				return nil, err
			}
		case isBoundToScope:
			// bound to scope
			if nil == value || nil == value.Dir {
				return nil, fmt.Errorf(
					"Unable to bind dir '%v' to '%v'. '%v' not a dir",
					scgContainerDirPath,
					scgContainerDirBind,
					scgContainerDirBind,
				)
			}

			if strings.HasPrefix(*value.Dir, df.rootFSPath) {
				// bound to rootFS dir
				dcgContainerCall.Dirs[scgContainerDirPath] = filepath.Join(scratchDirPath, scgContainerDirPath)

				// rootFS dirs must be passed by value
				if err := df.dirCopier.OS(
					*value.Dir,
					dcgContainerCall.Dirs[scgContainerDirPath],
				); nil != err {
					return nil, err
				}
			} else {
				// bound to non rootFS dir
				dcgContainerCall.Dirs[scgContainerDirPath] = *value.Dir
			}
		default:
			// unbound; create tree
			dcgContainerCall.Dirs[scgContainerDirPath] = filepath.Join(scratchDirPath, scgContainerDirPath)
			if err := df.os.MkdirAll(
				dcgContainerCall.Dirs[scgContainerDirPath],
				0700,
			); nil != err {
				return nil, err
			}
		}
	}

	// construct envVars
	for envVarName, scgContainerEnvVar := range scgContainerCall.EnvVars {
		if "" == scgContainerEnvVar {
			// implicitly bound
			value, ok := scope[envVarName]
			if !ok {
				return nil, fmt.Errorf("Unable to bind env var to '%v' via implicit ref. '%v' is not in scope", envVarName, envVarName)
			}

			switch {
			case nil != value.String:
				dcgContainerCall.EnvVars[envVarName] = *value.String
			case nil != value.Number:
				dcgContainerCall.EnvVars[envVarName] = strconv.FormatFloat(*value.Number, 'f', -1, 64)
			}
			continue
		}

		// otherwise interpolate value
		dcgContainerCall.EnvVars[envVarName] = df.interpolater.Interpolate(scgContainerEnvVar, scope)

	}

	// construct files
	for scgContainerFilePath, scgContainerFileBind := range scgContainerCall.Files {

		if "" == scgContainerFileBind {
			// bound implicitly
			scgContainerFileBind = scgContainerFilePath
		}

		isBoundToPkg := strings.HasPrefix(scgContainerFileBind, "/")
		value, isBoundToScope := scope[scgContainerFileBind]

		switch {
		case isBoundToPkg:
			// bound to pkg file
			dcgContainerCall.Files[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFileBind)

			// pkg files must be passed by value
			if err := df.fileCopier.OS(
				filepath.Join(pkgPath, scgContainerFileBind),
				dcgContainerCall.Files[scgContainerFilePath],
			); nil != err {
				return nil, err
			}
		case isBoundToScope:
			// bound to scope
			if nil == value || nil == value.File {
				return nil, fmt.Errorf(
					"Unable to bind file '%v' to '%v'. '%v' not a file",
					scgContainerFilePath,
					scgContainerFileBind,
					scgContainerFileBind,
				)
			}

			if strings.HasPrefix(*value.File, df.rootFSPath) {
        // bound to rootFS file
        dcgContainerCall.Files[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)

				// rootFS files must be passed by value
				if err := df.fileCopier.OS(
					*value.File,
					dcgContainerCall.Files[scgContainerFilePath],
				); nil != err {
					return nil, err
				}
			} else {
				// bound to non rootFS file
				dcgContainerCall.Files[scgContainerFilePath] = *value.File
			}
		default:
			// unbound; create tree
			dcgContainerCall.Files[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)
			// create dir
			if err := df.os.MkdirAll(
				path.Dir(dcgContainerCall.Files[scgContainerFilePath]),
				0700,
			); nil != err {
				return nil, err
			}
			// create file
			var outputFile *os.File
			outputFile, err := df.os.Create(dcgContainerCall.Files[scgContainerFilePath])
			outputFile.Close()
			if nil != err {
				return nil, err
			}
		}
	}

	// construct image
	if scgContainerCallImage := scgContainerCall.Image; scgContainerCallImage != nil {
		dcgContainerCall.Image = &model.DCGContainerCallImage{
			Ref: df.interpolater.Interpolate(scgContainerCall.Image.Ref, scope),
		}
		if "" != scgContainerCallImage.PullIdentity && "" != scgContainerCallImage.PullSecret {
			// fallback for deprecated cred format
			scgContainerCallImage.PullCreds = &model.SCGPullCreds{
				Username: scgContainerCallImage.PullIdentity,
				Password: scgContainerCallImage.PullSecret,
			}
		}

		if nil != scgContainerCallImage.PullCreds {
			dcgContainerCall.Image.PullCreds = &model.DCGPullCreds{
				Username: df.interpolater.Interpolate(scgContainerCall.Image.PullCreds.Username, scope),
				Password: df.interpolater.Interpolate(scgContainerCall.Image.PullCreds.Password, scope),
			}
		}
	}

	// construct sockets
	for scgContainerSocketAddress, scgContainerSocketBind := range scgContainerCall.Sockets {
		if boundArg, ok := scope[scgContainerSocketBind]; ok {
			// bound to var
			dcgContainerCall.Sockets[scgContainerSocketAddress] = *boundArg.Socket
		} else if isUnixSocketAddress(scgContainerSocketAddress) {
			// bound to output
			// create outputSocket on host so the output points to something
			dcgHostSocketAddress := filepath.Join(scratchDirPath, scgContainerSocketAddress)
			var outputSocket *os.File
			outputSocket, err := df.os.Create(dcgHostSocketAddress)
			outputSocket.Close()
			if nil != err {
				return nil, err
			}
			if err := os.Chmod(
				dcgHostSocketAddress,
				os.ModeSocket,
			); nil != err {
				return nil, err
			}
			dcgContainerCall.Sockets[scgContainerSocketAddress] = dcgHostSocketAddress
		}
	}

	return dcgContainerCall, nil

}
