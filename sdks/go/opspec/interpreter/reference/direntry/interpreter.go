package direntry

import (
	"fmt"
	"os"
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
		createTypeIfNotExists *string,
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
	createTypeIfNotExists *string,
) (string, *model.Value, error) {

	if !strings.HasPrefix(ref, "/") {
		return "", nil, fmt.Errorf("unable to interpret '%v'; expected '/'", ref)
	}

	valuePath := filepath.Join(*data.Dir, ref)

	fileInfo, err := itp.os.Stat(valuePath)
	if nil == err {
		if fileInfo.IsDir() {
			return "", &model.Value{Dir: &valuePath}, nil
		}

		return "", &model.Value{File: &valuePath}, nil
	} else if nil != createTypeIfNotExists && os.IsNotExist(err) {

		if "Dir" == *createTypeIfNotExists {
			err := os.MkdirAll(valuePath, 0700)
			if nil != err {
				return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref; error was %v", ref, err)
			}

			return "", &model.Value{Dir: &valuePath}, nil
		}

		// handle file ref
		err := os.MkdirAll(filepath.Dir(valuePath), 0700)
		if nil != err {
			return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref; error was %v", ref, err)
		}

		file, err := os.Create(valuePath)
		file.Close()
		if nil != err {
			return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref; error was %v", ref, err)
		}

		return "", &model.Value{File: &valuePath}, nil

	}

	return "", nil, fmt.Errorf("unable to interpret '%v' as dir entry ref; error was %v", ref, err)

}
