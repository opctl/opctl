package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/apiclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("killOp", func() {
	Context("Execute", func() {
		It("should call apiClient.KillOp w/ expected args", func() {
			/* arrange */
			fakeApiClient := new(apiclient.Fake)

			expectedReq := model.KillOpReq{
				OpGraphId: "dummyOpGraphId",
			}

			objectUnderTest := _core{
				apiClient: fakeApiClient,
			}

			/* act */
			objectUnderTest.KillOp(expectedReq.OpGraphId)

			/* assert */

			Expect(fakeApiClient.KillOpArgsForCall(0)).Should(BeEquivalentTo(expectedReq))
		})
		Context("apiClient.KillOp errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeApiClient := new(apiclient.Fake)
				expectedError := errors.New("dummyError")
				fakeApiClient.KillOpReturns(expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _core{
					apiClient: fakeApiClient,
					cliExiter: fakeCliExiter,
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
