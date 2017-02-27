package core

import (
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/sdk-golang/pkg/interpolate"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"path"
)

func newRunContainerReq(
	currentScope map[string]*model.Data,
	scgContainerCall *model.ScgContainerCall,
	containerId string,
	opGraphId string,
) *containerprovider.RunContainerReq {

	// create new slice so we don't cause side effects
	cmd := append([]string{}, scgContainerCall.Cmd...)
	dirs := map[string]string{}
	envVars := map[string]string{}
	files := map[string]string{}
	sockets := map[string]string{}

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
	err := os.MkdirAll(scratchDirPath, 0700)
	if nil != err {
		panic(err)
	}

	// construct dirs
	for scgContainerDirPath, scgContainerDir := range scgContainerCall.Dirs {
		if boundArg, ok := currentScope[scgContainerDir.Bind]; ok {
			// bound to input
			dirs[scgContainerDirPath] = boundArg.Dir
		} else {
			// bound to output
			// create placeholder dir on host so the output points to something
			dcgHostDirPath := path.Join(scratchDirPath, scgContainerDirPath)
			err := os.MkdirAll(dcgHostDirPath, 0700)
			if nil != err {
				panic(err)
			}
			files[scgContainerDirPath] = dcgHostDirPath
		}
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
			// create outputFile on host so the output points to something
			dcgHostFilePath := path.Join(scratchDirPath, scgContainerFilePath)
			// create dir
			err := os.MkdirAll(path.Dir(dcgHostFilePath), 0700)
			if nil != err {
				panic(err)
			}
			// create file
			outputFile, err := os.Create(dcgHostFilePath)
			outputFile.Close()
			if nil != err {
				panic(err)
			}
			files[scgContainerFilePath] = dcgHostFilePath
		}
	}

	// construct sockets
	for scgContainerSocketAddress, scgContainerSocket := range scgContainerCall.Sockets {
		if boundArg, ok := currentScope[scgContainerSocket.Bind]; ok {
			// bound to input
			sockets[scgContainerSocketAddress] = boundArg.Socket
		} else if isUnixSocketAddress(scgContainerSocketAddress) {
			// bound to output
			// create outputSocket on host so the output points to something
			if isUnixSocketAddress(scgContainerSocketAddress) {
				dcgHostSocketAddress := path.Join(scratchDirPath, scgContainerSocketAddress)
				outputSocket, err := os.Create(dcgHostSocketAddress)
				outputSocket.Close()
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

	// interpolate cmd & envVars w/ values from currentScope
	for varName, varData := range currentScope {
		switch {
		case 0 != varData.Number:
			numberVarData := varData.Number

			for cmdEntryIndex, cmdEntry := range cmd {
				cmd[cmdEntryIndex] = interpolate.NumberValue(cmdEntry, varName, numberVarData)
			}
			for containerEnvVarName, containerEnvVar := range envVars {
				envVars[containerEnvVarName] = interpolate.NumberValue(containerEnvVar, varName, numberVarData)
			}
		case "" != varData.String:
			stringVarData := varData.String

			for cmdEntryIndex, cmdEntry := range cmd {
				cmd[cmdEntryIndex] = interpolate.StringValue(cmdEntry, varName, stringVarData)
			}
			for containerEnvVarName, containerEnvVar := range envVars {
				envVars[containerEnvVarName] = interpolate.StringValue(containerEnvVar, varName, stringVarData)
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
