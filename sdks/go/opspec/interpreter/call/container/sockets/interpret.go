package sockets

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
)

// Interpret container sockets
func Interpret(
	scope map[string]*model.Value,
	containerCallSpecSockets map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	containerCallSockets := map[string]string{}
	for callSpecContainerSocketAddress, callSpecContainerSocketBind := range containerCallSpecSockets {
		// @TODO: use reference.interpret once reference syntax no longer optional
		callSpecContainerSocketBind = strings.TrimSuffix(strings.TrimPrefix(callSpecContainerSocketBind, reference.RefStart), reference.RefEnd)

		if boundArg, ok := scope[callSpecContainerSocketBind]; ok {
			// bound to var
			containerCallSockets[callSpecContainerSocketAddress] = *boundArg.Socket
		} else if isUnixSocketAddress(callSpecContainerSocketAddress) {
			// bound to output
			// create outputSocket on host so the output points to something
			dcgHostSocketAddress := filepath.Join(scratchDirPath, callSpecContainerSocketAddress)
			var outputSocket *os.File
			outputSocket, err := os.Create(dcgHostSocketAddress)
			outputSocket.Close()
			if err != nil {
				return nil, err
			}
			if err := os.Chmod(
				dcgHostSocketAddress,
				os.ModeSocket,
			); err != nil {
				return nil, err
			}
			containerCallSockets[callSpecContainerSocketAddress] = dcgHostSocketAddress
		}
	}
	return containerCallSockets, nil
}
