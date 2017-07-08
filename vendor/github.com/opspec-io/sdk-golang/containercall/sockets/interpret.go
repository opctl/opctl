package sockets

import (
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path/filepath"
)

func (soc _Sockets) Interpret(
	scope map[string]*model.Value,
	scgContainerCallSockets map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallSockets := map[string]string{}
	for scgContainerSocketAddress, scgContainerSocketBind := range scgContainerCallSockets {
		if boundArg, ok := scope[scgContainerSocketBind]; ok {
			// bound to var
			dcgContainerCallSockets[scgContainerSocketAddress] = *boundArg.Socket
		} else if isUnixSocketAddress(scgContainerSocketAddress) {
			// bound to output
			// create outputSocket on host so the output points to something
			dcgHostSocketAddress := filepath.Join(scratchDirPath, scgContainerSocketAddress)
			var outputSocket *os.File
			outputSocket, err := soc.os.Create(dcgHostSocketAddress)
			outputSocket.Close()
			if nil != err {
				return nil, err
			}
			if err := soc.os.Chmod(
				dcgHostSocketAddress,
				os.ModeSocket,
			); nil != err {
				return nil, err
			}
			dcgContainerCallSockets[scgContainerSocketAddress] = dcgHostSocketAddress
		}
	}
	return dcgContainerCallSockets, nil
}
