package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/containerprovider"
	"github.com/opspec-io/sdk-golang/util/pubsub"
)

var _ = Context("core", func() {
	Context("New", func() {
		It("should return Core", func() {
			/* arrange/act/assert */
			Expect(
				New(
					new(pubsub.Fake),
					new(containerprovider.Fake),
					"dummyFSRootPath",
				),
			).To(Not(BeNil()))
		})
	})
})
