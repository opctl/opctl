package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/virtual-go/vos"
	"path"
)

var _ = Describe("RefResolver", func() {
	Context("newRefResolver()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(newRefResolver(nil)).Should(Not(BeNil()))
		})
	})
	Context("Resolve", func() {
		It("should call os.Getwd", func() {
			/* arrange */
			fakeOS := new(vos.Fake)

			objectUnderTest := _refResolver{
				os: fakeOS,
			}

			/* act */
			objectUnderTest.Resolve("")

			/* assert */
			Expect(fakeOS.GetwdCallCount()).To(Equal(1))
		})
		Context("os.Getwd errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"
				expectedResult := providedPkgRef

				fakeOS := new(vos.Fake)
				fakeOS.GetwdReturns("", errors.New("dummyError"))

				objectUnderTest := _refResolver{
					os: fakeOS,
				}

				/* act */
				actualResult := objectUnderTest.Resolve(expectedResult)

				/* assert */
				Expect(actualResult).To(Equal(providedPkgRef))
			})
		})
		Context("os.Getwd doesn't err", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"

				expectedName := path.Join(DotOpspecDirName, providedPkgRef)

				fakeOS := new(vos.Fake)
				fakeOS.StatReturns(nil, nil)

				objectUnderTest := _refResolver{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.Resolve(providedPkgRef)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(expectedName))
			})
			Context("fs.Stat doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedPkgRef := "dummyPkgRef"
					wd := "dummyWd"

					expectedResult := path.Join(wd, DotOpspecDirName, providedPkgRef)

					fakeOS := new(vos.Fake)
					fakeOS.StatReturns(nil, nil)
					fakeOS.GetwdReturns(wd, nil)

					objectUnderTest := _refResolver{
						os: fakeOS,
					}

					/* act */
					actualResult := objectUnderTest.Resolve(providedPkgRef)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
				})
			})
			Context("fs.Stat errors", func() {
				It("should return expected result", func() {
					/* arrange */
					providedPkgRef := "dummyPkgRef"
					expectedResult := providedPkgRef

					fakeOS := new(vos.Fake)
					fakeOS.StatReturns(nil, errors.New("dummyError"))

					objectUnderTest := _refResolver{
						os: fakeOS,
					}

					/* act */
					actualResult := objectUnderTest.Resolve(providedPkgRef)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
				})
			})
		})
	})
})
