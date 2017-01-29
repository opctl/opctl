package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/pkg/containerengine/engines/fake"
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/opctl/util/pathnormalizer"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Context("core", func() {
	Context("GetEventStream", func() {
		It("should call eventBus.RegisterSubscriber w/ expected args", func() {
			/* arrange */
			providedReq := &model.GetEventStreamReq{
				Filter: &model.EventFilter{
					OpGraphIds: []string{
						"dummyOpGraphId",
					},
				},
			}

			providedEventStream := make(chan model.Event)

			fakeEventBus := new(eventbus.Fake)

			objectUnderTest := _core{
				containerEngine:     new(fake.ContainerEngine),
				eventBus:            fakeEventBus,
				opCaller:            new(fakeOpCaller),
				pathNormalizer:      pathnormalizer.NewPathNormalizer(),
				dcgNodeRepo:         new(fakeDcgNodeRepo),
				uniqueStringFactory: new(uniquestring.Fake),
			}

			/* act */
			objectUnderTest.GetEventStream(providedReq, providedEventStream)

			/* assert */

			// Call happens in go routine; wait 500ms to allow it to occur
			actualFilter,
				actualEventChannel := fakeEventBus.RegisterSubscriberArgsForCall(0)

			Expect(actualFilter).To(Equal(providedReq.Filter))
			Expect(actualEventChannel).To(Equal(providedEventStream))
		})
	})
})
