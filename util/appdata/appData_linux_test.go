package appdata

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/user"
)

var _ = Describe("appdata", func() {
	Context("GlobalPath", func() {
		It("should return expected path", func() {
			/* arrange */
			expected := "/var/lib"

			objectUnderTest := New()

			/* act */
			result := objectUnderTest.GlobalPath()

			/* assert */
			Expect(result).To(Equal(expected))
		})
	})
	Context("UserPath", func() {
		It("should return expected path", func() {
			/* arrange */
			currentUser, err := user.Current()
			if nil != err {
				panic(err)
			}
			expected := currentUser.HomeDir

			objectUnderTest := New()

			/* act */
			result := objectUnderTest.PerUserPath()

			/* assert */
			Expect(result).To(Equal(expected))
		})
	})
})
