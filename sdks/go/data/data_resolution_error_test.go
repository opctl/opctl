package data

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

func TestErrDataResolution(t *testing.T) {
	g := NewGomegaWithT(t)

	/* arrange */
	internalErr := errors.New("testing")
	err := ErrDataResolution{
		dataRef: "opref",
		errs: []error{
			ErrDataResolution{
				dataRef: "opref",
				errs: []error{
					errors.New("nested"),
					model.ErrDataProviderAuthorization{},
				},
			},
			internalErr,
		},
	}

	/* assert */
	g.Expect(err.Error()).To(Equal(`unable to resolve op "opref":
- unable to resolve op "opref":
  - nested
  - unauthorized
- testing`))
	g.Expect(err.Is(internalErr)).To(BeTrue())
	g.Expect(err.Is(errors.New("garbage"))).To(BeFalse())
}
