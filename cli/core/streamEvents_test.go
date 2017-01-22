package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/engineclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Describe("streamEvents", func() {
	Context("Execute", func() {
		It("should call bundle.GetEventStream", func() {
			/* arrange */
			fakeExiter := new(fakeExiter)

			fakeEngineClient := new(engineclient.FakeEngineClient)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				engineClient: fakeEngineClient,
				exiter:       fakeExiter,
			}

			/* act */
			objectUnderTest.StreamEvents()

			/* assert */
			Expect(fakeEngineClient.GetEventStreamCallCount()).Should(Equal(1))

		})
		Context("bundle.GetEventStream errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeExiter := new(fakeExiter)
				returnedError := errors.New("dummyError")

				fakeEngineClient := new(engineclient.FakeEngineClient)
				fakeEngineClient.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _core{
					engineClient: fakeEngineClient,
					exiter:       fakeExiter,
				}

				/* act */
				objectUnderTest.StreamEvents()

				/* assert */
				Expect(fakeExiter.ExitArgsForCall(0)).
					Should(Equal(ExitReq{Message: returnedError.Error(), Code: 1}))
			})
		})
		Context("bundle.GetEventStream doesn't error", func() {
			Context("channel closes unexpectedly", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeExiter := new(fakeExiter)

					fakeEngineClient := new(engineclient.FakeEngineClient)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						engineClient: fakeEngineClient,
						exiter:       fakeExiter,
					}

					/* act */
					objectUnderTest.StreamEvents()

					/* assert */
					Expect(fakeExiter.ExitArgsForCall(0)).
						Should(Equal(ExitReq{Message: "Event channel closed unexpectedly", Code: 1}))
				})
			})
		})
	})
})
