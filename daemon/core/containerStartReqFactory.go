package core

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/interpolater"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func newContainerStartReq(
	args map[string]*model.Data,
	containerCall *model.ContainerCall,
	containerId string,
	inputs []*model.Param,
	opGraphId string,
) *containerengine.StartContainerReq {

	cmd := append([]string{}, containerCall.Cmd...)
	env := []*model.ContainerEnvEntry{}
	fs := []*model.ContainerFsEntry{}
	net := []*model.ContainerNetEntry{}

	// construct fs
	for _, staticCallFsEntry := range containerCall.Fs {
		srcRef := ""
		if srcRefData, ok := args[staticCallFsEntry.Bind]; ok {
			switch {
			case "" != srcRefData.Dir:
				srcRef = srcRefData.Dir
			case "" != srcRefData.File:
				srcRef = srcRefData.File
			}
		}
		fsEntry := &model.ContainerFsEntry{
			Path:   staticCallFsEntry.Path,
			SrcRef: srcRef,
		}
		fs = append(fs, fsEntry)
	}
	fmt.Printf("startFactory: containerFs\n%#v\n", fs)

	for _, input := range inputs {
		fmt.Printf("containerStartReqFactory.input:\n%#v\n", input)
		switch {
		case nil != input.String:
			stringInput := input.String
			inputValue := ""
			if _, isArgForInput := args[stringInput.Name]; isArgForInput {
				// use provided arg for param
				inputValue = args[stringInput.Name].String
			} else {
				// no provided arg for param; fallback to default
				inputValue = stringInput.Default
			}
			for cmdEntryIndex, cmdEntry := range cmd {
				cmd[cmdEntryIndex] = interpolater.Interpolate(cmdEntry, stringInput.Name, inputValue)
			}
			for _, envEntry := range containerCall.Env {
				// append bound strings to env
				if envEntry.Bind == stringInput.Name {
					env = append(
						env,
						&model.ContainerEnvEntry{
							Name:  stringInput.Name,
							Value: inputValue,
						},
					)
					break
				}
			}
		case nil != input.NetSocket:
			netSocketInput := input.NetSocket
			netSocketArg := args[netSocketInput.Name].NetSocket
			for _, netEntry := range containerCall.Net {
				// append bound sockets to net
				if netEntry.Bind == netSocketInput.Name {
					net = append(
						net,
						&model.ContainerNetEntry{
							Host:        netSocketArg.Host,
							Port:        netSocketArg.Port,
							HostAliases: netEntry.HostAliases,
						},
					)
					break
				}
			}
		}
	}
	return &containerengine.StartContainerReq{
		Cmd:         cmd,
		Env:         env,
		Fs:          fs,
		Image:       containerCall.Image,
		Net:         net,
		WorkDir:     containerCall.WorkDir,
		ContainerId: containerId,
		OpGraphId:   opGraphId,
	}

}
