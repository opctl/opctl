package op

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("Installer", func() {
	Context("Install", func() {
		It("should call dataResolver w/ expected args", func() {
			/* arrange */
			providedPkgRef := "providedPkgRef#0.0.0"
			providedOpRef := fmt.Sprintf("%s/subpath", providedPkgRef)
			providedPullCreds := &model.Creds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			fakeDataResolver := new(dataresolver.Fake)
			fakeDataResolver.ResolveReturns(new(FakeDataHandle))

			objectUnderTest := _installer{
				dataResolver: fakeDataResolver,
			}

			/* act */
			objectUnderTest.Install(
				context.Background(),
				"dummyPath",
				providedOpRef,
				providedPullCreds.Username,
				providedPullCreds.Password,
			)

			/* assert */
			actualPkgRef,
				actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)

			Expect(actualPkgRef).To(Equal(providedPkgRef))
			Expect(actualPullCreds).To(Equal(providedPullCreds))
		})
		Context("op.Install errs", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.ListDescendantsReturns(nil, errors.New(""))

				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				objectUnderTest := _installer{
					cliExiter:    fakeCliExiter,
					dataResolver: fakeDataResolver,
				}

				/* act */
				objectUnderTest.Install(
					context.Background(),
					"",
					"",
					"",
					"",
				)

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: "", Code: 1}))

			})
		})
	})
})
