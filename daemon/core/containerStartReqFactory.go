package core

import (
	appdatapath "github.com/appdataspec/sdk-golang/pkg/path"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/sdk-golang/pkg/interpolate"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"path/filepath"
)

func newContainerStartReq(
	currentScope map[string]*model.Data,
	scgContainerCall *model.ScgContainerCall,
	containerId string,
	inputs []*model.Param,
	opGraphId string,
) *containerengine.StartContainerReq {

	// create new slice so we don't mutate containerCall.Cmd @ caller
	cmd := append([]string{}, scgContainerCall.Cmd...)
	dirs := map[string]string{}
	envVars := map[string]string{}
	files := map[string]string{}
	sockets := map[string]string{}

	// create scratch dir for container
	scratchDirPath := filepath.Join(
		appdatapath.New().PerUser(),
		"opctl",
		"dcg",
		opGraphId,
		"containers",
		containerId,
		"fs",
	)
	err := os.MkdirAll(scratchDirPath, 0700)
	if nil != err {
		panic(err)
	}

	// construct envVars
	for scgContainerEnvVarName, scgContainerEnvVar := range scgContainerCall.EnvVars {
		envVars[scgContainerEnvVarName] = scgContainerEnvVar.Value
	}

	// construct files
	for scgContainerFilePath, scgContainerFile := range scgContainerCall.Files {
		if boundArg, ok := currentScope[scgContainerFile.Bind]; ok {
			// bound to input
			files[scgContainerFilePath] = boundArg.File
		} else {
			// bound to output
			dcgHostFilePath := filepath.Join(scratchDirPath, scgContainerFilePath)
			_, err = os.Create(dcgHostFilePath)
			if nil != err {
				panic(err)
			}
			files[scgContainerFilePath] = dcgHostFilePath
		}
	}

	// construct dirs
	for scgContainerDirPath, scgContainerDir := range scgContainerCall.Dirs {
		if boundArg, ok := currentScope[scgContainerDir.Bind]; ok {
			// bound to input
			dirs[scgContainerDirPath] = boundArg.Dir
		} else {
			// bound to output
			dcgHostDirPath := filepath.Join(scratchDirPath, scgContainerDirPath)
			err := os.MkdirAll(dcgHostDirPath, 0700)
			if nil != err {
				panic(err)
			}
			files[scgContainerDirPath] = dcgHostDirPath
		}
	}

	// construct sockets
	for scgContainerSocketAddress, scgContainerSocket := range scgContainerCall.Sockets {
		if boundArg, ok := currentScope[scgContainerSocket.Bind]; ok {
			// bound to input
			switch {
			case isUnixSocketAddress(scgContainerSocketAddress):
				sockets[scgContainerSocketAddress] = boundArg.Socket
			default:
				// @TODO: handle network sockets
			}
		} else {
			// bound to output
			switch {
			case isUnixSocketAddress(scgContainerSocketAddress):
				dcgHostSocketAddress := filepath.Join(scratchDirPath, scgContainerSocketAddress)
				_, err = os.Create(dcgHostSocketAddress)
				if nil != err {
					panic(err)
				}
				err = os.Chmod(dcgHostSocketAddress, os.ModeSocket)
				if nil != err {
					panic(err)
				}
				sockets[scgContainerSocketAddress] = dcgHostSocketAddress
			default:
				// @TODO: handle network sockets
			}
		}
	}

	for _, input := range inputs {
		switch {
		case nil != input.String:
			stringInput := input.String

			// obtain inputValue
			inputValue := ""
			if _, isArgForInput := currentScope[stringInput.Name]; isArgForInput {
				// use provided arg for param
				inputValue = currentScope[stringInput.Name].String
			} else {
				// no provided arg for param; fallback to default
				inputValue = stringInput.Default
			}

			// interpolate interpolatedStrings w/ inputValue
			for cmdEntryIndex, cmdEntry := range cmd {
				cmd[cmdEntryIndex] = interpolate.Interpolate(cmdEntry, stringInput.Name, inputValue)
			}
			for containerEnvVarName, containerEnvVar := range envVars {
				envVars[containerEnvVarName] = interpolate.Interpolate(containerEnvVar, stringInput.Name, inputValue)
			}
		}
	}
	return &containerengine.StartContainerReq{
		Cmd:         cmd,
		Dirs:        dirs,
		Env:         envVars,
		Files:       files,
		Image:       scgContainerCall.Image,
		Sockets:     sockets,
		WorkDir:     scgContainerCall.WorkDir,
		ContainerId: containerId,
		OpGraphId:   opGraphId,
	}

}
