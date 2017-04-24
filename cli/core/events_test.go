package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/consumenodeapi"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("events", func() {
	Context("Execute", func() {
		It("should call pkg.GetEventStream", func() {
			/* arrange */
			fakeCliExiter := new(cliexiter.Fake)

			fakeConsumeNodeApi := new(consumenodeapi.Fake)
			eventChannel := make(chan model.Event)
			close(eventChannel)
			fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

			objectUnderTest := _core{
				consumeNodeApi: fakeConsumeNodeApi,
				cliExiter:      fakeCliExiter,
				nodeProvider:   new(nodeprovider.Fake),
			}

			/* act */
			objectUnderTest.Events()

			/* assert */
			Expect(fakeConsumeNodeApi.GetEventStreamCallCount()).Should(Equal(1))

		})
		Context("pkg.GetEventStream errors", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeCliExiter := new(cliexiter.Fake)
				returnedError := errors.New("dummyError")

				fakeConsumeNodeApi := new(consumenodeapi.Fake)
				fakeConsumeNodeApi.GetEventStreamReturns(nil, returnedError)

				objectUnderTest := _core{
					consumeNodeApi: fakeConsumeNodeApi,
					cliExiter:      fakeCliExiter,
					nodeProvider:   new(nodeprovider.Fake),
				}

				/* act */
				objectUnderTest.Events()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					Should(Equal(cliexiter.ExitReq{Message: returnedError.Error(), Code: 1}))
			})
		})
		Context("pkg.GetEventStream doesn't error", func() {
			Context("channel closes unexpectedly", func() {
				It("should call exiter w/ expected args", func() {
					/* arrange */
					fakeCliExiter := new(cliexiter.Fake)

					fakeConsumeNodeApi := new(consumenodeapi.Fake)
					eventChannel := make(chan model.Event)
					close(eventChannel)
					fakeConsumeNodeApi.GetEventStreamReturns(eventChannel, nil)

					objectUnderTest := _core{
						consumeNodeApi: fakeConsumeNodeApi,
						cliExiter:      fakeCliExiter,
						nodeProvider:   new(nodeprovider.Fake),
					}

					/* act */
					objectUnderTest.Events()

					/* assert */
					Expect(fakeCliExiter.ExitArgsForCall(0)).
						Should(Equal(cliexiter.ExitReq{Message: "Connection to event stream lost", Code: 1}))
				})
			})
		})
	})
})
