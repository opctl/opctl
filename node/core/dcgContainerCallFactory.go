package core

import (
	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-utils/dircopier"
	"github.com/golang-utils/filecopier"
	interpolatePkg "github.com/opspec-io/sdk-golang/interpolate"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path"
	"strings"
)

func constructDCGContainerCall(
	currentScope map[string]*model.Data,
	scgContainerCall *model.SCGContainerCall,
	containerId string,
	rootOpId string,
	pkgRef string,
) (*model.DCGContainerCall, error) {
	fileCopier := filecopier.New()
	dirCopier := dircopier.New()
	interpolate := interpolatePkg.New()

	dcgContainerCall := &model.DCGContainerCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId: rootOpId,
			PkgRef:   pkgRef,
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

	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		return nil, err
	}

	// create scratch dir for container
	scratchDirPath := path.Join(
		perUserAppDataPath,
		"opctl",
		"dcg",
		rootOpId,
		"containers",
		containerId,
		"fs",
	)
	err = os.MkdirAll(scratchDirPath, 0700)
	if nil != err {
		return nil, err
	}

	// construct cmd
	for _, cmdEntry := range scgContainerCall.Cmd {
		// interpolate each entry
		dcgContainerCall.Cmd = append(dcgContainerCall.Cmd, interpolate.Interpolate(cmdEntry, currentScope))
	}

	// construct dirs
	for scgContainerDirPath, scgContainerDirBind := range scgContainerCall.Dirs {
		if "" == scgContainerDirBind {
			// use container dir path as pkg dir path if not provided
			scgContainerDirBind = scgContainerDirPath
		}

		if strings.HasPrefix(scgContainerDirBind, "/") {
			// is bound to pkg path
			dcgContainerCall.Dirs[scgContainerDirPath] = path.Join(scratchDirPath, scgContainerDirBind)
			err = dirCopier.OS(
				path.Join(pkgRef, scgContainerDirBind),
				dcgContainerCall.Dirs[scgContainerDirPath],
			)
		} else {
			// is bound to variable
			if varData, ok := currentScope[scgContainerDirBind]; ok {
				// bound to input
				dcgContainerCall.Dirs[scgContainerDirPath] = *varData.Dir
			} else {
				// bound to output
				// create placeholder dir on host so the output points to something
				dcgContainerCall.Dirs[scgContainerDirPath] = path.Join(scratchDirPath, scgContainerDirPath)
				err = os.MkdirAll(dcgContainerCall.Dirs[scgContainerDirPath], 0700)
			}
		}
		if nil != err {
			return nil, err
		}
	}

	// construct envVars
	for scgContainerEnvVarName, scgContainerEnvVar := range scgContainerCall.EnvVars {
		// interpolate each value
		dcgContainerCall.EnvVars[scgContainerEnvVarName] = interpolate.Interpolate(scgContainerEnvVar, currentScope)
	}

	// construct files
	for scgContainerFilePath, scgContainerFileBind := range scgContainerCall.Files {
		if "" == scgContainerFileBind {
			// use container file path as pkg file path if not provided
			scgContainerFileBind = scgContainerFilePath
		}

		if strings.HasPrefix(scgContainerFileBind, "/") {
			// is bound to pkg path
			dcgContainerCall.Files[scgContainerFilePath] = path.Join(scratchDirPath, scgContainerFileBind)
			err = fileCopier.OS(
				path.Join(pkgRef, scgContainerFileBind),
				dcgContainerCall.Files[scgContainerFilePath],
			)
		} else {
			// is bound to variable
			if varData, ok := currentScope[scgContainerFileBind]; ok {
				// bound to input
				dcgContainerCall.Files[scgContainerFilePath] = *varData.File
			} else {
				// bound to output
				// create outputFile on host so the output points to something
				dcgContainerCall.Files[scgContainerFilePath] = path.Join(scratchDirPath, scgContainerFilePath)
				// create dir
				err = os.MkdirAll(path.Dir(dcgContainerCall.Files[scgContainerFilePath]), 0700)
				if nil != err {
					return nil, err
				}
				// create file
				var outputFile *os.File
				outputFile, err = os.Create(dcgContainerCall.Files[scgContainerFilePath])
				outputFile.Close()
			}
		}
		if nil != err {
			return nil, err
		}
	}

	// construct image
	if scgContainerCallImage := scgContainerCall.Image; scgContainerCallImage != nil {
		dcgContainerCall.Image = &model.DCGContainerCallImage{
			// interpolate all properties
			Ref: interpolate.Interpolate(scgContainerCall.Image.Ref, currentScope),
		}
		if nil != scgContainerCallImage.PullAuth {
			dcgContainerCall.Image.PullAuth = &model.DCGUsernamePasswordAuth{
				Username: interpolate.Interpolate(scgContainerCall.Image.PullAuth.Username, currentScope),
				Password: interpolate.Interpolate(scgContainerCall.Image.PullAuth.Password, currentScope),
			}
		}
	}

	// construct sockets
	for scgContainerSocketAddress, scgContainerSocketBind := range scgContainerCall.Sockets {
		if boundArg, ok := currentScope[scgContainerSocketBind]; ok {
			// bound to var
			dcgContainerCall.Sockets[scgContainerSocketAddress] = *boundArg.Socket
		} else if isUnixSocketAddress(scgContainerSocketAddress) {
			// bound to output
			// create outputSocket on host so the output points to something
			if isUnixSocketAddress(scgContainerSocketAddress) {
				dcgHostSocketAddress := path.Join(scratchDirPath, scgContainerSocketAddress)
				var outputSocket *os.File
				outputSocket, err = os.Create(dcgHostSocketAddress)
				outputSocket.Close()
				if nil != err {
					return nil, err
				}
				err = os.Chmod(dcgHostSocketAddress, os.ModeSocket)
				if nil != err {
					return nil, err
				}
				dcgContainerCall.Sockets[scgContainerSocketAddress] = dcgHostSocketAddress
			}
		}
	}

	return dcgContainerCall, nil

}
