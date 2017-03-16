package validate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Param", func() {
	objectUnderTest := New()
	Context("invoked w/ nil param", func() {
		It("should panic", func() {
			/* arrange/act/assert */
			Expect(
				func() {
					objectUnderTest.Param(&model.Data{}, nil)
				},
			).To(Panic())
		})
	})
})
