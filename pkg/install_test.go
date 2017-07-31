package pkg

import (
	"context"
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Context("pkg", func() {

	Context("Install", func() {

		It("should call os.MkdirAll w/ expected args", func() {
			/* arrange */
			providedPath := "dummyPath"

			fakeOS := new(ios.Fake)
			// error to trigger immediate return
			fakeOS.MkdirAllReturns(errors.New("dummyError"))

			objectUnderTest := _Pkg{
				os: fakeOS,
			}

			/* act */
			objectUnderTest.Install(nil, providedPath, nil)

			/* assert */
			actualPath, actualPerm := fakeOS.MkdirAllArgsForCall(0)
			Expect(actualPath).To(Equal(providedPath))
			Expect(actualPerm).To(Equal(os.FileMode(0777)))
		})
		Context("os.MkdirAll errs", func() {
			It("should return error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeOS := new(ios.Fake)
				fakeOS.MkdirAllReturns(expectedError)

				objectUnderTest := _Pkg{
					os: fakeOS,
				}

				/* act */
				actualError := objectUnderTest.Install(nil, "", nil)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("os.MkdirAll doesn't err", func() {
			It("should call handle.ListContents w/ expected args", func() {
				/* arrange */
				providedCtx := context.TODO()

				fakeHandle := new(FakeHandle)

				objectUnderTest := _Pkg{
					os: new(ios.Fake),
				}

				/* act */
				objectUnderTest.Install(providedCtx, "", fakeHandle)

				/* assert */
				Expect(fakeHandle.ListContentsArgsForCall(0)).To(Equal(providedCtx))
			})
			Context("handle.ListContents errs", func() {
				It("should return error", func() {
					/* arrange */
					expectedError := errors.New("dummyError")

					fakeHandle := new(FakeHandle)
					fakeHandle.ListContentsReturns(nil, expectedError)

					objectUnderTest := _Pkg{
						os: new(ios.Fake),
					}

					/* act */
					actualError := objectUnderTest.Install(nil, "", fakeHandle)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("handle.ListContents doesn't err", func() {
				It("should call handle.GetContent w/ expectet args", func() {
					/* arrange */
					providedCtx := context.TODO()

					fakeHandle := new(FakeHandle)
					contentsList := []*model.PkgContent{
						{
							Path: "pkgContent1Path",
						},
					}

					fakeHandle.ListContentsReturns(
						contentsList,
						nil,
					)

					// error to trigger immediate return
					fakeHandle.GetContentReturns(nil, errors.New("dummyError"))

					objectUnderTest := _Pkg{
						os: new(ios.Fake),
					}

					/* act */
					objectUnderTest.Install(providedCtx, "", fakeHandle)

					/* assert */
					actualContext,
						actualPath := fakeHandle.GetContentArgsForCall(0)

					Expect(actualContext).To(Equal(providedCtx))
					Expect(actualPath).To(Equal(contentsList[0].Path))
				})
				Context("handle.GetContent errs", func() {
					It("should return error", func() {
						/* arrange */
						expectedError := errors.New("dummyError")

						fakeHandle := new(FakeHandle)
						fakeHandle.ListContentsReturns([]*model.PkgContent{{}}, expectedError)

						fakeHandle.GetContentReturns(nil, expectedError)

						objectUnderTest := _Pkg{
							os: new(ios.Fake),
						}

						/* act */
						actualError := objectUnderTest.Install(nil, "", fakeHandle)

						/* assert */
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("handle.GetContent doesn't err", func() {
					It("should call os.MkdirAll w/ expected args", func() {
						/* arrange */
						providedPath := "dummyPath"

						fakeHandle := new(FakeHandle)
						contentsList := []*model.PkgContent{
							{
								Path: "pkgContent1Path",
							},
						}

						fakeHandle.ListContentsReturns(
							contentsList,
							nil,
						)

						fakeOS := new(ios.Fake)
						// error to trigger immediate return
						fakeOS.MkdirAllReturnsOnCall(1, errors.New("dummyError"))

						objectUnderTest := _Pkg{
							os: fakeOS,
						}

						/* act */
						objectUnderTest.Install(nil, providedPath, fakeHandle)

						/* assert */
						actualPath, actualPerm := fakeOS.MkdirAllArgsForCall(1)

						Expect(actualPath).To(Equal(
							filepath.Dir(
								filepath.Join(providedPath, contentsList[0].Path),
							),
						))

						Expect(actualPerm).To(Equal(os.FileMode(0777)))
					})
				})
			})
			Context("os.MkdirAll errs", func() {
				It("should return error", func() {
					/* arrange */
					expectedError := errors.New("dummyError")

					fakeHandle := new(FakeHandle)
					fakeHandle.ListContentsReturns([]*model.PkgContent{{}}, nil)

					fakeOS := new(ios.Fake)
					fakeOS.MkdirAllReturnsOnCall(1, expectedError)

					objectUnderTest := _Pkg{
						os: fakeOS,
					}

					/* act */
					actualError := objectUnderTest.Install(nil, "", fakeHandle)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("os.MkdirAll doesn't err", func() {
				It("should call os.Create w/ expected args", func() {
					/* arrange */
					providedPath := "dummyPath"

					fakeHandle := new(FakeHandle)
					contentsList := []*model.PkgContent{
						{
							Path: "pkgContent1Path",
						},
					}

					fakeHandle.ListContentsReturns(
						contentsList,
						nil,
					)

					fakeOS := new(ios.Fake)
					// error to trigger immediate return
					fakeOS.CreateReturns(nil, errors.New("dummyError"))

					objectUnderTest := _Pkg{
						os: fakeOS,
					}

					/* act */
					objectUnderTest.Install(nil, providedPath, fakeHandle)

					/* assert */
					actualPath := fakeOS.CreateArgsForCall(0)

					Expect(actualPath).To(Equal(filepath.Join(providedPath, contentsList[0].Path)))
				})
			})
			Context("os.Create errs", func() {
				It("should return error", func() {
					/* arrange */
					expectedError := errors.New("dummyError")

					fakeHandle := new(FakeHandle)
					fakeHandle.ListContentsReturns([]*model.PkgContent{{}}, nil)

					fakeOS := new(ios.Fake)
					fakeOS.CreateReturns(nil, expectedError)

					objectUnderTest := _Pkg{
						os: fakeOS,
					}

					/* act */
					actualError := objectUnderTest.Install(nil, "", fakeHandle)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("os.Create doesn't err", func() {
				It("should copy content", func() {
					/* arrange */
					fakeHandle := new(FakeHandle)
					fakeHandle.ListContentsReturns([]*model.PkgContent{{}}, nil)

					// create tmpfile to use as src
					contentSrc, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}
					defer contentSrc.Close()

					expectedContent := []byte("dummyString")
					err = ioutil.WriteFile(contentSrc.Name(), expectedContent, os.FileMode(0666))
					if nil != err {
						panic(err)
					}

					fakeHandle.GetContentReturns(contentSrc, nil)

					fakeOS := new(ios.Fake)

					// create tmpfile to use as dst
					contentDst, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}
					defer contentDst.Close()

					fakeOS.CreateReturns(contentDst, nil)

					objectUnderTest := _Pkg{
						os: fakeOS,
					}

					/* act */
					objectUnderTest.Install(nil, "", fakeHandle)

					/* assert */
					actualContent, err := ioutil.ReadFile(contentDst.Name())
					if nil != err {
						panic(err)
					}

					Expect(actualContent).To(Equal(expectedContent))
				})

				It("shouldn't err", func() {
					/* arrange */
					fakeHandle := new(FakeHandle)
					fakeHandle.ListContentsReturns([]*model.PkgContent{{}}, nil)

					// create tmpfile to use as src
					file, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}

					fakeHandle.GetContentReturns(file, nil)

					fakeOS := new(ios.Fake)
					fakeOS.CreateReturns(file, nil)

					objectUnderTest := _Pkg{
						os: fakeOS,
					}

					/* act */
					actualErr := objectUnderTest.Install(nil, "", fakeHandle)

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})

})
