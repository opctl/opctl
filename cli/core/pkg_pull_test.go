package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("core", func() {
	Context("PkgPull", func() {
		It("should call pkg.Pull w/ expected args", func() {
			/* arrange */
			fakePkg := new(pkg.Fake)

			providedPkgRef := "dummyPkgRef"
			providedUsername := "dummyUsername"
			providedPassword := "dummyPassword"

			expectedPkgRef := providedPkgRef
			expectedPullOpts := &pkg.PullOpts{
				Username: providedUsername,
				Password: providedPassword,
			}

			objectUnderTest := _core{
				pkg: fakePkg,
			}

			/* act */
			objectUnderTest.PkgPull(providedPkgRef, providedUsername, providedPassword)

			/* assert */
			actualPkgRef,
				actualPullOpts := fakePkg.PullArgsForCall(0)

			Expect(actualPkgRef).To(Equal(expectedPkgRef))
			Expect(actualPullOpts).To(Equal(expectedPullOpts))
		})
		Context("pkg.Pull errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakePkg := new(pkg.Fake)
				expectedError := errors.New("dummyError")
				fakePkg.PullReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					pkg:       fakePkg,
					cliExiter: fakeCliExiter,
				}

				/* act */
				objectUnderTest.PkgPull("", "", "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))

			})
		})
	})
})
