package dir

import (
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
)

// Interpret interprets an expression to a dir value.
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
		// @TODO: this incorrectly treats $(inScope)$(inScope) as ref
		if strings.HasPrefix(expression, interpolater.RefStart) && strings.HasSuffix(expression, interpolater.RefEnd) {
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
			if nil != err {
				return nil, fmt.Errorf("unable to interpret %+v to dir; error was %v", expression, err)
			}

			if nil == value.Dir {
				return nil, fmt.Errorf("unable to interpret %+v to dir", expression)
			}

			return value, nil

		}
	}

	return nil, fmt.Errorf("unable to interpret %v to dir", expression)

}
