package urlpath

import (
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("NextSegment", func() {
	Context("Path has no segments", func() {
		Context("no escaped characters", func() {
			It("should return expected result", func() {
				/* arrange */
				providedURL, err := url.Parse("")
				if err != nil {
					panic(err.Error)
				}

				/* act */
				actualPath, err := NextSegment(providedURL)

				/* assert */
				Expect(actualPath).To(Equal(""))
			})
		})
	})
	Context("Path has single segment", func() {
		Context("no escaped characters", func() {
			It("should return expected result", func() {
				/* arrange */
				providedSegment := "1"

				providedURL, err := url.Parse(providedSegment)
				if err != nil {
					panic(err.Error)
				}

				/* act */
				actualPath, err := NextSegment(providedURL)

				/* assert */
				Expect(err).To(BeNil())
				Expect(providedURL.Path).To(BeEmpty())
				Expect(providedURL.RawPath).To(BeEmpty())
				Expect(actualPath).To(Equal(providedSegment))
			})
		})
		Context("invalid escape characters", func() {
			It("should return expected result", func() {
				/* arrange */
				providedSegment := "%%%"

				providedURL := &url.URL{
					Path:    providedSegment,
					RawPath: providedSegment,
				}

				/* act */
				actualPath, err := NextSegment(providedURL)

				/* assert */
				Expect(err).To(BeNil())
				Expect(providedURL.Path).To(BeEmpty())
				Expect(providedURL.RawPath).To(BeEmpty())
				Expect(actualPath).To(Equal(providedSegment))
			})
		})
	})
	Context("Path has two segments", func() {
		Context("no escaped characters", func() {
			It("should return expected result", func() {
				/* arrange */
				providedSegments := []string{
					"1",
					"2",
				}

				providedURL, err := url.Parse(strings.Join(providedSegments, "/"))
				if err != nil {
					panic(err.Error)
				}

				/* act */
				actualPath, err := NextSegment(providedURL)

				/* assert */
				Expect(err).To(BeNil())
				Expect(providedURL.Path).To(Equal(providedSegments[1]))
				Expect(providedURL.RawPath).To(BeEmpty())
				Expect(actualPath).To(Equal(providedSegments[0]))
			})
		})
	})
	Context("Path has multiple segments", func() {
		Context("no escaped characters", func() {
			It("should return expected result", func() {
				/* arrange */
				providedSegments := []string{
					"1",
					"2",
					"3",
				}

				providedURL, err := url.Parse(strings.Join(providedSegments, "/"))
				if err != nil {
					panic(err.Error)
				}

				/* act */
				actualPath, err := NextSegment(providedURL)

				/* assert */
				Expect(err).To(BeNil())
				Expect(providedURL.Path).To(Equal(strings.Join(providedSegments[1:], "/")))
				Expect(providedURL.RawPath).To(BeEmpty())
				Expect(actualPath).To(Equal(providedSegments[0]))
			})
		})
	})
})
