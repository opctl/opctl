package hostruntime

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("HostPathMap", func() {
	Context("when empty", func() {
		objectUnderTest := newHostPathMap([]string{})

		It("should return original path", func() {
			/* arrange */
			path := "/some/dummy/path"

			/* act */
			actual := objectUnderTest.ToHostPath(path)

			/* assert */
			Expect(actual).To(Equal(path))
		})
	})

	Context("when set", func() {
		objectUnderTest := newHostPathMap([]string{"/host/some/dummy/path:/some/dummy/path"})

		It("should remap path", func() {
			/* arrange */
			path := "/some/dummy/path"
			expected := "/host/some/dummy/path"

			/* act */
			actual := objectUnderTest.ToHostPath(path)

			/* assert */
			Expect(actual).To(Equal(expected))
		})

		It("should remap subpath", func() {
			/* arrange */
			path := "/some/dummy/path/subpath"
			expected := "/host/some/dummy/path/subpath"

			/* act */
			actual := objectUnderTest.ToHostPath(path)

			/* assert */
			Expect(actual).To(Equal(expected))
		})

		It("should not modify not mapped path", func() {
			/* arrange */
			path := "/proc/0"

			/* act */
			actual := objectUnderTest.ToHostPath(path)

			/* assert */
			Expect(actual).To(Equal(path))
		})
	})
})
