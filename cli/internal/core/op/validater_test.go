package op

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	dataresolver "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("Validater", func() {
	Context("Validate", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"

			fakeDataResolver := new(dataresolver.FakeDataResolver)
			opPath := "opPath"
			fakeOpHandle := new(FakeDataHandle)
			fakeOpHandle.PathReturns(&opPath)
			fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

			objectUnderTest := _validater{
				dataResolver: fakeDataResolver,
			}

			/* act */
			objectUnderTest.Validate(
				context.Background(),
				providedPkgRef,
			)

			/* assert */
			actualPkgRef,
				actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)

			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("returns errors from op resolution", func() {
			/* arrange */
			expectedError := errors.New("expected")

			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(nil, expectedError)

			objectUnderTest := newValidater(fakeDataResolver)

			/* act */
			_, err := objectUnderTest.Validate(context.Background(), "dummy")

			/* assert */
			Expect(err).To(MatchError(expectedError))
		})
		Context("op.Validate returns errors", func() {
			It("should call cliExiter.Exit w/ expected args", func() {
				/* arrange */
				fakeDataResolver := new(dataresolver.FakeDataResolver)

				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.PathReturns(new(string))
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _validater{
					dataResolver: fakeDataResolver,
				}

				/* act */
				_, err := objectUnderTest.Validate(
					context.Background(),
					"dummyPkgRef",
				)

				/* assert */
				Expect(err).To(MatchError("open op.yml: no such file or directory"))
			})
		})
		Context("pkg.Validate doesn't return errors", func() {
			It("should call cliExiter.Exit w/ expected args", func() {
				/* arrange */
				wd, err := os.Getwd()
				if nil != err {
					panic(err)
				}
				opRef := filepath.Join(wd, "testdata/validater_valid")

				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.PathReturns(&opRef)
				fakeOpHandle.RefReturns(opRef)

				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _validater{
					dataResolver: fakeDataResolver,
				}

				/* act */
				message, err := objectUnderTest.Validate(
					context.Background(),
					opRef,
				)

				/* assert */
				Expect(err).To(BeNil())
				Expect(message).To(Equal(fmt.Sprintf("%v is valid", opRef)))
			})
		})
	})
})
