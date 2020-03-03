package opfile

import (
	"context"

	"github.com/golang-interfaces/iioutil"
	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o fakes/getter.go . Getter
type Getter interface {
	// Get gets the validated, deserialized representation of an "op.yml" file
	Get(
		ctx context.Context,
		opHandle model.DataHandle,
	) (
		*model.OpFile,
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
	*model.OpFile,
	error,
) {
	opFileReader, err := opHandle.GetContent(ctx, FileName)
	if nil != err {
		return nil, err
	}
	defer opFileReader.Close()

	opFileBytes, err := gtr.ioUtil.ReadAll(opFileReader)
	if nil != err {
		return nil, err
	}

	return gtr.unmarshaller.Unmarshal(opFileBytes)
}
