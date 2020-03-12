package opfile

import (
	"context"
	"path/filepath"

	"github.com/golang-interfaces/iioutil"
	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o fakes/getter.go . Getter
type Getter interface {
	// Get gets the validated, deserialized representation of an "op.yml" file
	Get(
		ctx context.Context,
		opPath string,
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
	opPath string,
) (
	*model.OpFile,
	error,
) {
	opFileBytes, err := gtr.ioUtil.ReadFile(filepath.Join(opPath, FileName))
	if nil != err {
		return nil, err
	}

	return gtr.unmarshaller.Unmarshal(opFileBytes)
}
