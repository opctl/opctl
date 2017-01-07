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
		It("should invoke bundle.GetEventStream", func() {
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
		It("should call exiter with expected args when bundle.GetEventStream returns error", func() {
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
		It("should call exiter with expected args when event channel closes unexpectedly", func() {
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
