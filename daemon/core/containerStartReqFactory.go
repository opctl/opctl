package core

import (
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/sdk-golang/pkg/interpolate"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func newContainerStartReq(
	currentScope map[string]*model.Data,
	containerCall *model.ScgContainerCall,
	containerId string,
	inputs []*model.Param,
	opGraphId string,
) *containerengine.StartContainerReq {

	// create new slice so we don't mutate containerCall.Cmd @ caller
	cmd := append([]string{}, containerCall.Cmd...)
	dirs := map[string]string{}
	envVars := map[string]string{}
	files := map[string]string{}
	sockets := map[string]string{}

	// construct envVars
	for scgContainerEnvVarName, scgContainerEnvVar := range containerCall.EnvVars {
		envVars[scgContainerEnvVarName] = scgContainerEnvVar.Value
	}

	// construct files
	for scgContainerFilePath, scgContainerFile := range containerCall.Files {
		// @TODO: handle output files; (create temp file & pass path)
		if boundArg, ok := currentScope[scgContainerFile.Bind]; ok {
			files[scgContainerFilePath] = boundArg.File
		}
	}

	// construct dirs
	for scgContainerDirPath, scgContainerDir := range containerCall.Dirs {
		// @TODO: handle output dirs; (create temp dir & pass path)
		if boundArg, ok := currentScope[scgContainerDir.Bind]; ok {
			dirs[scgContainerDirPath] = boundArg.Dir
		}
	}

	// construct sockets
	for scgContainerSocketAddress, scgContainerSocket := range containerCall.Sockets {
		// @TODO: handle output (unix) sockets; (create temp unix socket & pass path)
		if boundArg, ok := currentScope[scgContainerSocket.Bind]; ok {
			switch {
			case "" != boundArg.Socket:
				sockets[scgContainerSocketAddress] = boundArg.Socket
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
		Image:       containerCall.Image,
		Sockets:     sockets,
		WorkDir:     containerCall.WorkDir,
		ContainerId: containerId,
		OpGraphId:   opGraphId,
	}

}
