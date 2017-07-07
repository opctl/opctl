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
		dcgContainerCallImage := &model.DCGContainerCallImage{
			Ref: img.interpolater.Interpolate(scgContainerCallImage.Ref, scope),
		}
		if "" != scgContainerCallImage.PullIdentity && "" != scgContainerCallImage.PullSecret {
			// fallback for deprecated cred format
			scgContainerCallImage.PullCreds = &model.SCGPullCreds{
				Username: scgContainerCallImage.PullIdentity,
				Password: scgContainerCallImage.PullSecret,
			}
		}

		if nil != scgContainerCallImage.PullCreds {
			dcgContainerCallImage.PullCreds = &model.DCGPullCreds{
				Username: img.interpolater.Interpolate(scgContainerCallImage.PullCreds.Username, scope),
				Password: img.interpolater.Interpolate(scgContainerCallImage.PullCreds.Password, scope),
			}
		}
		return dcgContainerCallImage, nil
	}
	return nil, errors.New("image required")
}
