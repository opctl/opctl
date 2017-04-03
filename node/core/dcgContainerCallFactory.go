package core

import (
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	interpolatePkg "github.com/opspec-io/sdk-golang/interpolate"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/virtual-go/dircopier"
	"github.com/virtual-go/filecopier"
	"github.com/virtual-go/fs/osfs"
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
) (dcgContainerCall *model.DCGContainerCall, err error) {
	fs := osfs.New()
	fileCopier := filecopier.New()
	dirCopier := dircopier.New()
	interpolate := interpolatePkg.New()

	dcgContainerCall = &model.DCGContainerCall{
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

	// create scratch dir for container
	scratchDirPath := path.Join(
		appdatapath.New().PerUser(),
		"opctl",
		"dcg",
		rootOpId,
		"containers",
		containerId,
		"fs",
	)
	err = fs.MkdirAll(scratchDirPath, 0700)
	if nil != err {
		return
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
			err = dirCopier.Fs(
				path.Join(pkgRef, scgContainerDirBind),
				dcgContainerCall.Dirs[scgContainerDirPath],
			)
		} else {
			// is bound to variable
			if varData, ok := currentScope[scgContainerDirBind]; ok {
				// bound to input
				dcgContainerCall.Dirs[scgContainerDirPath] = varData.Dir
			} else {
				// bound to output
				// create placeholder dir on host so the output points to something
				dcgContainerCall.Dirs[scgContainerDirPath] = path.Join(scratchDirPath, scgContainerDirPath)
				err = fs.MkdirAll(dcgContainerCall.Dirs[scgContainerDirPath], 0700)
			}
		}
		if nil != err {
			return
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
			err = fileCopier.Fs(
				path.Join(pkgRef, scgContainerFileBind),
				dcgContainerCall.Files[scgContainerFilePath],
			)
		} else {
			// is bound to variable
			if varData, ok := currentScope[scgContainerFileBind]; ok {
				// bound to input
				dcgContainerCall.Files[scgContainerFilePath] = varData.File
			} else {
				// bound to output
				// create outputFile on host so the output points to something
				dcgContainerCall.Files[scgContainerFilePath] = path.Join(scratchDirPath, scgContainerFilePath)
				// create dir
				err = fs.MkdirAll(path.Dir(dcgContainerCall.Files[scgContainerFilePath]), 0700)
				if nil != err {
					return
				}
				// create file
				var outputFile *os.File
				outputFile, err = fs.Create(dcgContainerCall.Files[scgContainerFilePath])
				outputFile.Close()
			}
		}
		if nil != err {
			return
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
			dcgContainerCall.Sockets[scgContainerSocketAddress] = boundArg.Socket
		} else if isUnixSocketAddress(scgContainerSocketAddress) {
			// bound to output
			// create outputSocket on host so the output points to something
			if isUnixSocketAddress(scgContainerSocketAddress) {
				dcgHostSocketAddress := path.Join(scratchDirPath, scgContainerSocketAddress)
				var outputSocket *os.File
				outputSocket, err = fs.Create(dcgHostSocketAddress)
				outputSocket.Close()
				if nil != err {
					return
				}
				err = fs.Chmod(dcgHostSocketAddress, os.ModeSocket)
				if nil != err {
					return
				}
				dcgContainerCall.Sockets[scgContainerSocketAddress] = dcgHostSocketAddress
			}
		}
	}

	return

}
