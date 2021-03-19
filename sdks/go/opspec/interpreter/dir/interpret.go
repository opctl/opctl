package dir

import (
	"fmt"
	"regexp"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
	"github.com/pkg/errors"
)

// Interpret an expression to a dir value.
// Expression must be of type string.
//
// Examples of valid dir expressions:
// scope ref: $(scope-ref)
// scope ref w/ path: $(scope-ref/sub-dir)
// pkg fs ref: $(/pkg-fs-ref)
// pkg fs ref w/ path: $(/pkg-fs-ref/sub-dir)
func Interpret(
	scope map[string]*model.Value,
	expression interface{},
	scratchDir string,
	createIfNotExist bool,
) (*model.Value, error) {
	switch expression := expression.(type) {
	case string:
		if regexp.MustCompile("^\\$\\(.+\\)$").MatchString(expression) {
			var opts *model.ReferenceOpts
			if createIfNotExist {
				opts = &model.ReferenceOpts{
					Type:       "Dir",
					ScratchDir: scratchDir,
				}
			}

			value, err := reference.Interpret(
				expression,
				scope,
				opts,
			)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to dir", expression))
			}

			result, err := coerce.ToDir(value, scratchDir)
			if err != nil {
				err = errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to dir", expression))
			}
			return result, err

		}
	}

	value, err := value.Interpret(
		expression,
		scope,
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to dir", expression))
	}

	result, err := coerce.ToDir(&value, scratchDir)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to dir", expression))
	}
	return result, err
}
