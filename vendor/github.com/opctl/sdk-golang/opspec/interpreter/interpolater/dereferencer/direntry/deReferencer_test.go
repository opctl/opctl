package direntry

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
	"path/filepath"
)

var _ = Context("DeReferencer", func() {
	Context("NewDeReferencer", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewDeReferencer()).Should(Not(BeNil()))
		})
	})
	Context("DeReferencer", func() {
		Context("ref doesn't start w/ '/'", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				objectUnderTest := _deReferencer{}
				expectedErr := fmt.Errorf("unable to deReference '%v'; expected '/'", providedRef)

				/* act */
				_, _, actualErr := objectUnderTest.DeReference(
					providedRef,
					nil,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("ref is scope file path ref", func() {
			It("should return expected result", func() {
				/* arrange */
				value := "/dummyScopeValue"

				providedRef := "/providedRef"

				expectedPath := filepath.Join(value, providedRef)

				objectUnderTest := _deReferencer{}

				/* act */

				actualRefRemainder, actualValue, actualErr := objectUnderTest.DeReference(
					providedRef,
					&model.Value{Dir: &value},
				)

				/* assert */
				Expect(actualRefRemainder).To(BeEmpty())
				Expect(*actualValue).To(Equal(model.Value{File: &expectedPath}))
				Expect(actualErr).To(BeNil())

			})
		})
	})
})
