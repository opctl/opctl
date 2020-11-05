package item

import (
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("ParseIndex", func() {
	Context("index doesn't parse to integer", func() {
		It("should return expected result", func() {
			/* arrange */
			providedIdentifier := "blah"

			_, parseIntErr := strconv.ParseInt(providedIdentifier, 10, 64)

			/* act */
			_, actualErr := ParseIndex(
				providedIdentifier,
				[]interface{}{},
			)

			/* assert */
			Expect(actualErr).To(Equal(parseIntErr))
		})
	})
	Context("index parses to integer", func() {
		Context("index negative", func() {
			Context("index within range of array", func() {
				It("should return expected result", func() {
					/* arrange */
					arrayItemIndex := -1
					providedArray := []interface{}{"hello"}

					expectedIndex := int64(arrayItemIndex + len(providedArray))

					/* act */
					actualIndex, actualErr := ParseIndex(
						fmt.Sprintf("%v", arrayItemIndex),
						providedArray,
					)

					/* assert */
					Expect(actualIndex).To(Equal(expectedIndex))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("index outside range of array", func() {
				It("should return expected result", func() {
					/* arrange */
					arrayItemIndex := -1
					providedArray := []interface{}{}

					expectedErr := fmt.Errorf("array index %v out of range 0-%v", arrayItemIndex, len(providedArray)-1)

					/* act */
					_, actualErr := ParseIndex(
						fmt.Sprintf("%v", arrayItemIndex),
						providedArray,
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
		})
		Context("index positive", func() {
			Context("index within range of array", func() {
				It("should return expected result", func() {
					/* arrange */
					arrayItemIndex := int64(0)
					providedArray := []interface{}{"hello"}

					/* act */
					actualIndex, actualErr := ParseIndex(
						fmt.Sprintf("%v", arrayItemIndex),
						providedArray,
					)

					/* assert */
					Expect(actualIndex).To(Equal(arrayItemIndex))
					Expect(actualErr).To(BeNil())
				})
			})
			Context("index outside range of array", func() {
				It("should return expected result", func() {
					/* arrange */
					arrayItemIndex := -1
					providedArray := []interface{}{}

					expectedErr := fmt.Errorf("array index %v out of range 0-%v", arrayItemIndex, len(providedArray)-1)

					/* act */
					_, actualErr := ParseIndex(
						fmt.Sprintf("%v", arrayItemIndex),
						providedArray,
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
		})
	})
})
