package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/engineclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Describe("killOp", func() {
	Context("Execute", func() {
		It("should call engineClient.KillOp w/ expected args", func() {
			/* arrange */
			fakeEngineClient := new(engineclient.FakeEngineClient)

			expectedReq := model.KillOpReq{
				OpGraphId: "dummyOpGraphId",
			}

			objectUnderTest := _core{
				engineClient: fakeEngineClient,
			}

			/* act */
			objectUnderTest.KillOp(expectedReq.OpGraphId)

			/* assert */

			Expect(fakeEngineClient.KillOpArgsForCall(0)).Should(BeEquivalentTo(expectedReq))
		})
		Context("engineClient.KillOp errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeEngineClient := new(engineclient.FakeEngineClient)
				expectedError := errors.New("dummyError")
				fakeEngineClient.KillOpReturns(expectedError)

				fakeExiter := new(fakeExiter)

				objectUnderTest := _core{
					engineClient: fakeEngineClient,
					exiter:       fakeExiter,
				}

				/* act */
				objectUnderTest.KillOp("")

				/* assert */
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
