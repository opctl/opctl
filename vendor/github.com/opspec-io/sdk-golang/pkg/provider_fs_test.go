package pkg

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Context("fsProvider", func() {
	Context("TryResolve", func() {
		Context("pkgRef is absolute path", func() {
			It("should call fs.Stat w/ expected args", func() {
				/* arrange */
				providedPkgRef := "/dummyFullyQualifiedName"

				fakeOS := new(ios.Fake)
				// error to trigger immediate return
				fakeOS.StatReturns(nil, errors.New("dummyError"))

				objectUnderTest := fsProvider{
					os: fakeOS,
				}

				/* act */
				objectUnderTest.TryResolve(providedPkgRef)

				/* assert */
				Expect(fakeOS.StatArgsForCall(0)).To(Equal(providedPkgRef))
			})
			Context("os.Stat errs", func() {
				It("should return err", func() {
					/* arrange */
					expectedErr := errors.New("dummyError")

					fakeOS := new(ios.Fake)
					fakeOS.StatReturns(nil, expectedErr)

					objectUnderTest := fsProvider{
						os: fakeOS,
					}

					/* act */
					_, actualError := objectUnderTest.TryResolve(
						"/dummyPkgRef",
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

					expectedHandle := newFSHandle(file.Name())

					fakeOS := new(ios.Fake)
					fakeOS.StatReturns(nil, nil)

					objectUnderTest := fsProvider{
						os: fakeOS,
					}

					/* act */
					actualHandle, actualError := objectUnderTest.TryResolve(file.Name())

					/* assert */
					Expect(actualHandle).To(Equal(expectedHandle))
					Expect(actualError).To(BeNil())
				})
			})
		})
		Context("pkgRef isn't absolute path", func() {
			Context("basePaths not empty", func() {
				Context("pkgRef/.opspec exists", func() {
					It("should call fs.Stat w/ expected args", func() {
						/* arrange */
						providedPkgRef := "dummyPkgRef"
						basePath := "dummyBasePath"

						expectedPath := filepath.Join(
							basePath,
							DotOpspecDirName,
							providedPkgRef,
						)

						fakeOS := new(ios.Fake)
						fakeOS.StatReturns(nil, nil)

						objectUnderTest := fsProvider{
							basePaths: []string{basePath},
							os:        fakeOS,
						}

						/* act */
						objectUnderTest.TryResolve(providedPkgRef)

						/* assert */
						Expect(fakeOS.StatArgsForCall(0)).To(Equal(expectedPath))
					})
					Context("fs.Stat errors", func() {
						It("should return err", func() {
							/* arrange */
							expectedErr := errors.New("dummyError")

							fakeOS := new(ios.Fake)
							fakeOS.StatReturns(nil, expectedErr)

							objectUnderTest := fsProvider{
								basePaths: []string{"dummyBasePath"},
								os:        fakeOS,
							}

							/* act */
							_, actualError := objectUnderTest.TryResolve("dummyPkgRef")

							/* assert */
							Expect(actualError).To(Equal(expectedErr))
						})
					})
					Context("fs.Stat doesn't err", func() {
						It("should return expected result", func() {
							/* arrange */
							providedPkgRef := "dummyPkgRef"
							basePath := "dummyBasePath"

							expectedHandle := newFSHandle(filepath.Join(
								basePath,
								DotOpspecDirName,
								providedPkgRef,
							))

							fakeOS := new(ios.Fake)
							fakeOS.StatReturns(nil, nil)

							objectUnderTest := fsProvider{
								basePaths: []string{basePath},
								os:        fakeOS,
							}

							/* act */
							actualHandle, actualError := objectUnderTest.TryResolve(providedPkgRef)

							/* assert */
							Expect(actualHandle).To(Equal(expectedHandle))
							Expect(actualError).To(BeNil())
						})
					})
				})
				Context("pkgRef/.opspec doesn't exist", func() {
					It("should call fs.Stat w/ expected args", func() {
						/* arrange */
						providedPkgRef := "dummyPkgRef"
						basePath := "dummyBasePath"

						expectedPath := filepath.Join(
							basePath,
							providedPkgRef,
						)

						fakeOS := new(ios.Fake)
						fakeOS.StatReturnsOnCall(0, nil, os.ErrNotExist)

						objectUnderTest := fsProvider{
							basePaths: []string{basePath},
							os:        fakeOS,
						}

						/* act */
						objectUnderTest.TryResolve(providedPkgRef)

						/* assert */
						Expect(fakeOS.StatArgsForCall(1)).To(Equal(expectedPath))
					})
					Context("fs.Stat errors", func() {
						It("should return err", func() {
							/* arrange */
							expectedErr := errors.New("dummyError")

							fakeOS := new(ios.Fake)
							fakeOS.StatReturnsOnCall(0, nil, os.ErrNotExist)
							fakeOS.StatReturnsOnCall(1, nil, expectedErr)

							objectUnderTest := fsProvider{
								basePaths: []string{"dummyBasePath"},
								os:        fakeOS,
							}

							/* act */
							_, actualError := objectUnderTest.TryResolve("dummyPkgRef")

							/* assert */
							Expect(actualError).To(Equal(expectedErr))
						})
					})
					Context("fs.Stat doesn't err", func() {
						It("should return expected result", func() {
							/* arrange */
							providedPkgRef := "dummyPkgRef"
							basePath := "dummyBasePath"

							expectedHandle := newFSHandle(filepath.Join(
								basePath,
								providedPkgRef,
							))

							fakeOS := new(ios.Fake)
							fakeOS.StatReturnsOnCall(0, nil, os.ErrNotExist)

							objectUnderTest := fsProvider{
								basePaths: []string{basePath},
								os:        fakeOS,
							}

							/* act */
							actualHandle, actualError := objectUnderTest.TryResolve(providedPkgRef)

							/* assert */
							Expect(actualHandle).To(Equal(expectedHandle))
							Expect(actualError).To(BeNil())
						})
					})
				})
			})
		})
	})
})
