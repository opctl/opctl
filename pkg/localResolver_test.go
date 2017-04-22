package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/virtual-go/vos"
	"path"
)

var _ = Describe("LocalResolver", func() {
	Context("newLocalResolver()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(newLocalResolver(nil)).Should(Not(BeNil()))
		})
	})
	Context("Resolve", func() {
		It("should call fs.Stat w/ expected args", func() {
			/* arrange */
			providedBasePath := "dummyBasePath"
			providedPkgRef := "dummyPkgRef"

			expectedName := path.Join(providedBasePath, DotOpspecDirName, providedPkgRef)

			fakeOS := new(vos.Fake)
			fakeOS.StatReturns(nil, nil)

			objectUnderTest := _localResolver{
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

				fakeOS := new(vos.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _localResolver{
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

				fakeOS := new(vos.Fake)
				fakeOS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := _localResolver{
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
