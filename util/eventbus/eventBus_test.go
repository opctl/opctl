package eventbus

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("eventBus", func() {
	Context("New", func() {
		It("should return EventBus", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
})
