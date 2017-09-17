package opcall

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("OpCall", func() {
	Context("New", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(New("")).To(Not(BeNil()))
		})
	})
})
