package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/util/pubsub"
)

var _ = Context("core", func() {
	Context("New", func() {
		It("should return Core", func() {
			/* arrange/act/assert */
			Expect(
				New(
					new(pubsub.Fake),
					new(containerruntime.Fake),
					"dummyFSRootPath",
				),
			).To(Not(BeNil()))
		})
	})
})
