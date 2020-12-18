package op

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	dataresolver "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("Installer", func() {
	Context("Install", func() {
		It("should return dataResolver errors", func() {
			/* arrange */
			expectedError := errors.New("expected")
			providedPkgRef := "providedPkgRef#0.0.0"
			providedOpRef := fmt.Sprintf("%s/subpath", providedPkgRef)
			providedPullCreds := &model.Creds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}
			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(nil, expectedError)

			objectUnderTest := newInstaller(fakeDataResolver)

			/* act */
			err := objectUnderTest.Install(
				context.Background(),
				"dummyPath",
				providedOpRef,
				providedPullCreds.Username,
				providedPullCreds.Password,
			)

			/* assert */
			Expect(err).To(MatchError(expectedError))
		})
		It("should call dataResolver w/ expected args", func() {
			/* arrange */
			providedPkgRef := "providedPkgRef#0.0.0"
			providedOpRef := fmt.Sprintf("%s/subpath", providedPkgRef)
			providedPullCreds := &model.Creds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(new(FakeDataHandle), nil)

			objectUnderTest := _installer{
				dataResolver: fakeDataResolver,
			}

			/* act */
			err := objectUnderTest.Install(
				context.Background(),
				"dummyPath",
				providedOpRef,
				providedPullCreds.Username,
				providedPullCreds.Password,
			)

			/* assert */
			actualPkgRef,
				actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)
			Expect(err).To(BeNil())
			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(Equal(providedPullCreds))
		})
		Context("op.Install errs", func() {
			It("should return expected error", func() {
				/* arrange */
				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.ListDescendantsReturns(nil, errors.New(""))

				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _installer{
					dataResolver: fakeDataResolver,
				}

				/* act */
				err := objectUnderTest.Install(
					context.Background(),
					"",
					"",
					"",
					"",
				)

				/* assert */
				Expect(err).To(MatchError(""))
			})
		})
	})
})
