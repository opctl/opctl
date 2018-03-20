package op

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = XContext("Lister", func() {
	Context("NewLister", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewLister()).Should(Not(BeNil()))
		})
	})
	Context("List", func() {
	})
})
