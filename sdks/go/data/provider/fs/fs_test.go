package fs

import (
	"context"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("_fs", func() {
	Context("TryResolve", func() {
		Context("dataRef is absolute path", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedDataRef := "/dummyFullyQualifiedName"

				fakeOS := new(ios.Fake)
				// error to trigger immediate return
				fakeOS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := _fs{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.TryResolve(
					context.Background(),
					providedDataRef,
				)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(providedDataRef))
			})
			Context("os.Stat errs", func() {
				It("should return err", func() {
					/* arrange */
					expectedErr := errors.New("dummyError")

					fakeOS := new(ios.Fake)
					fakeOS.StatReturns(nil, expectedErr)

					objectUnderTest := _fs{
						os: fakeOS,
					}

					/* act */
					_, actualError := objectUnderTest.TryResolve(
						context.Background(),
						"/dummyDataRef",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedErr))
				})
			})
			Context("os.Stat doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					file, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}

					expectedHandle := newHandle(file.Name())

					fakeOS := new(ios.Fake)
					fakeOS.StatReturns(nil, nil)

					objectUnderTest := _fs{
						os: fakeOS,
					}

					/* act */
					actualHandle, actualError := objectUnderTest.TryResolve(
						context.Background(),
						file.Name(),
					)

					/* assert */
					Expect(actualHandle).To(Equal(expectedHandle))
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("dataRef isn't absolute path", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedDataRef := "dummyDataRef"
				basePath := "dummyBasePath"

				expectedPath := filepath.Join(
					basePath,
					providedDataRef,
				)

				fakeOS := new(ios.Fake)

				objectUnderTest := _fs{
					basePaths: []string{basePath},
					os:        fakeOS,
				}

				/* act */
				objectUnderTest.TryResolve(
					context.Background(),
					providedDataRef,
				)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(expectedPath))
			})
			Context("fs.Stat errors", func() {
				It("should return err", func() {
					/* arrange */
					expectedErr := errors.New("dummyError")

					fakeOS := new(ios.Fake)
					fakeOS.StatReturnsOnCall(0, nil, expectedErr)

					objectUnderTest := _fs{
						basePaths: []string{"dummyBasePath"},
						os:        fakeOS,
					}

					/* act */
					_, actualError := objectUnderTest.TryResolve(
						context.Background(),
						"dummyDataRef",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedErr))
				})
			})
			Context("fs.Stat doesn't err", func() {
				It("should return expected result", func() {
					/* arrange */
					providedDataRef := "dummyDataRef"
					basePath := "dummyBasePath"

					expectedHandle := newHandle(filepath.Join(
						basePath,
						providedDataRef,
					))

					objectUnderTest := _fs{
						basePaths: []string{basePath},
						os:        new(ios.Fake),
					}

					/* act */
					actualHandle, actualError := objectUnderTest.TryResolve(
						context.Background(),
						providedDataRef,
					)

					/* assert */
					Expect(actualHandle).To(Equal(expectedHandle))
					Expect(actualError).To(BeNil())
				})
			})
		})
	})
})
