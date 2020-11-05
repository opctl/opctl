package op

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("Validater", func() {
	Context("Validate", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedPkgRef := "dummyPkgRef"

			fakeDataResolver := new(dataresolver.Fake)
			opPath := "opPath"
			fakeOpHandle := new(FakeDataHandle)
			fakeOpHandle.PathReturns(&opPath)
			fakeDataResolver.ResolveReturns(fakeOpHandle)

			objectUnderTest := _validater{
				cliExiter:    new(cliexiterFakes.FakeCliExiter),
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
		Context("op.Validate returns errors", func() {
			It("should call cliExiter.Exit w/ expected args", func() {
				/* arrange */
				fakeDataResolver := new(dataresolver.Fake)

				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.PathReturns(new(string))
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				expectedExitReq := cliexiter.ExitReq{
					Message: "open op.yml: no such file or directory",
					Code:    1,
				}

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _validater{
					cliExiter:    fakeCliExiter,
					dataResolver: fakeDataResolver,
				}

				/* act */
				objectUnderTest.Validate(
					context.Background(),
					"dummyPkgRef",
				)

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
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

				fakeDataResolver := new(dataresolver.Fake)
				fakeDataResolver.ResolveReturns(fakeOpHandle)

				expectedExitReq := cliexiter.ExitReq{
					Message: fmt.Sprintf("%v is valid", opRef),
				}

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _validater{
					cliExiter:    fakeCliExiter,
					dataResolver: fakeDataResolver,
				}

				/* act */
				objectUnderTest.Validate(
					context.Background(),
					opRef,
				)

				/* assert */

				Expect(fakeCliExiter.ExitArgsForCall(0)).To(Equal(expectedExitReq))
			})
		})
	})
})
