package pkg

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("_Pkg", func() {
	Context("Resolve", func() {
		It("should call fs.Stat w/ expected args", func() {
			/* arrange */
			providedBasePath := "dummyBasePath"
			providedPkgRef := &PkgRef{
				FullyQualifiedName: "dummyPkgRef",
				Version:            "0.0.0",
			}

			expectedPath := filepath.Join(
				providedBasePath,
				DotOpspecDirName,
				providedPkgRef.FullyQualifiedName,
				providedPkgRef.Version,
			)

			fakeOS := new(ios.Fake)
			fakeOS.StatReturns(nil, nil)

			objectUnderTest := _Pkg{
				os: fakeOS,
			}

			/* act */
			objectUnderTest.Resolve(providedPkgRef, providedBasePath)

			/* assert */
			Expect(fakeOS.StatArgsForCall(0)).To(Equal(expectedPath))
		})
		Context("fs.Stat doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedBasePath := "dummyBasePath"
				providedPkgRef := &PkgRef{
					FullyQualifiedName: "dummyPkgRef",
					Version:            "0.0.0",
				}

				expectedPath := filepath.Join(
					providedBasePath,
					DotOpspecDirName,
					providedPkgRef.FullyQualifiedName,
					providedPkgRef.Version,
				)
				expectedOk := true

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _Pkg{
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

				objectUnderTest := _Pkg{
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
