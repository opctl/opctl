package sockets

import (
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"os"
	"path/filepath"
	"strings"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		scope map[string]*model.Value,
		callContainerSpecSockets map[string]string,
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
	callContainerSpecSockets map[string]string,
	scratchDirPath string,
) (map[string]string, error) {
	dcgContainerCallSockets := map[string]string{}
	for callSpecContainerSocketAddress, callSpecContainerSocketBind := range callContainerSpecSockets {
		// @TODO: use reference.interpret once reference syntax no longer optional
		callSpecContainerSocketBind = strings.TrimSuffix(strings.TrimPrefix(callSpecContainerSocketBind, reference.RefStart), reference.RefEnd)

		if boundArg, ok := scope[callSpecContainerSocketBind]; ok {
			// bound to var
			dcgContainerCallSockets[callSpecContainerSocketAddress] = *boundArg.Socket
		} else if isUnixSocketAddress(callSpecContainerSocketAddress) {
			// bound to output
			// create outputSocket on host so the output points to something
			dcgHostSocketAddress := filepath.Join(scratchDirPath, callSpecContainerSocketAddress)
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
			dcgContainerCallSockets[callSpecContainerSocketAddress] = dcgHostSocketAddress
		}
	}
	return dcgContainerCallSockets, nil
}
