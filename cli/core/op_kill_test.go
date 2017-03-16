package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/consumenodeapi"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("opKill", func() {
	Context("Execute", func() {
		It("should call consumeNodeApi.OpKill w/ expected args", func() {
			/* arrange */
			fakeConsumeNodeApi := new(consumenodeapi.Fake)

			expectedReq := model.KillOpReq{
				OpId: "dummyOpId",
			}

			objectUnderTest := _core{
				consumeNodeApi: fakeConsumeNodeApi,
			}

			/* act */
			objectUnderTest.OpKill(expectedReq.OpId)

			/* assert */

			Expect(fakeConsumeNodeApi.KillOpArgsForCall(0)).Should(BeEquivalentTo(expectedReq))
		})
		Context("consumeNodeApi.OpKill errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeConsumeNodeApi := new(consumenodeapi.Fake)
				expectedError := errors.New("dummyError")
				fakeConsumeNodeApi.KillOpReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					consumeNodeApi: fakeConsumeNodeApi,
					cliExiter:      fakeCliExiter,
				}

				/* act */
				objectUnderTest.OpKill("")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
