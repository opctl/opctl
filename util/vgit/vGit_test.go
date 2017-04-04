package vgit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("_VGit", func() {
	Context("New", func() {
		It("should return VGIT", func() {
			/* arrange/act/assert */
			Expect(New()).
				Should(Not(BeNil()))
		})
	})
})
