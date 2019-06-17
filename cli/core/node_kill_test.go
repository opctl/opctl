package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/nodeprovider"
	"github.com/opctl/opctl/cli/util/cliexiter"
)

var _ = Context("nodeKill", func() {
	Context("Execute", func() {
		It("should call nodeProvider.NodeKill w/ expected args", func() {
			/* arrange */
			fakeNodeProvider := new(nodeprovider.Fake)

			objectUnderTest := _core{
				nodeProvider: fakeNodeProvider,
			}

			/* act */
			objectUnderTest.NodeKill()

			/* assert */
			Expect(fakeNodeProvider.KillNodeIfExistsArgsForCall(0)).To(BeEquivalentTo(""))
		})
		Context("nodeProvider.NodeKill errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeNodeProvider := new(nodeprovider.Fake)
				expectedError := errors.New("dummyError")
				fakeNodeProvider.KillNodeIfExistsReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					nodeProvider: fakeNodeProvider,
					cliExiter:    fakeCliExiter,
				}

				/* act */
				objectUnderTest.NodeKill()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
