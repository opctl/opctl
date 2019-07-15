package direntry

import (
	"fmt"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).Should(Not(BeNil()))
		})
	})
	Context("Interpreter", func() {
		Context("ref doesn't start w/ '/'", func() {
			It("should return expected result", func() {
				/* arrange */
				providedRef := "dummyRef"

				objectUnderTest := _interpreter{}
				expectedErr := fmt.Errorf("unable to interpret '%v'; expected '/'", providedRef)

				/* act */
				_, _, actualErr := objectUnderTest.Interpret(
					providedRef,
					nil,
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("ref is file", func() {
			It("should return expected result", func() {
				/* arrange */
				value := "/dummyScopeValue"

				providedRef := "/providedRef"

				expectedPath := filepath.Join(value, providedRef)

				// no good way to fake fileinfo
				tmpFile, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}

				tmpFileStat, err := tmpFile.Stat()
				if nil != err {
					panic(err)
				}

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(tmpFileStat, nil)

				objectUnderTest := _interpreter{
					os: fakeOS,
				}

				/* act */

				actualRefRemainder, actualValue, actualErr := objectUnderTest.Interpret(
					providedRef,
					&model.Value{Dir: &value},
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
				value := "/dummyScopeValue"

				providedRef := "/providedRef"

				expectedPath := filepath.Join(value, providedRef)

				// no good way to fake fileinfo
				tmpDir, err := ioutil.TempDir("", "")
				if nil != err {
					panic(err)
				}

				tmpDirStat, err := os.Stat(tmpDir)
				if nil != err {
					panic(err)
				}

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(tmpDirStat, nil)

				objectUnderTest := _interpreter{
					os: fakeOS,
				}

				/* act */

				actualRefRemainder, actualValue, actualErr := objectUnderTest.Interpret(
					providedRef,
					&model.Value{Dir: &value},
				)

				/* assert */
				Expect(actualRefRemainder).To(BeEmpty())
				Expect(*actualValue).To(Equal(model.Value{Dir: &expectedPath}))
				Expect(actualErr).To(BeNil())

			})
		})
	})
})
