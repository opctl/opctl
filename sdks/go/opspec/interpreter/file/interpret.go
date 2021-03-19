package file

import (
	"fmt"
	"regexp"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/pkg/errors"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

// Interpret an expression to a file value.
// Expression must be a type supported by coerce.ToFile
// scratchDir will be used as the containing dir if file creation necessary
//
// Examples of valid file expressions:
// scope ref: $(scope-ref)
// scope ref w/ path: $(scope-ref/file.txt)
// pkg fs ref: $(/pkg-fs-ref)
// pkg fs ref w/ path: $(/pkg-fs-ref/file.txt)
func Interpret(
	scope map[string]*model.Value,
	expression interface{},
	scratchDir string,
	createIfNotExist bool,
) (*model.Value, error) {
	expressionAsString, expressionIsString := expression.(string)

	if expressionIsString && regexp.MustCompile("^\\$\\(.+\\)$").MatchString(expressionAsString) {
		var opts *model.ReferenceOpts
		if createIfNotExist {
			opts = &model.ReferenceOpts{
				Type:       "File",
				ScratchDir: scratchDir,
			}
		}

		value, err := reference.Interpret(
			expressionAsString,
			scope,
			opts,
		)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to file", expression))
		}
		return coerce.ToFile(value, scratchDir)
	}

	value, err := value.Interpret(
		expression,
		scope,
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to interpret %+v to file", expression))
	}

	return coerce.ToFile(&value, scratchDir)
}
