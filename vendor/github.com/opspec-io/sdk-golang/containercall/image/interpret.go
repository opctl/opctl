package image

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

func (img _Image) Interpret(
	scope map[string]*model.Value,
	scgContainerCallImage *model.SCGContainerCallImage,
) (*model.DCGContainerCallImage, error) {
	// construct image
	if scgContainerCallImage := scgContainerCallImage; scgContainerCallImage != nil {
		ref, err := img.string.Interpret(scope, scgContainerCallImage.Ref)
		if nil != err {
			return nil, err
		}

		dcgContainerCallImage := &model.DCGContainerCallImage{
			Ref: ref,
		}

		if nil != scgContainerCallImage.PullCreds {
			username, err := img.string.Interpret(scope, scgContainerCallImage.PullCreds.Username)
			if nil != err {
				return nil, err
			}

			password, err := img.string.Interpret(scope, scgContainerCallImage.PullCreds.Password)
			if nil != err {
				return nil, err
			}

			dcgContainerCallImage.PullCreds = &model.DCGPullCreds{
				Username: username,
				Password: password,
			}
		}
		return dcgContainerCallImage, nil
	}
	return nil, errors.New("image required")
}
