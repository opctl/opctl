package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/consumenodeapi"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("killOp", func() {
	Context("Execute", func() {
		It("should call consumeNodeApi.KillOp w/ expected args", func() {
			/* arrange */
			fakeConsumeNodeApi := new(consumenodeapi.Fake)

			expectedReq := model.KillOpReq{
				OpId: "dummyOpId",
			}

			objectUnderTest := _core{
				consumeNodeApi: fakeConsumeNodeApi,
			}

			/* act */
			objectUnderTest.KillOp(expectedReq.OpId)

			/* assert */

			Expect(fakeConsumeNodeApi.KillOpArgsForCall(0)).Should(BeEquivalentTo(expectedReq))
		})
		Context("consumeNodeApi.KillOp errors", func() {
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
				objectUnderTest.KillOp("")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
