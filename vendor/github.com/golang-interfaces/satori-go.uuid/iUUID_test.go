package iuuid

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Context("_IUUID", func() {
	Context("New", func() {
		It("should return IUUID", func() {
			/* arrange/act/assert */
			Expect(New()).
				Should(Not(BeNil()))
		})
	})
	Context("NewV4", func() {
		It("should return uuid parseable by uuid.FromString", func() {
			/* arrange */
			objectUnderTest := _IUUID{}

			/* act */
			returnedUUID := objectUnderTest.NewV4()
			_, parseErr := uuid.FromString(returnedUUID.String())

			/* assert */
			Expect(parseErr).To(BeNil())
		})
	})
})
