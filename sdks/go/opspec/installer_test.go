package opspec

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("Installer", func() {
	Context("NewInstaller", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewInstaller()).Should(Not(BeNil()))
		})
	})
	Context("Install", func() {
		It("should call handle.ListDescendants w/ expected args", func() {
			/* arrange */
			providedCtx := context.TODO()

			fakeHandle := new(modelFakes.FakeDataHandle)

			objectUnderTest := _installer{
				os: new(ios.Fake),
			}

			/* act */
			objectUnderTest.Install(providedCtx, "", fakeHandle)

			/* assert */
			Expect(fakeHandle.ListDescendantsArgsForCall(0)).To(Equal(providedCtx))
		})
		Context("handle.ListDescendants errs", func() {
			It("should return error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeHandle := new(modelFakes.FakeDataHandle)
				fakeHandle.ListDescendantsReturns(nil, expectedError)

				objectUnderTest := _installer{
					os: new(ios.Fake),
				}

				/* act */
				actualError := objectUnderTest.Install(nil, "", fakeHandle)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("handle.ListDescendants doesn't err", func() {
			It("should call handle.GetContent w/ expected args", func() {
				/* arrange */
				providedCtx := context.TODO()

				fakeHandle := new(modelFakes.FakeDataHandle)
				contentsList := []*model.DirEntry{
					{
						Path: "dirEntry1Path",
					},
				}

				fakeHandle.ListDescendantsReturns(
					contentsList,
					nil,
				)

				// error to trigger immediate return
				fakeHandle.GetContentReturns(nil, errors.New("dummyError"))

				objectUnderTest := _installer{
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

					fakeHandle := new(modelFakes.FakeDataHandle)
					fakeHandle.ListDescendantsReturns([]*model.DirEntry{{}}, expectedError)

					fakeHandle.GetContentReturns(nil, expectedError)

					objectUnderTest := _installer{
						os: new(ios.Fake),
					}

					/* act */
					actualError := objectUnderTest.Install(nil, "", fakeHandle)

					/* assert */
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("handle.GetContent doesn't err", func() {
				Context("content.Mode.IsDir() == true", func() {
					It("should call os.MkdirAll w/ expected args", func() {
						/* arrange */
						providedPath := "dummyPath"

						fakeHandle := new(modelFakes.FakeDataHandle)

						contentsList := []*model.DirEntry{
							{
								Path: "dirEntry1Path",
								Mode: os.ModeDir,
							},
						}

						fakeHandle.ListDescendantsReturns(
							contentsList,
							nil,
						)

						fakeOS := new(ios.Fake)
						// error to trigger immediate return
						fakeOS.MkdirAllReturns(errors.New("dummyError"))

						objectUnderTest := _installer{
							os: fakeOS,
						}

						/* act */
						objectUnderTest.Install(nil, providedPath, fakeHandle)

						/* assert */
						actualPath, actualPerm := fakeOS.MkdirAllArgsForCall(0)

						Expect(actualPath).To(Equal(
							filepath.Join(providedPath, contentsList[0].Path),
						))

						Expect(actualPerm).To(Equal(contentsList[0].Mode))
					})
					Context("os.MkdirAll errs", func() {
						It("should return error", func() {
							/* arrange */
							expectedError := errors.New("dummyError")

							fakeHandle := new(modelFakes.FakeDataHandle)
							fakeHandle.ListDescendantsReturns([]*model.DirEntry{{Mode: os.ModeDir}}, nil)

							fakeOS := new(ios.Fake)
							fakeOS.MkdirAllReturns(expectedError)

							objectUnderTest := _installer{
								os: fakeOS,
							}

							/* act */
							actualError := objectUnderTest.Install(nil, "", fakeHandle)

							/* assert */
							Expect(actualError).To(Equal(expectedError))
						})
					})
				})
				Context("content.Mode.IsDir() == false", func() {
					It("should call os.MkdirAll w/ expected args", func() {
						/* arrange */
						providedPath := "dummyPath"

						fakeHandle := new(modelFakes.FakeDataHandle)

						contentsList := []*model.DirEntry{
							{
								Path: "dirEntry1Path",
							},
						}

						fakeHandle.ListDescendantsReturns(
							contentsList,
							nil,
						)

						fakeOS := new(ios.Fake)
						// error to trigger immediate return
						fakeOS.MkdirAllReturns(errors.New("dummyError"))

						objectUnderTest := _installer{
							os: fakeOS,
						}

						/* act */
						objectUnderTest.Install(nil, providedPath, fakeHandle)

						/* assert */
						actualPath, actualPerm := fakeOS.MkdirAllArgsForCall(0)

						Expect(actualPath).To(Equal(
							filepath.Dir(
								filepath.Join(providedPath, contentsList[0].Path),
							),
						))

						Expect(actualPerm).To(Equal(os.FileMode(0777)))
					})
					Context("os.MkdirAll errs", func() {
						It("should return error", func() {
							/* arrange */
							expectedError := errors.New("dummyError")

							fakeHandle := new(modelFakes.FakeDataHandle)
							fakeHandle.ListDescendantsReturns([]*model.DirEntry{{}}, nil)

							fakeOS := new(ios.Fake)
							fakeOS.MkdirAllReturns(expectedError)

							objectUnderTest := _installer{
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

							fakeHandle := new(modelFakes.FakeDataHandle)
							contentsList := []*model.DirEntry{
								{
									Path: "dirEntry1Path",
								},
							}

							fakeHandle.ListDescendantsReturns(
								contentsList,
								nil,
							)

							fakeOS := new(ios.Fake)
							// error to trigger immediate return
							fakeOS.CreateReturns(nil, errors.New("dummyError"))

							objectUnderTest := _installer{
								os: fakeOS,
							}

							/* act */
							objectUnderTest.Install(nil, providedPath, fakeHandle)

							/* assert */
							actualPath := fakeOS.CreateArgsForCall(0)

							Expect(actualPath).To(Equal(filepath.Join(providedPath, contentsList[0].Path)))
						})
						Context("os.Create errs", func() {
							It("should return error", func() {
								/* arrange */
								expectedError := errors.New("dummyError")

								fakeHandle := new(modelFakes.FakeDataHandle)
								fakeHandle.ListDescendantsReturns([]*model.DirEntry{{}}, nil)

								fakeOS := new(ios.Fake)
								fakeOS.CreateReturns(nil, expectedError)

								objectUnderTest := _installer{
									os: fakeOS,
								}

								/* act */
								actualError := objectUnderTest.Install(nil, "", fakeHandle)

								/* assert */
								Expect(actualError).To(Equal(expectedError))
							})
						})
						Context("os.Create doesn't err", func() {
							It("should call os.Chmod w/ expected args", func() {
								/* arrange */
								providedPath := "dummyPath"

								fakeHandle := new(modelFakes.FakeDataHandle)
								contentsList := []*model.DirEntry{
									{
										Mode: os.FileMode(0777),
										Path: "dirEntry1Path",
									},
								}

								fakeHandle.ListDescendantsReturns(
									contentsList,
									nil,
								)

								fakeOS := new(ios.Fake)
								// error to trigger immediate return
								fakeOS.ChmodReturns(errors.New("dummyError"))

								objectUnderTest := _installer{
									os: fakeOS,
								}

								/* act */
								objectUnderTest.Install(nil, providedPath, fakeHandle)

								/* assert */
								actualPath, actualMode := fakeOS.ChmodArgsForCall(0)

								Expect(actualPath).To(Equal(filepath.Join(providedPath, contentsList[0].Path)))
								Expect(actualMode).To(Equal(contentsList[0].Mode))

							})
							Context("os.Chmod errs", func() {

								It("should return error", func() {
									/* arrange */
									expectedError := errors.New("dummyError")

									fakeHandle := new(modelFakes.FakeDataHandle)
									fakeHandle.ListDescendantsReturns([]*model.DirEntry{{}}, nil)

									fakeOS := new(ios.Fake)
									fakeOS.ChmodReturns(expectedError)

									objectUnderTest := _installer{
										os: fakeOS,
									}

									/* act */
									actualError := objectUnderTest.Install(nil, "", fakeHandle)

									/* assert */
									Expect(actualError).To(Equal(expectedError))
								})
							})
							Context("os.Chmod doesn't err", func() {
								It("should copy content", func() {
									/* arrange */
									fakeHandle := new(modelFakes.FakeDataHandle)
									fakeHandle.ListDescendantsReturns([]*model.DirEntry{{}}, nil)

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

									objectUnderTest := _installer{
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
									fakeHandle := new(modelFakes.FakeDataHandle)
									fakeHandle.ListDescendantsReturns([]*model.DirEntry{{}}, nil)

									// create tmpfile to use as src
									file, err := ioutil.TempFile("", "")
									if nil != err {
										panic(err)
									}

									fakeHandle.GetContentReturns(file, nil)

									fakeOS := new(ios.Fake)
									fakeOS.CreateReturns(file, nil)

									objectUnderTest := _installer{
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
			})
		})
	})
})
