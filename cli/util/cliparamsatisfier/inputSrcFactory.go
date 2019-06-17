package cliparamsatisfier

import (
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
)

//go:generate counterfeiter -o ./fakeInputSrc.go --fake-name FakeInputSrc ./ InputSrc

type InputSrc interface {
	// ReadString returns the value (if any) and true/false to indicate whether the read was successful
	ReadString(
		inputName string,
	) (*string, bool)
}

type InputSrcFactory interface {
	NewCliPromptInputSrc(
		inputs map[string]*model.Param,
	) InputSrc

	NewEnvVarInputSrc() InputSrc

	NewParamDefaultInputSrc(
		inputs map[string]*model.Param,
	) InputSrc

	NewSliceInputSrc(
		args []string,
		sep string,
	) InputSrc

	NewYMLFileInputSrc(
		filePath string,
	) (InputSrc, error)
}

func newInputSrcFactory() InputSrcFactory {
	return _InputSrcFactory{
		json:   ijson.New(),
		os:     ios.New(),
		ioutil: iioutil.New(),
	}
}

type _InputSrcFactory struct {
	json   ijson.IJSON
	os     ios.IOS
	ioutil iioutil.IIOUtil
}
