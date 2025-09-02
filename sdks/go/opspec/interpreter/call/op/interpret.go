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
	if pullCredsSpec := opCallSpec.PullCreds; pullCredsSpec != nil {
		pkgPullCreds = &model.Creds{}
		var err error
		interpretdUsername, err := str.Interpret(scope, pullCredsSpec.Username)
		if err != nil {
			return nil, err
		}
		pkgPullCreds.Username = *interpretdUsername.String

		interpretdPassword, err := str.Interpret(scope, pullCredsSpec.Password)
		if err != nil {
			return nil, err
		}
		pkgPullCreds.Password = *interpretdPassword.String
	}

	var opDir model.DataHandle
	if regexp.MustCompile(`^\$\(.+\)$`).MatchString(opCallSpec.Ref) {
		// attempt to process as a variable reference since its variable reference like.
		dirValue, err := dir.Interpret(
			scope,
			opCallSpec.Ref,
			scratchDirPath,
			false,
		)
		if err != nil {
			return nil, fmt.Errorf("error encountered interpreting image src: %w", err)
		}

		opDir, err = fs.New().TryResolve(ctx, *dirValue.Dir)
		if err != nil {
			return nil, fmt.Errorf("error encountered interpreting image src: %w", err)
		}
	} else {
		var err error
		opDir, err = data.Resolve(
			ctx,
			opCallSpec.Ref,
			fs.New(parentOpPath, filepath.Dir(parentOpPath)),
			git.New(filepath.Join(dataDirPath, "ops"), pkgPullCreds),
		)
		if err != nil {
			return nil, err
		}
	}

	opFile, err := opfile.Get(
		ctx,
		opDir,
	)
	if err != nil {
		return nil, err
	}

	childCallID, err := uniquestring.Construct()
	if err != nil {
		return nil, err
	}

	// this relies on op existing locally which is currently always true since we only use fs & git as data providers
	opPath := opDir.Ref()
	if !filepath.IsAbs(opPath) {
		opPath = filepath.Join(dataDirPath, "ops", opPath)
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
	if err != nil {
		return nil, fmt.Errorf("unable to interpret call to %v: %w", opCallSpec.Ref, err)
	}

	return opCall, nil
}
