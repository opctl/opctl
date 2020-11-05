package image

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
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

	if nil == containerCallImageSpec {
		return nil, fmt.Errorf("image required")
	}

	// try to interpret as dir
	src, err := dir.Interpret(
		scope,
		containerCallImageSpec.Ref,
		scratchDir,
		false,
	)
	if nil == err {
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
	if nil != err {
		return nil, err
	}

	parsedRef, err := reference.ParseAnyReference(strings.ToLower(*ref.String))
	if nil != err {
		return nil, err
	}
	parsedRefString := parsedRef.String()

	containerCallImage.Ref = &parsedRefString

	if nil != containerCallImageSpec.PullCreds {
		username, err := str.Interpret(scope, containerCallImageSpec.PullCreds.Username)
		if nil != err {
			return nil, fmt.Errorf("error encountered interpreting image pullcreds username; error was: %v", err)
		}

		password, err := str.Interpret(scope, containerCallImageSpec.PullCreds.Password)
		if nil != err {
			return nil, fmt.Errorf("error encountered interpreting image pullcreds password; error was: %v", err)
		}

		containerCallImage.PullCreds = &model.Creds{
			Username: *username.String,
			Password: *password.String,
		}
	}
	return containerCallImage, nil
}
