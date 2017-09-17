package iio

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("IIO", func() {
	Context("New", func() {
		It("should return IIO", func() {
			/* arrange/act/assert */
			Expect(New()).
				Should(Not(BeNil()))
		})
	})
	Context("Pipe", func() {
		It("should create non nil reader & writer", func() {
			/* arrange */
			objectUnderTest := New()

			/* act */
			actualReader, actualWriter := objectUnderTest.Pipe()
			actualWriter.Close()
			actualReader.Close()

			/* assert */
			Expect(actualReader).NotTo(BeNil())
			Expect(actualWriter).NotTo(BeNil())
		})
	})
	Context("Copy", func() {
		It("should copy src to dst", func() {
			/* arrange */
			objectUnderTest := New()
			expectedBytes := []byte{0, 2, 3, 43}

			src := bytes.NewBuffer(expectedBytes)
			dst := bytes.NewBuffer(nil)

			/* act */
			objectUnderTest.Copy(dst, src)

			/* assert */
			Expect(dst.Bytes()).To(Equal(expectedBytes))
		})
		It("should return expected result", func() {
			/* arrange */
			objectUnderTest := New()
			srcBytes := []byte{0, 2, 3, 43}

			src := bytes.NewBuffer(srcBytes)
			dst := bytes.NewBuffer(make([]byte, src.Len()))

			/* act */
			actualWritten, actualErr := objectUnderTest.Copy(dst, src)

			/* assert */
			Expect(int(actualWritten)).To(Equal(len(srcBytes)))
			Expect(actualErr).To(BeNil())
		})
	})
})
