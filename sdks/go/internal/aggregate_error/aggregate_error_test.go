package errors

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

func TestAggregateError(t *testing.T) {
	g := NewGomegaWithT(t)

	/* arrange */
	internalErr := errors.New("testing")
	err := ErrAggregate{
		errs: []error{
			errors.Wrap(ErrAggregate{
				errs: []error{
					errors.New("nested"),
					model.ErrDataProviderAuthorization{},
				},
			}, "container"),
			internalErr,
		},
	}

	/* assert */
	t.Log(err.Error())
	g.Expect(err.Error()).To(Equal(`
- container:` + " " + `
  - nested
  - unauthorized
- testing`))
	g.Expect(err.Is(internalErr)).To(BeTrue())
	g.Expect(err.Is(errors.New("garbage"))).To(BeFalse())
}

func TestAddError(t *testing.T) {
	g := NewGomegaWithT(t)

	innerErr := errors.New("testing")

	err := ErrAggregate{}
	err.AddError(innerErr)

	g.Expect(err).To(MatchError(ErrAggregate{
		errs: []error{innerErr},
	}))
}
