package core

import (
	appdatapath "github.com/appdataspec/sdk-golang/pkg/path"
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/sdk-golang/pkg/interpolate"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"path/filepath"
)

func newRunContainerReq(
	currentScope map[string]*model.Data,
	scgContainerCall *model.ScgContainerCall,
	containerId string,
	inputs map[string]*model.Param,
	opGraphId string,
) *containerprovider.RunContainerReq {

	// create new slice so we don't cause side effects
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
		envVars[scgContainerEnvVarName] = scgContainerEnvVar
	}

	// construct files
	for scgContainerFilePath, scgContainerFile := range scgContainerCall.Files {
		if boundArg, ok := currentScope[scgContainerFile.Bind]; ok {
			// bound to input
			files[scgContainerFilePath] = boundArg.File
		} else {
			// bound to output
			// create placeholder file on host so the output points to something
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
			// create placeholder dir on host so the output points to something
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
			sockets[scgContainerSocketAddress] = boundArg.Socket
		} else if isUnixSocketAddress(scgContainerSocketAddress) {
			// bound to output
			// create placeholder unix socket on host so the output points to something
			if isUnixSocketAddress(scgContainerSocketAddress) {
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
			}
		}
	}

	for inputName, input := range inputs {
		switch {
		case nil != input.String:
			stringInput := input.String

			// obtain inputValue
			inputValue := ""
			if _, isArgForInput := currentScope[inputName]; isArgForInput {
				// use provided arg for param
				inputValue = currentScope[inputName].String
			} else {
				// no provided arg for param; fallback to default
				inputValue = stringInput.Default
			}

			// interpolate interpolatedStrings w/ inputValue
			for cmdEntryIndex, cmdEntry := range cmd {
				cmd[cmdEntryIndex] = interpolate.Interpolate(cmdEntry, inputName, inputValue)
			}
			for containerEnvVarName, containerEnvVar := range envVars {
				envVars[containerEnvVarName] = interpolate.Interpolate(containerEnvVar, inputName, inputValue)
			}
		}
	}
	return &containerprovider.RunContainerReq{
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
