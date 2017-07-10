package pkg

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"path/filepath"
)

var _ = Context("Resolver", func() {
	Context("Resolve", func() {
		Context("pkgRef is absolute path", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedPkgRef := "/dummyFullyQualifiedName"

				fakeOS := new(ios.Fake)
				// error to trigger immediate return
				fakeOS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.Resolve(providedPkgRef, nil)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(providedPkgRef))
			})
			It("should return expected result", func() {
				/* arrange */
				file, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}

				expectedHandle := newLocalHandle(file.Name())

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				actualHandle, actualError := objectUnderTest.Resolve(file.Name(), nil)

				/* assert */
				Expect(actualHandle).To(Equal(expectedHandle))
				Expect(actualError).To(BeNil())
			})
		})
		Context("pkgRef isn't absolute path", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"
				providedOpts := &ResolveOpts{
					BasePath: "dummyBasePath",
				}

				expectedPath := filepath.Join(
					providedOpts.BasePath,
					DotOpspecDirName,
					providedPkgRef,
				)

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.Resolve(providedPkgRef, providedOpts)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(expectedPath))
			})
		})
		Context("fs.Stat errors", func() {
			It("should return err", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, expectedErr)

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				_, actualError := objectUnderTest.Resolve(
					"dummyPkgRef",
					nil,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedErr))
			})
		})
		Context("fs.Stat doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"
				providedOpts := &ResolveOpts{
					BasePath: "dummyBasePath",
				}

				expectedHandle := newLocalHandle(filepath.Join(
					providedOpts.BasePath,
					DotOpspecDirName,
					providedPkgRef,
				))

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				actualHandle, actualError := objectUnderTest.Resolve(
					providedPkgRef,
					providedOpts,
				)

				/* assert */
				Expect(actualHandle).To(Equal(expectedHandle))
				Expect(actualError).To(BeNil())
			})
		})
	})
})
