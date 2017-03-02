package core

import (
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/dircopier"
	"github.com/opspec-io/opctl/util/filecopier"
	osfs "github.com/opspec-io/opctl/util/vfs/os"
	"github.com/opspec-io/sdk-golang/pkg/interpolate"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"path"
	"strings"
)

func newRunContainerReq(
	currentScope map[string]*model.Data,
	scgContainerCall *model.ScgContainerCall,
	containerId string,
	opGraphId string,
	opRef string,
) (req *containerprovider.RunContainerReq, err error) {
	fs := osfs.New()
	fileCopier := filecopier.New()
	dirCopier := dircopier.New()

	req = &containerprovider.RunContainerReq{
		Cmd:         append([]string{}, scgContainerCall.Cmd...), // create new slice so we don't cause side effects
		Dirs:        map[string]string{},
		Env:         map[string]string{},
		Files:       map[string]string{},
		Image:       scgContainerCall.Image,
		Sockets:     map[string]string{},
		WorkDir:     scgContainerCall.WorkDir,
		ContainerId: containerId,
		OpGraphId:   opGraphId,
	}

	// create scratch dir for container
	scratchDirPath := path.Join(
		appdatapath.New().PerUser(),
		"opctl",
		"dcg",
		opGraphId,
		"containers",
		containerId,
		"fs",
	)
	err = fs.MkdirAll(scratchDirPath, 0700)
	if nil != err {
		return
	}

	// construct dirs
	for scgContainerDirPath, scgContainerDirBind := range scgContainerCall.Dirs {
		if "" == scgContainerDirBind {
			// use container dir path as bundle dir path if not provided
			scgContainerDirBind = scgContainerDirPath
		}

		if strings.HasPrefix(scgContainerDirBind, "/") {
			// is bound to bundle path
			req.Dirs[scgContainerDirPath] = path.Join(scratchDirPath, scgContainerDirBind)
			err = dirCopier.Fs(
				path.Join(opRef, scgContainerDirBind),
				req.Dirs[scgContainerDirPath],
			)
		} else {
			// is bound to variable
			if boundArg, ok := currentScope[scgContainerDirBind]; ok {
				// bound to input
				req.Dirs[scgContainerDirPath] = boundArg.Dir
			} else {
				// bound to output
				// create placeholder dir on host so the output points to something
				req.Files[scgContainerDirPath] = path.Join(scratchDirPath, scgContainerDirPath)
				err = fs.MkdirAll(req.Files[scgContainerDirPath], 0700)
			}
		}
		if nil != err {
			return
		}
	}

	// construct envVars
	for scgContainerEnvVarName, scgContainerEnvVar := range scgContainerCall.EnvVars {
		req.Env[scgContainerEnvVarName] = scgContainerEnvVar
	}

	// construct files
	for scgContainerFilePath, scgContainerFileBind := range scgContainerCall.Files {
		if "" == scgContainerFileBind {
			// use container file path as bundle file path if not provided
			scgContainerFileBind = scgContainerFilePath
		}

		if strings.HasPrefix(scgContainerFileBind, "/") {
			// is bound to bundle path
			req.Files[scgContainerFilePath] = path.Join(scratchDirPath, scgContainerFileBind)
			err = fileCopier.Fs(
				path.Join(opRef, scgContainerFileBind),
				req.Files[scgContainerFilePath],
			)
		} else {
			// is bound to variable
			if boundArg, ok := currentScope[scgContainerFileBind]; ok {
				// bound to input
				req.Files[scgContainerFilePath] = boundArg.File
			} else {
				// bound to output
				// create outputFile on host so the output points to something
				req.Files[scgContainerFilePath] = path.Join(scratchDirPath, scgContainerFilePath)
				// create dir
				err = fs.MkdirAll(path.Dir(req.Files[scgContainerFilePath]), 0700)
				if nil != err {
					return
				}
				// create file
				var outputFile *os.File
				outputFile, err = fs.Create(req.Files[scgContainerFilePath])
				outputFile.Close()
			}
		}
		if nil != err {
			return
		}
	}

	// construct sockets
	for scgContainerSocketAddress, scgContainerSocketBind := range scgContainerCall.Sockets {
		if boundArg, ok := currentScope[scgContainerSocketBind]; ok {
			// bound to input
			req.Sockets[scgContainerSocketAddress] = boundArg.Socket
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
				req.Sockets[scgContainerSocketAddress] = dcgHostSocketAddress
			}
		}
	}

	// interpolate cmd & envVars w/ values from currentScope
	for varName, varData := range currentScope {
		switch {
		case 0 != varData.Number:
			numberVarData := varData.Number

			for cmdEntryIndex, cmdEntry := range req.Cmd {
				req.Cmd[cmdEntryIndex] = interpolate.NumberValue(cmdEntry, varName, numberVarData)
			}
			for containerEnvVarName, containerEnvVar := range req.Env {
				req.Env[containerEnvVarName] = interpolate.NumberValue(containerEnvVar, varName, numberVarData)
			}
		case "" != varData.String:
			stringVarData := varData.String

			for cmdEntryIndex, cmdEntry := range req.Cmd {
				req.Cmd[cmdEntryIndex] = interpolate.StringValue(cmdEntry, varName, stringVarData)
			}
			for containerEnvVarName, containerEnvVar := range req.Env {
				req.Env[containerEnvVarName] = interpolate.StringValue(containerEnvVar, varName, stringVarData)
			}
		}
	}

	return

}
