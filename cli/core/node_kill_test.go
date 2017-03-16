package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/nodeprovider"
	"github.com/opspec-io/opctl/util/cliexiter"
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
			Expect(fakeNodeProvider.KillNodeIfExistsArgsForCall(0)).Should(BeEquivalentTo(""))
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
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
