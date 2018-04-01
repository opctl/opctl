package direntry

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
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
			It("should call ioutil.ReadFile w/ expected args", func() {
				/* arrange */
				value := "/dummyScopeValue"

				providedRef := "/providedRef"

				expectedPath := filepath.Join(value, providedRef)

				fakeIoUtil := new(iioutil.Fake)
				// err to trigger immediate return
				fakeIoUtil.ReadFileReturns([]byte{}, errors.New("dummyError"))

				objectUnderTest := _deReferencer{
					ioutil: fakeIoUtil,
				}

				/* act */

				objectUnderTest.DeReference(
					providedRef,
					&model.Value{Dir: &value},
				)

				/* assert */
				actualPath := fakeIoUtil.ReadFileArgsForCall(0)

				Expect(actualPath).To(Equal(expectedPath))
			})
		})
	})
})
