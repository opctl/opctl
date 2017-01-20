package core

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/interpolater"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func newContainerStartReq(
	args map[string]*model.Data,
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
		if boundArg, ok := args[scgContainerFile.Bind]; ok {
			switch {
			case "" != boundArg.File:
				files[scgContainerFilePath] = boundArg.File
			}
		}
	}
	fmt.Printf("startFactory: req.files\n%#v\n", files)

	// construct dirs
	for scgContainerDirPath, scgContainerDir := range containerCall.Dirs {
		if boundArg, ok := args[scgContainerDir.Bind]; ok {
			switch {
			case "" != boundArg.Dir:
				dirs[scgContainerDirPath] = boundArg.Dir
			}
		}
	}
	fmt.Printf("startFactory: req.dirs\n%#v\n", dirs)

	// construct sockets
	for scgContainerSocketAddress, scgContainerSocket := range containerCall.Sockets {
		if boundArg, ok := args[scgContainerSocket.Bind]; ok {
			switch {
			case "" != boundArg.Socket:
				sockets[scgContainerSocketAddress] = boundArg.Dir
			}
		}
	}
	fmt.Printf("startFactory: req.sockets\n%#v\n", sockets)

	for _, input := range inputs {
		fmt.Printf("containerStartReqFactory.input:\n%#v\n", input)
		switch {
		case nil != input.String:
			stringInput := input.String

			// obtain inputValue
			inputValue := ""
			if _, isArgForInput := args[stringInput.Name]; isArgForInput {
				// use provided arg for param
				inputValue = args[stringInput.Name].String
			} else {
				// no provided arg for param; fallback to default
				inputValue = stringInput.Default
			}

			// interpolate interpolatedStrings w/ inputValue
			for cmdEntryIndex, cmdEntry := range cmd {
				cmd[cmdEntryIndex] = interpolater.Interpolate(cmdEntry, stringInput.Name, inputValue)
			}
			for containerEnvVarName, containerEnvVar := range envVars {
				envVars[containerEnvVarName] = interpolater.Interpolate(containerEnvVar, stringInput.Name, inputValue)
			}
		}
	}
	return &containerengine.StartContainerReq{
		Cmd:         cmd,
		Dirs:        dirs,
		Env:         envVars,
		Files:       files,
		Image:       containerCall.Image,
		Net:         sockets,
		WorkDir:     containerCall.WorkDir,
		ContainerId: containerId,
		OpGraphId:   opGraphId,
	}

}
