package image

import (
	"errors"
	"fmt"
	"strings"

	"github.com/distribution/reference"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/str"
)

// Interpret container image
func Interpret(
	scope map[string]*model.Value,
	containerCallImageSpec *model.ContainerCallImageSpec,
	scratchDir string,
) (*model.ContainerCallImage, error) {

	if containerCallImageSpec == nil {
		return nil, errors.New("image required")
	}

	// try to interpret as dir
	src, err := dir.Interpret(
		scope,
		containerCallImageSpec.Ref,
		scratchDir,
		false,
	)
	if err == nil {
		return &model.ContainerCallImage{
			Src: src,
		}, nil
	}

	// fallback to string
	containerCallImage := &model.ContainerCallImage{}
	ref, err := str.Interpret(
		scope,
		containerCallImageSpec.Ref,
	)
	if err != nil {
		return nil, err
	}

	parsedRef, err := reference.ParseAnyReference(strings.ToLower(*ref.String))
	if err != nil {
		return nil, err
	}
	parsedRefString := parsedRef.String()

	containerCallImage.Ref = &parsedRefString

	if containerCallImageSpec.PullCreds != nil {
		username, err := str.Interpret(scope, containerCallImageSpec.PullCreds.Username)
		if err != nil {
			return nil, fmt.Errorf("error encountered interpreting image pullcreds username: %w", err)
		}

		password, err := str.Interpret(scope, containerCallImageSpec.PullCreds.Password)
		if err != nil {
			return nil, fmt.Errorf("error encountered interpreting image pullcreds password: %w", err)
		}

		containerCallImage.PullCreds = &model.Creds{
			Username: *username.String,
			Password: *password.String,
		}
	}
	return containerCallImage, nil
}
