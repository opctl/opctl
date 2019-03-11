package dotyml

//go:generate counterfeiter -o ./fakeGetter.go --fake-name FakeGetter ./ Getter

import (
	"context"
	"github.com/golang-interfaces/iioutil"
	"github.com/opctl/sdk-golang/model"
)

type Getter interface {
	// Get gets the validated, deserialized representation of an "op.yml" file
	Get(
		ctx context.Context,
		opHandle model.DataHandle,
	) (
		*model.OpDotYml,
		error,
	)
}

func NewGetter() Getter {
	return _getter{
		ioUtil:       iioutil.New(),
		unmarshaller: NewUnmarshaller(),
	}
}

type _getter struct {
	ioUtil       iioutil.IIOUtil
	unmarshaller Unmarshaller
}

func (gtr _getter) Get(
	ctx context.Context,
	opHandle model.DataHandle,
) (
	*model.OpDotYml,
	error,
) {
	manifestReader, err := opHandle.GetContent(ctx, FileName)
	if nil != err {
		return nil, err
	}
	defer manifestReader.Close()

	manifestBytes, err := gtr.ioUtil.ReadAll(manifestReader)
	if nil != err {
		return nil, err
	}

	return gtr.unmarshaller.Unmarshal(manifestBytes)
}
