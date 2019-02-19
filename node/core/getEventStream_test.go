package core

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core/containerruntime"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

var _ = Context("core", func() {
	Context("GetEventStream", func() {
		It("should call pubSub.Subscribe w/ expected args", func() {
			/* arrange */
			providedCtx := context.TODO()
			providedReq := &model.GetEventStreamReq{
				Filter: model.EventFilter{
					Roots: []string{
						"dummyRootOpID",
					},
				},
			}

			fakePubSub := new(pubsub.Fake)

			objectUnderTest := _core{
				containerRuntime:    new(containerruntime.Fake),
				pubSub:              fakePubSub,
				opCaller:            new(fakeOpCaller),
				dcgNodeRepo:         new(fakeDCGNodeRepo),
				uniqueStringFactory: new(uniquestring.Fake),
			}

			/* act */
			objectUnderTest.GetEventStream(
				providedCtx,
				providedReq,
			)

			/* assert */

			actualCtx,
				actualFilter := fakePubSub.SubscribeArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualFilter).To(Equal(providedReq.Filter))
		})
	})
})
