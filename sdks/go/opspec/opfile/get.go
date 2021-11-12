package opfile

import (
	"context"
	"os"
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
	opFileBytes, err := os.ReadFile(filepath.Join(opPath, FileName))
	if err != nil {
		return nil, err
	}

	return Unmarshal(opFileBytes)
}
