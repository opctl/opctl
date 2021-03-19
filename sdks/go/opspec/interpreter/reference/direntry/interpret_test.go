package direntry

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("ref doesn't start w/ '/'", func() {
		It("should return expected result", func() {
			/* arrange */
			providedRef := "dummyRef"

			expectedErr := fmt.Errorf("unable to interpret '%v' as dir entry ref: expected '/'", providedRef)

			/* act */
			_, _, actualErr := Interpret(
				providedRef,
				nil,
				nil,
			)

			/* assert */
			Expect(actualErr).To(Equal(expectedErr))
		})
	})
	Context("ref is file", func() {
		It("should return expected result", func() {
			/* arrange */
			providedRef := "/providedRef"
			dirData, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			expectedPath := filepath.Join(dirData, providedRef)
			err = ioutil.WriteFile(expectedPath, []byte(""), 0777)
			if err != nil {
				panic(err)
			}

			/* act */

			actualRefRemainder, actualValue, actualErr := Interpret(
				providedRef,
				&model.Value{Dir: &dirData},
				nil,
			)

			/* assert */
			Expect(actualRefRemainder).To(BeEmpty())
			Expect(*actualValue).To(Equal(model.Value{File: &expectedPath}))
			Expect(actualErr).To(BeNil())

		})
	})
	Context("ref is dir", func() {
		It("should return expected result", func() {
			/* arrange */
			providedRef := "/providedRef"
			dirData, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			expectedPath := filepath.Join(dirData, providedRef)
			err = os.Mkdir(expectedPath, 0777)
			if err != nil {
				panic(err)
			}

			/* act */

			actualRefRemainder, actualValue, actualErr := Interpret(
				providedRef,
				&model.Value{Dir: &dirData},
				nil,
			)

			/* assert */
			Expect(actualRefRemainder).To(BeEmpty())
			Expect(*actualValue).To(Equal(model.Value{Dir: &expectedPath}))
			Expect(actualErr).To(BeNil())

		})
	})
})
