package sockets

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
	"os"
	"path/filepath"
)

type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallSockets map[string]string,
		scratchDirPath string,
	) (map[string]string, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		os: ios.New(),
	}
}

type _interpreter struct {
	os ios.IOS
}

func (itp _interpreter) Interpret(
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
			outputSocket, err := itp.os.Create(dcgHostSocketAddress)
			outputSocket.Close()
			if nil != err {
				return nil, err
			}
			if err := itp.os.Chmod(
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
