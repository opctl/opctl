package direntry

import (
	"fmt"
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
			dirData, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

			expectedPath := filepath.Join(dirData, providedRef)
			err = os.WriteFile(expectedPath, []byte(""), 0777)
			if err != nil {
				panic(err)
			}

			/* act */

			actualRefRemainder, actualValue, actualErr := Interpret(
				providedRef,
				&ipld.Node{Dir: &dirData},
				nil,
			)

			/* assert */
			Expect(actualRefRemainder).To(BeEmpty())
			Expect(*actualValue).To(Equal(ipld.Node{File: &expectedPath}))
			Expect(actualErr).To(BeNil())

		})
	})
	Context("ref is dir", func() {
		It("should return expected result", func() {
			/* arrange */
			providedRef := "/providedRef"
			dirData, err := os.MkdirTemp("", "")
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
				&ipld.Node{Dir: &dirData},
				nil,
			)

			/* assert */
			Expect(actualRefRemainder).To(BeEmpty())
			Expect(*actualValue).To(Equal(ipld.Node{Dir: &expectedPath}))
			Expect(actualErr).To(BeNil())

		})
	})
})
