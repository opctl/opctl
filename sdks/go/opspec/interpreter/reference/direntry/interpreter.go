package direntry

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/types"
)

// Interpreter interprets a dir entry ref i.e. refs of the form name/sub/file.ext
// it's an error if ref doesn't start with '/'
// returns ref remainder, dereferenced data, and error if one occurred
type Interpreter interface {
	Interpret(
		ref string,
		data *types.Value,
	) (string, *types.Value, error)
}

func NewInterpreter() Interpreter {
	return _interpreter{
		os: ios.New(),
	}
}

type _interpreter struct {
	os ios.IOS
}

func (itp _interpreter) Interpret(
	ref string,
	data *types.Value,
) (string, *types.Value, error) {

	if !strings.HasPrefix(ref, "/") {
		return "", nil, fmt.Errorf("unable to interpret '%v'; expected '/'", ref)
	}

	valuePath := filepath.Join(*data.Dir, ref)

	fileInfo, err := itp.os.Stat(valuePath)
	if nil != err {
		return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref; error was %v", ref, err)
	}

	if fileInfo.IsDir() {
		return "", &types.Value{Dir: &valuePath}, nil
	}

	return "", &types.Value{File: &valuePath}, nil

}
