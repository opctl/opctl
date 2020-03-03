package node

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	cliexiterFakes "github.com/opctl/opctl/cli/internal/cliexiter/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

var _ = Context("Killer", func() {
	Context("Kill", func() {
		It("should call nodeProvider.Invoke w/ expected args", func() {
			/* arrange */
			fakeNodeProvider := new(nodeprovider.Fake)

			objectUnderTest := _killer{
				nodeProvider: fakeNodeProvider,
			}

			/* act */
			objectUnderTest.Kill()

			/* assert */
			Expect(fakeNodeProvider.KillNodeIfExistsArgsForCall(0)).To(BeEquivalentTo(""))
		})
		Context("nodeProvider.Invoke errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeNodeProvider := new(nodeprovider.Fake)
				expectedError := errors.New("dummyError")
				fakeNodeProvider.KillNodeIfExistsReturns(expectedError)

				fakeCliExiter := new(cliexiterFakes.FakeCliExiter)

				objectUnderTest := _killer{
					nodeProvider: fakeNodeProvider,
					cliExiter:    fakeCliExiter,
				}

				/* act */
				objectUnderTest.Kill()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
