package image

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("Image", func() {
	Context("New", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(New()).To(Not(BeNil()))
		})
	})
})
