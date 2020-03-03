package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
	. "github.com/opctl/opctl/sdks/go/pubsub/fakes"
)

var _ = Context("core", func() {
	Context("New", func() {
		It("should return Core", func() {
			/* arrange/act/assert */
			Expect(
				New(
					new(FakePubSub),
					new(FakeContainerRuntime),
					"dummyFSRootPath",
				),
			).To(Not(BeNil()))
		})
	})
})
