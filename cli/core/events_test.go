package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/engineclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("streamEvents", func() {
	Context("Execute", func() {
		It("should call bundle.GetEventStream", func() {
			/* arrange */
			fakeCliExiter := new(cliexiter.Fake)

			fakeEngineClient := new(engineclient.Fake)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				engineClient: fakeEngineClient,
				cliExiter:    fakeCliExiter,
			}

			/* act */
			objectUnderTest.StreamEvents()

			/* assert */
			Expect(fakeEngineClient.GetEventStreamCallCount()).Should(Equal(1))

		})
		Context("bundle.GetEventStream errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiter.Fake)
				returnedError := errors.New("dummyError")

				fakeEngineClient := new(engineclient.Fake)
				fakeEngineClient.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _core{
					engineClient: fakeEngineClient,
					cliExiter:    fakeCliExiter,
				}

				/* act */
				objectUnderTest.StreamEvents()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
			})
		})
		Context("bundle.GetEventStream doesn't error", func() {
			Context("channel closes unexpectedly", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)

					fakeEngineClient := new(engineclient.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeEngineClient.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						engineClient: fakeEngineClient,
						cliExiter:    fakeCliExiter,
					}

					/* act */
					objectUnderTest.StreamEvents()

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1}))
				})
			})
		})
	})
})
