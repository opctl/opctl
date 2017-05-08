package pkg

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path"
)

var _ = Describe("Pkg", func() {
	Context("Resolve", func() {
		It("should call fs.Stat w/ expected args", func() {
			/* arrange */
			providedBasePath := "dummyBasePath"
			providedPkgRef := "dummyPkgRef"

			expectedName := path.Join(providedBasePath, DotOpspecDirName, providedPkgRef)

			fakeOS := new(ios.Fake)
			fakeOS.StatReturns(nil, nil)

			objectUnderTest := pkg{
				os: fakeOS,
			}

			/* act */
			objectUnderTest.Resolve(providedBasePath, providedPkgRef)

			/* assert */
			Expect(fakeOS.StatArgsForCall(0)).To(Equal(expectedName))
		})
		Context("fs.Stat doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedBasePath := "dummyBasePath"
				providedPkgRef := "dummyPkgRef"

				expectedPath := path.Join(providedBasePath, DotOpspecDirName, providedPkgRef)
				expectedOk := true

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := pkg{
					os: fakeOS,
				}

				/* act */
				actualPath, actualOk := objectUnderTest.Resolve(providedBasePath, providedPkgRef)

				/* assert */
				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualOk).To(Equal(expectedOk))
			})
		})
		Context("fs.Stat errors", func() {
			It("should return expected result", func() {
				/* arrange */
				providedBasePath := "dummyBasePath"
				providedPkgRef := "dummyPkgRef"

				expectedPath := ""
				expectedOk := false

				fakeOS := new(ios.Fake)
				fakeOS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := pkg{
					os: fakeOS,
				}

				/* act */
				actualPath, actualOk := objectUnderTest.Resolve(providedBasePath, providedPkgRef)

				/* assert */
				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualOk).To(Equal(expectedOk))
			})
		})
	})
})
