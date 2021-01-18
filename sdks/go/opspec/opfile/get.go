package opfile

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
)

// Get gets the validated, deserialized representation of an "op.yml" file
func Get(
	ctx context.Context,
	opPath string,
) (
	*model.OpSpec,
	error,
) {
	opFileBytes, err := ioutil.ReadFile(filepath.Join(opPath, FileName))
	if nil != err {
		return nil, err
	}

	return Unmarshal(opFileBytes)
}
