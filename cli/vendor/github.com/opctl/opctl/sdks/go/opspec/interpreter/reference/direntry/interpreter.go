package direntry

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/model"
)

// Interpreter interprets a dir entry ref i.e. refs of the form name/sub/file.ext
// it's an error if ref doesn't start with '/'
// returns ref remainder, dereferenced data, and error if one occurred
//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		ref string,
		data *model.Value,
	) (string, *model.Value, error)
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
	data *model.Value,
) (string, *model.Value, error) {

	if !strings.HasPrefix(ref, "/") {
		return "", nil, fmt.Errorf("unable to interpret '%v'; expected '/'", ref)
	}

	valuePath := filepath.Join(*data.Dir, ref)

	fileInfo, err := itp.os.Stat(valuePath)
	if nil != err {
		return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref; error was %v", ref, err)
	}

	if fileInfo.IsDir() {
		return "", &model.Value{Dir: &valuePath}, nil
	}

	return "", &model.Value{File: &valuePath}, nil

}
