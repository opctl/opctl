package pkg

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"path/filepath"
)

var _ = Context("resolver", func() {
	Context("Resolve", func() {
		Context("pkgRef is absolute path", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedPkgRef := &PkgRef{
					FullyQualifiedName: "/dummyFullyQualifiedName",
				}

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.Resolve(providedPkgRef)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(providedPkgRef.FullyQualifiedName))
			})
			It("should return expected result", func() {
				/* arrange */
				providedPkgRef := &PkgRef{
					FullyQualifiedName: "/dummyFullyQualifiedName",
				}

				file, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}

				expectedFileInfo, err := file.Stat()
				if nil != err {
					panic(err)
				}

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(expectedFileInfo, nil)

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				actualPkgPath, actualIsResolved := objectUnderTest.Resolve(providedPkgRef)

				/* assert */
				Expect(actualPkgPath).To(Equal(providedPkgRef.FullyQualifiedName))
				Expect(actualIsResolved).To(Equal(true))
			})
		})
		Context("pkgRef isn't absolute path", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedBasePath := "dummyBasePath"
				providedPkgRef := &PkgRef{
					FullyQualifiedName: "dummyPkgRef",
					Version:            "0.0.0",
				}

				expectedPath := providedPkgRef.ToPath(
					filepath.Join(
						providedBasePath,
						DotOpspecDirName,
					),
				)

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.Resolve(providedPkgRef, providedBasePath)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(expectedPath))
			})
		})
		Context("fs.Stat doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedBasePath := "dummyBasePath"
				providedPkgRef := &PkgRef{
					FullyQualifiedName: "dummyPkgRef",
					Version:            "0.0.0",
				}

				expectedPath := providedPkgRef.ToPath(
					filepath.Join(
						providedBasePath,
						DotOpspecDirName,
					),
				)
				expectedOk := true

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				actualPath, actualOk := objectUnderTest.Resolve(providedPkgRef, providedBasePath)

				/* assert */
				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualOk).To(Equal(expectedOk))
			})
		})
		Context("fs.Stat errors", func() {
			It("should return expected result", func() {
				/* arrange */
				providedBasePath := "dummyBasePath"
				providedPkgRef := &PkgRef{
					FullyQualifiedName: "dummyPkgRef",
					Version:            "0.0.0",
				}

				expectedPath := ""
				expectedOk := false

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := _resolver{
					os: fakeOS,
				}

				/* act */
				actualPath, actualOk := objectUnderTest.Resolve(providedPkgRef, providedBasePath)

				/* assert */
				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualOk).To(Equal(expectedOk))
			})
		})
	})
})
