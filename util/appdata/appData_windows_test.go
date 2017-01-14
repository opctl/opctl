package appdata

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("appdata", func() {
	Context("GlobalPath", func() {
		It("should return expected path", func() {
			/* arrange */
			expected := os.Getenv("PROGRAMDATA")

			objectUnderTest := New()

			/* act */
			result := objectUnderTest.GlobalPath()

			/* assert */
			Expect(result).To(Equal(expected))
		})
	})
	Context("PerUserPath", func() {
		It("should return expected path", func() {
			/* arrange */
			expected := os.Getenv("LOCALAPPDATA")

			objectUnderTest := New()

			/* act */
			result := objectUnderTest.PerUserPath()

			/* assert */
			Expect(result).To(Equal(expected))
		})
	})
})
