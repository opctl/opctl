package op

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

var _ = Context("Killer", func() {
	Context("Invoke", func() {
		It("should call apiClient.Invoke w/ expected args", func() {
			/* arrange */
			fakeAPIClient := new(client.Fake)

			providedCtx := context.TODO()

			expectedCtx := providedCtx
			expectedReq := model.KillOpReq{
				OpID: "dummyOpID",
			}

			objectUnderTest := _killer{
				apiClient: fakeAPIClient,
			}

			/* act */
			objectUnderTest.Kill(expectedCtx, expectedReq.OpID)

			/* assert */
			actualCtx, actualReq := fakeAPIClient.KillOpArgsForCall(0)
			Expect(actualCtx).To(Equal(expectedCtx))
			Expect(actualReq).To(BeEquivalentTo(expectedReq))
		})
		Context("apiClient.Invoke errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeAPIClient := new(client.Fake)
				expectedError := errors.New("dummyError")
				fakeAPIClient.KillOpReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _killer{
					apiClient: fakeAPIClient,
					cliExiter: fakeCliExiter,
				}

				/* act */
				objectUnderTest.Kill(context.TODO(), "")

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
