package core

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api/client"
)

var _ = Context("opKill", func() {
	Context("Execute", func() {
		It("should call opspecNodeAPIClient.OpKill w/ expected args", func() {
			/* arrange */
			fakeOpspecNodeAPIClient := new(client.Fake)

			providedCtx := context.TODO()

			expectedCtx := providedCtx
			expectedReq := model.KillOpReq{
				OpID: "dummyOpID",
			}

			objectUnderTest := _core{
				opspecNodeAPIClient: fakeOpspecNodeAPIClient,
			}

			/* act */
			objectUnderTest.OpKill(expectedCtx, expectedReq.OpID)

			/* assert */
			actualCtx, actualReq := fakeOpspecNodeAPIClient.KillOpArgsForCall(0)
			Expect(actualCtx).To(Equal(expectedCtx))
			Expect(actualReq).To(BeEquivalentTo(expectedReq))
		})
		Context("opspecNodeAPIClient.OpKill errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeOpspecNodeAPIClient := new(client.Fake)
				expectedError := errors.New("dummyError")
				fakeOpspecNodeAPIClient.KillOpReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					opspecNodeAPIClient: fakeOpspecNodeAPIClient,
					cliExiter:           fakeCliExiter,
				}

				/* act */
				objectUnderTest.OpKill(context.TODO(), "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
