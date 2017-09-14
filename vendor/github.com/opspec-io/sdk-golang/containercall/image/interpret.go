package image

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

func (img _Image) Interpret(
	scope map[string]*model.Value,
	scgContainerCallImage *model.SCGContainerCallImage,
	pkgHandle model.PkgHandle,
) (*model.DCGContainerCallImage, error) {
	// construct image
	if scgContainerCallImage := scgContainerCallImage; scgContainerCallImage != nil {
		ref, err := img.expression.EvalToString(scope, scgContainerCallImage.Ref, pkgHandle)
		if nil != err {
			return nil, err
		}

		dcgContainerCallImage := &model.DCGContainerCallImage{
			Ref: ref,
		}

		if nil != scgContainerCallImage.PullCreds {
			username, err := img.expression.EvalToString(scope, scgContainerCallImage.PullCreds.Username, pkgHandle)
			if nil != err {
				return nil, err
			}

			password, err := img.expression.EvalToString(scope, scgContainerCallImage.PullCreds.Password, pkgHandle)
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
