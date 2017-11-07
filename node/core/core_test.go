package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/containerprovider"
	"github.com/opspec-io/sdk-golang/util/pubsub"
)

var _ = Context("core", func() {
	Context("New", func() {
		It("should return Core", func() {
			/* arrange */
			pubSub := new(pubsub.Fake)
			pubSub.SubscribeStub = func(filter *model.EventFilter, eventChannel chan *model.Event) {
				close(eventChannel)
			}

			/* arrange/act/assert */
			Expect(
				New(
					pubSub,
					new(containerprovider.Fake),
					"dummyFSRootPath",
				),
			).To(Not(BeNil()))
		})
	})
})
