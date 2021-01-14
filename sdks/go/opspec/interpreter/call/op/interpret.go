package op

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/data/git"
	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/inputs"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// Interpret interprets an OpCallSpec into a OpCall
func Interpret(
	ctx context.Context,
	scope map[string]*model.Value,
	opCallSpec *model.OpCallSpec,
	opID string,
	parentOpPath string,
	dataDirPath string,
) (*model.OpCall, error) {

	scratchDirPath := filepath.Join(dataDirPath, "dcg", opID)

	var pkgPullCreds *model.Creds
	if pullCredsSpec := opCallSpec.PullCreds; nil != pullCredsSpec {
		pkgPullCreds = &model.Creds{}
		var err error
		interpretdUsername, err := str.Interpret(scope, pullCredsSpec.Username)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Username = *interpretdUsername.String

		interpretdPassword, err := str.Interpret(scope, pullCredsSpec.Password)
		if nil != err {
			return nil, err
		}
		pkgPullCreds.Password = *interpretdPassword.String
	}

	var opPath string
	if regexp.MustCompile("^\\$\\(.+\\)$").MatchString(opCallSpec.Ref) {
		// attempt to process as a variable reference since its variable reference like.
		dirValue, err := dir.Interpret(
			scope,
			opCallSpec.Ref,
			scratchDirPath,
			false,
		)
		if nil != err {
			return nil, fmt.Errorf("error encountered interpreting image src; error was: %v", err)
		}
		opPath = *dirValue.Dir
	} else {
		opHandle, err := data.Resolve(
			ctx,
			opCallSpec.Ref,
			fs.New(parentOpPath, filepath.Dir(parentOpPath)),
			git.New(filepath.Join(dataDirPath, "ops"), pkgPullCreds),
		)
		if nil != err {
			return nil, err
		}
		opPath = *opHandle.Path()
	}

	opFile, err := opfile.Get(
		ctx,
		opPath,
	)
	if nil != err {
		return nil, err
	}

	childCallID, err := uniquestring.Construct()
	if nil != err {
		return nil, err
	}

	opCall := &model.OpCall{
		BaseCall: model.BaseCall{
			OpPath: opPath,
		},
		ChildCallID:       childCallID,
		ChildCallCallSpec: opFile.Run,
		OpID:              opID,
	}

	opCall.Inputs, err = inputs.Interpret(
		opCallSpec.Inputs,
		opFile.Inputs,
		opPath,
		scope,
		scratchDirPath,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret call to %v; error was: %v", opCallSpec.Ref, err)
	}

	return opCall, nil
}
