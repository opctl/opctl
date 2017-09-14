package files

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-interfaces/iio"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _ = Context("Files", func() {
	tempFile, err := ioutil.TempFile("", "")
	if nil != err {
		panic(err)
	}
	Context("Interpret", func() {
		Context("bound value is absolute path", func() {
			It("should call pkgHandle.GetContent w/ expected args", func() {
				/* arrange */

				containerFilePath := "/dummyFile1Path.txt"

				providedSCGContainerCallFiles := map[string]string{
					// implicitly bound
					containerFilePath: "",
				}

				providedParentPkgHandle := new(pkg.FakeHandle)
				// error to trigger immediate return
				providedParentPkgHandle.GetContentReturns(tempFile, errors.New("dummyError"))

				objectUnderTest := _Files{}

				/* act */
				objectUnderTest.Interpret(
					providedParentPkgHandle,
					map[string]*model.Value{},
					providedSCGContainerCallFiles,
					"dummyScratchDir",
				)

				/* assert */
				actualContext, actualContentPath := providedParentPkgHandle.GetContentArgsForCall(0)
				Expect(actualContext).To(Equal(context.TODO()))
				Expect(actualContentPath).To(Equal(containerFilePath))
			})
			Context("pkgHandle.GetContent errs", func() {

				It("should return expected error", func() {
					/* arrange */
					containerFilePath := "/dummyFile1Path.txt"
					providedSCGContainerCallFiles := map[string]string{
						// implicitly bound
						containerFilePath: "",
					}

					getContentErr := fmt.Errorf("dummyError")

					providedParentPkgHandle := new(pkg.FakeHandle)
					providedParentPkgHandle.GetContentReturns(tempFile, getContentErr)

					expectedErr := fmt.Errorf(
						"unable to bind file '%v' to pkg content '%v'; error was: %v",
						containerFilePath,
						containerFilePath,
						getContentErr,
					)

					objectUnderTest := _Files{}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						providedParentPkgHandle,
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("pkgHandle.GetContent doesn't err", func() {
				It("should call pkgHandle.ListContents w/ expected args", func() {
					/* arrange */
					containerFilePath := "/dummyFile1Path.txt"
					providedSCGContainerCallFiles := map[string]string{
						// implicitly bound
						containerFilePath: "",
					}

					providedScratchDirPath := "dummyScratchDirPath"

					providedParentPkgHandle := new(pkg.FakeHandle)
					providedParentPkgHandle.GetContentReturns(tempFile, nil)
					// err to trigger immediate return
					providedParentPkgHandle.ListContentsReturns(nil, errors.New("dummyErr"))

					fakeOS := new(ios.Fake)
					// err to trigger immediate return
					fakeOS.OpenFileReturns(nil, errors.New("dummyError"))

					objectUnderTest := _Files{
						os: fakeOS,
					}

					/* act */
					objectUnderTest.Interpret(
						providedParentPkgHandle,
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						providedScratchDirPath,
					)

					/* assert */
					Expect(providedParentPkgHandle.ListContentsArgsForCall(0)).To(Equal(context.TODO()))
				})
				Context("pkgHandle.ListContents errs", func() {

					It("should return expected error", func() {
						/* arrange */
						containerFilePath := "/dummyFile1Path.txt"
						providedSCGContainerCallFiles := map[string]string{
							// implicitly bound
							containerFilePath: "",
						}

						providedParentPkgHandle := new(pkg.FakeHandle)
						providedParentPkgHandle.GetContentReturns(tempFile, nil)

						getContentErr := fmt.Errorf("dummyError")
						providedParentPkgHandle.ListContentsReturns(nil, getContentErr)

						expectedErr := fmt.Errorf(
							"unable to bind file '%v' to pkg content '%v'; error was: %v",
							containerFilePath,
							containerFilePath,
							getContentErr,
						)

						objectUnderTest := _Files{}

						/* act */
						_, actualErr := objectUnderTest.Interpret(
							providedParentPkgHandle,
							map[string]*model.Value{},
							providedSCGContainerCallFiles,
							"dummyScratchDirPath",
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("pkgHandle.ListContents doesn't err", func() {
					It("should call os.OpenFile w/ expected args", func() {
						/* arrange */
						containerFilePath := "/dummyFile1Path.txt"
						providedSCGContainerCallFiles := map[string]string{
							// implicitly bound
							containerFilePath: "",
						}

						providedScratchDirPath := "dummyScratchDirPath"

						providedParentPkgHandle := new(pkg.FakeHandle)
						providedParentPkgHandle.GetContentReturns(tempFile, nil)

						fakeOS := new(ios.Fake)
						// err to trigger immediate return
						fakeOS.OpenFileReturns(nil, errors.New("dummyError"))

						objectUnderTest := _Files{
							os: fakeOS,
						}

						/* act */
						objectUnderTest.Interpret(
							providedParentPkgHandle,
							map[string]*model.Value{},
							providedSCGContainerCallFiles,
							providedScratchDirPath,
						)

						/* assert */
						actualPath, actualFlags, actualPerm := fakeOS.OpenFileArgsForCall(0)
						Expect(actualPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))
						Expect(actualFlags).To(Equal(os.O_RDWR | os.O_CREATE))
						Expect(actualPerm).To(Equal(os.FileMode(0666)))

					})
					Context("os.OpenFile errs", func() {
						It("should return expected err", func() {
							/* arrange */
							containerFilePath := "/dummyFile1Path.txt"
							providedSCGContainerCallFiles := map[string]string{
								// implicitly bound
								containerFilePath: "",
							}

							providedParentPkgHandle := new(pkg.FakeHandle)
							providedParentPkgHandle.GetContentReturns(tempFile, nil)

							openErr := fmt.Errorf("dummyError")

							fakeOS := new(ios.Fake)
							fakeOS.OpenFileReturns(tempFile, openErr)

							objectUnderTest := _Files{
								os: fakeOS,
							}

							/* act */
							_, actualErr := objectUnderTest.Interpret(
								providedParentPkgHandle,
								map[string]*model.Value{},
								providedSCGContainerCallFiles,
								"dummyScratchDirPath",
							)

							/* assert */
							Expect(actualErr).To(Equal(openErr))
						})
					})
					Context("os.OpenFile doesn't err", func() {
						It("should call io.Copy w/ expected args", func() {
							/* arrange */
							containerFilePath := "/dummyFile1Path.txt"
							providedSCGContainerCallFiles := map[string]string{
								// implicitly bound
								containerFilePath: "",
							}

							readSeekCloser, err := ioutil.TempFile("", "")
							if nil != err {
								panic(err)
							}

							providedParentPkgHandle := new(pkg.FakeHandle)
							providedParentPkgHandle.GetContentReturns(readSeekCloser, nil)

							fakeOS := new(ios.Fake)
							fakeOS.OpenFileReturns(tempFile, nil)

							copyErr := fmt.Errorf("dummyError")

							fakeIO := new(iio.Fake)
							// err to trigger immediate return
							fakeIO.CopyReturns(0, copyErr)

							objectUnderTest := _Files{
								os: fakeOS,
								io: fakeIO,
							}

							/* act */
							objectUnderTest.Interpret(
								providedParentPkgHandle,
								map[string]*model.Value{},
								providedSCGContainerCallFiles,
								"dummyScratchDirPath",
							)

							/* assert */
							actualWriter, actualReader := fakeIO.CopyArgsForCall(0)
							Expect(actualWriter).To(Equal(tempFile))
							Expect(actualReader).To(Equal(readSeekCloser))
						})
						Context("io.Copy errs", func() {
							It("should return expected err", func() {
								/* arrange */
								containerFilePath := "/dummyFile1Path.txt"
								providedSCGContainerCallFiles := map[string]string{
									// implicitly bound
									containerFilePath: "",
								}

								providedParentPkgHandle := new(pkg.FakeHandle)
								providedParentPkgHandle.GetContentReturns(tempFile, nil)

								fakeOS := new(ios.Fake)
								fakeOS.OpenFileReturns(tempFile, nil)

								copyErr := fmt.Errorf("dummyError")

								fakeIO := new(iio.Fake)
								fakeIO.CopyReturns(0, copyErr)

								objectUnderTest := _Files{
									os: fakeOS,
									io: fakeIO,
								}

								/* act */
								_, actualErr := objectUnderTest.Interpret(
									providedParentPkgHandle,
									map[string]*model.Value{},
									providedSCGContainerCallFiles,
									"dummyScratchDirPath",
								)

								/* assert */
								Expect(actualErr).To(Equal(copyErr))
							})
						})
					})
				})
			})
		})
		Context("bound value matches scope name", func() {
			Context("value is nil", func() {
				It("should return expected error", func() {
					/* arrange */
					scopeName := "dummyScopeName"

					providedScope := map[string]*model.Value{
						scopeName: nil,
					}

					containerFilePath := "dummyContainerFilePath"
					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						containerFilePath: scopeName,
					}

					expectedErr := fmt.Errorf(
						"unable to bind file '%v' to '%v'; '%v' null",
						containerFilePath,
						scopeName,
						scopeName,
					)

					objectUnderTest := _Files{
						os: new(ios.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						new(pkg.FakeHandle),
						providedScope,
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("value isn't nil", func() {
				Context("value.File isn't nil", func() {
					Context("value.File not prefixed by rootFSPath", func() {
						It("should return expected results", func() {

							/* arrange */
							scopeName := "dummyScopeName"
							scopeValue := "dummyScopeValue"
							providedScope := map[string]*model.Value{
								scopeName: {File: &scopeValue},
							}

							containerFilePath := "dummyContainerFilePath"
							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								containerFilePath: scopeName,
							}

							expectedDCGContainerCallFiles := map[string]string{
								containerFilePath: scopeValue,
							}

							fakeOS := new(ios.Fake)

							objectUnderTest := _Files{
								io:         new(iio.Fake),
								os:         fakeOS,
								rootFSPath: "dummyRootFSPath",
							}

							/* act */
							actualDCGContainerCallFiles, actualErr := objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								providedScope,
								providedSCGContainerCallFiles,
								"dummyScratchDirPath",
							)

							/* assert */
							Expect(actualDCGContainerCallFiles).To(Equal(expectedDCGContainerCallFiles))
							Expect(actualErr).To(BeNil())
						})
					})
					Context("value.File prefixed by rootFSPath", func() {
						It("should call os.OpenFile w/ expected args", func() {
							/* arrange */
							providedRootFSPath := "dummyRootFSPath"

							scopeName := "dummyScopeName"
							scopeValue := providedRootFSPath

							providedScope := map[string]*model.Value{
								scopeName: {File: &scopeValue},
							}

							containerFilePath := "dummyContainerFilePath"
							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								containerFilePath: scopeName,
							}

							providedScratchDirPath := "dummyScratchDirPath"

							fakeOS := new(ios.Fake)
							// err to trigger immediate return
							fakeOS.OpenFileReturns(nil, errors.New("dummyError"))

							objectUnderTest := _Files{
								os:         fakeOS,
								rootFSPath: providedRootFSPath,
							}

							/* act */
							objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								providedScope,
								providedSCGContainerCallFiles,
								providedScratchDirPath,
							)

							/* assert */
							actualSrcPath, actualFlags, actualPerm := fakeOS.OpenFileArgsForCall(0)
							Expect(actualSrcPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))
							Expect(actualFlags).To(Equal(os.O_RDWR | os.O_CREATE))
							Expect(actualPerm).To(Equal(os.FileMode(0666)))

						})
						Context("os.OpenFile errs", func() {
							It("should return expected error", func() {
								/* arrange */
								scopeName := "dummyScopeName"

								providedScope := map[string]*model.Value{
									scopeName: {File: new(string)},
								}

								containerFilePath := "dummyContainerFilePath"
								providedSCGContainerCallFiles := map[string]string{
									// explicitly bound
									containerFilePath: scopeName,
								}

								fakeOS := new(ios.Fake)
								openErr := errors.New("dummyError")
								fakeOS.OpenFileReturns(nil, openErr)

								objectUnderTest := _Files{
									os: fakeOS,
								}

								/* act */
								_, actualErr := objectUnderTest.Interpret(
									new(pkg.FakeHandle),
									providedScope,
									providedSCGContainerCallFiles,
									"dummyScratchDirPath",
								)

								/* assert */
								Expect(actualErr).To(Equal(openErr))
							})
						})
						Context("os.OpenFile doesn't err", func() {
							It("should call os.MkdirAll w/ expected args", func() {

								/* arrange */
								providedRootFSPath := "dummyRootFSPath"

								scopeName := "dummyScopeName"
								scopeValue := providedRootFSPath

								providedScope := map[string]*model.Value{
									scopeName: {File: &scopeValue},
								}

								containerFilePath := "dummyContainerFilePath"
								providedSCGContainerCallFiles := map[string]string{
									// explicitly bound
									containerFilePath: scopeName,
								}

								providedScratchDirPath := "dummyScratchDirPath"

								fakeOS := new(ios.Fake)
								// err to trigger immediate return
								fakeOS.MkdirAllReturns(errors.New("dummyError"))

								objectUnderTest := _Files{
									os:         fakeOS,
									rootFSPath: providedRootFSPath,
								}

								/* act */
								objectUnderTest.Interpret(
									new(pkg.FakeHandle),
									providedScope,
									providedSCGContainerCallFiles,
									providedScratchDirPath,
								)

								/* assert */
								actualPath, actualPerm := fakeOS.MkdirAllArgsForCall(0)
								Expect(actualPath).To(Equal(filepath.Dir(filepath.Join(providedScratchDirPath, containerFilePath))))
								Expect(actualPerm).To(Equal(os.FileMode(0777)))

							})
						})
					})
				})
				Context("value.File is nil", func() {
					It("should call data.CoerceToString w/ expected args", func() {

						/* arrange */
						scopeName := "dummyScopeName"
						containerFilePath := "dummyContainerFilePath"
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							containerFilePath: scopeName,
						}

						fakeData := new(data.Fake)
						// err to trigger immediate return
						fakeData.CoerceToStringReturns("", errors.New("dummyError"))

						objectUnderTest := _Files{
							data: fakeData,
						}

						expectedValue := &model.Value{String: new(string)}

						/* act */
						objectUnderTest.Interpret(
							new(pkg.FakeHandle),
							map[string]*model.Value{
								scopeName: expectedValue,
							},
							providedSCGContainerCallFiles,
							"dummyScratchDirPath",
						)

						/* assert */
						actualValue := fakeData.CoerceToStringArgsForCall(0)
						Expect(actualValue).To(Equal(expectedValue))

					})
					Context("data.CoerceToString errs", func() {
						It("should return expected result", func() {

							/* arrange */
							scopeName := "dummyScopeName"

							containerFilePath := "dummyContainerFilePath"
							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								containerFilePath: scopeName,
							}

							fakeData := new(data.Fake)

							coerceToStringErr := errors.New("dummyError")
							fakeData.CoerceToStringReturns("", coerceToStringErr)

							expectedErrors := fmt.Errorf(
								"unable to bind file '%v' to '%v'; error was: %v",
								containerFilePath,
								scopeName,
								coerceToStringErr.Error(),
							)

							objectUnderTest := _Files{
								data: fakeData,
							}

							/* act */
							_, actualErr := objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								map[string]*model.Value{
									scopeName: {},
								},
								providedSCGContainerCallFiles,
								"dummyScratchDirPath",
							)

							/* assert */
							Expect(actualErr).To(Equal(expectedErrors))

						})
					})
					Context("data.CoerceToString doesn't err", func() {
						It("should call os.MkdirAll w/ expected args", func() {
							/* arrange */
							scopeName := "dummyScopeName"

							providedScope := map[string]*model.Value{
								scopeName: {},
							}

							containerFilePath := "dummyContainerFilePath"
							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								containerFilePath: scopeName,
							}

							providedScratchDirPath := "dummyScratchDirPath"

							fakeOS := new(ios.Fake)

							objectUnderTest := _Files{
								data: new(data.Fake),
								io:   new(iio.Fake),
								os:   fakeOS,
							}

							/* act */
							objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								providedScope,
								providedSCGContainerCallFiles,
								providedScratchDirPath,
							)

							/* assert */
							actualPath, actualPerm := fakeOS.MkdirAllArgsForCall(0)
							Expect(actualPath).To(Equal(filepath.Dir(filepath.Join(providedScratchDirPath, containerFilePath))))
							Expect(actualPerm).To(Equal(os.FileMode(0777)))

						})
						Context("os.MkdirAll errs", func() {
							It("should return error", func() {

								/* arrange */
								scopeName := "dummyScopeName"

								providedScope := map[string]*model.Value{
									scopeName: {},
								}

								containerFilePath := "dummyContainerFilePath"
								providedSCGContainerCallFiles := map[string]string{
									// explicitly bound
									containerFilePath: scopeName,
								}

								fakeOS := new(ios.Fake)
								mkdirAllErr := errors.New("dummyError")
								fakeOS.MkdirAllReturns(mkdirAllErr)

								expectedErr := fmt.Errorf(
									"unable to mkdir for bound file '%v'; error was: %v",
									containerFilePath,
									mkdirAllErr.Error(),
								)

								objectUnderTest := _Files{
									data: new(data.Fake),
									os:   fakeOS,
								}

								/* act */
								_, actualErr := objectUnderTest.Interpret(
									new(pkg.FakeHandle),
									providedScope,
									providedSCGContainerCallFiles,
									"dummyScratchDirPath",
								)

								/* assert */
								Expect(actualErr).To(Equal(expectedErr))
							})
						})
						Context("os.MkdirAll doesn't err", func() {
							It("should call os.OpenFile w/ expected args", func() {
								/* arrange */
								scopeName := "dummyScopeName"

								providedScope := map[string]*model.Value{
									scopeName: {},
								}

								containerFilePath := "dummyContainerFilePath"
								providedSCGContainerCallFiles := map[string]string{
									// explicitly bound
									containerFilePath: scopeName,
								}

								providedScratchDirPath := "dummyScratchDirPath"

								fakeOS := new(ios.Fake)

								objectUnderTest := _Files{
									data: new(data.Fake),
									io:   new(iio.Fake),
									os:   fakeOS,
								}

								/* act */
								objectUnderTest.Interpret(
									new(pkg.FakeHandle),
									providedScope,
									providedSCGContainerCallFiles,
									providedScratchDirPath,
								)

								/* assert */
								actualPath, actualFlags, actualPerm := fakeOS.OpenFileArgsForCall(0)
								Expect(actualPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))
								Expect(actualFlags).To(Equal(os.O_RDWR | os.O_CREATE))
								Expect(actualPerm).To(Equal(os.FileMode(0666)))

							})
							Context("os.OpenFile errs", func() {
								It("should return error", func() {

									/* arrange */
									scopeName := "dummyScopeName"

									providedScope := map[string]*model.Value{
										scopeName: {},
									}

									providedSCGContainerCallFiles := map[string]string{
										// explicitly bound
										"dummyContainerFilePath": scopeName,
									}

									expectedErr := errors.New("dummyError")

									fakeOS := new(ios.Fake)
									fakeOS.OpenFileReturns(nil, expectedErr)

									objectUnderTest := _Files{
										data: new(data.Fake),
										os:   fakeOS,
									}

									/* act */
									_, actualErr := objectUnderTest.Interpret(
										new(pkg.FakeHandle),
										providedScope,
										providedSCGContainerCallFiles,
										"dummyScratchDirPath",
									)

									/* assert */
									Expect(actualErr).To(Equal(expectedErr))
								})
							})
							Context("os.OpenFile doesn't err", func() {
								It("should call io.Copy w/ expected args", func() {
									/* arrange */
									scopeName := "dummyScopeName"

									providedScope := map[string]*model.Value{
										scopeName: {},
									}

									providedSCGContainerCallFiles := map[string]string{
										// explicitly bound
										"dummyContainerFilePath": scopeName,
									}

									providedScratchDirPath := "dummyScratchDirPath"
									expectedCopyReader := strings.NewReader("")

									fakeIO := new(iio.Fake)

									fakeOS := new(ios.Fake)
									expectedCopyWriter, err := ioutil.TempFile("", "")
									fakeOS.OpenFileReturns(expectedCopyWriter, err)

									objectUnderTest := _Files{
										data: new(data.Fake),
										io:   fakeIO,
										os:   fakeOS,
									}

									/* act */
									objectUnderTest.Interpret(
										new(pkg.FakeHandle),
										providedScope,
										providedSCGContainerCallFiles,
										providedScratchDirPath,
									)

									/* assert */
									actualCopyWriter, actualCopyReader := fakeIO.CopyArgsForCall(0)
									Expect(actualCopyReader).To(Equal(expectedCopyReader))
									Expect(actualCopyWriter).To(Equal(expectedCopyWriter))
								})
								Context("io.Copy errs", func() {
									It("should return error", func() {

										/* arrange */
										scopeName := "dummyScopeName"

										providedScope := map[string]*model.Value{
											scopeName: {},
										}

										providedSCGContainerCallFiles := map[string]string{
											// explicitly bound
											"dummyContainerFilePath": scopeName,
										}

										expectedErr := errors.New("dummyError")

										fakeIO := new(iio.Fake)
										fakeIO.CopyReturns(0, expectedErr)

										objectUnderTest := _Files{
											data: new(data.Fake),
											io:   fakeIO,
											os:   new(ios.Fake),
										}

										/* act */
										_, actualErr := objectUnderTest.Interpret(
											new(pkg.FakeHandle),
											providedScope,
											providedSCGContainerCallFiles,
											"dummyScratchDirPath",
										)

										/* assert */
										Expect(actualErr).To(Equal(expectedErr))
									})
								})
								Context("io.Copy doesn't err", func() {
									It("should return expected results", func() {

										/* arrange */
										scopeName := "dummyScopeName"

										providedScope := map[string]*model.Value{
											scopeName: {},
										}

										containerFilePath := "dummyContainerFilePath"
										providedSCGContainerCallFiles := map[string]string{
											// explicitly bound
											containerFilePath: scopeName,
										}

										providedScratchDirPath := "dummyScratchDirPath"

										expectedDCGContainerCallFiles := map[string]string{
											containerFilePath: filepath.Join(providedScratchDirPath, containerFilePath),
										}

										fakeOS := new(ios.Fake)

										objectUnderTest := _Files{
											data: new(data.Fake),
											io:   new(iio.Fake),
											os:   fakeOS,
										}

										/* act */
										actualDCGContainerCallFiles, actualErr := objectUnderTest.Interpret(
											new(pkg.FakeHandle),
											providedScope,
											providedSCGContainerCallFiles,
											providedScratchDirPath,
										)

										/* assert */
										Expect(actualDCGContainerCallFiles).To(Equal(expectedDCGContainerCallFiles))
										Expect(actualErr).To(BeNil())
									})
								})
							})
						})
					})
				})
			})
		})
		Context("bound value doesn't match scope name", func() {
			It("should call os.MkdirAll w/ expected args", func() {
				/* arrange */
				containerFilePath := "dummyContainerFilePath"
				providedSCGContainerCallFiles := map[string]string{
					// explicitly bound
					containerFilePath: "dummyScopeName",
				}

				providedScratchDirPath := "dummyScratchDirPath"

				fakeOS := new(ios.Fake)

				objectUnderTest := _Files{
					io: new(iio.Fake),
					os: fakeOS,
				}

				/* act */
				objectUnderTest.Interpret(
					new(pkg.FakeHandle),
					map[string]*model.Value{},
					providedSCGContainerCallFiles,
					providedScratchDirPath,
				)

				/* assert */
				actualPath, actualPerm := fakeOS.MkdirAllArgsForCall(0)
				Expect(actualPath).To(Equal(filepath.Dir(filepath.Join(providedScratchDirPath, containerFilePath))))
				Expect(actualPerm).To(Equal(os.FileMode(0777)))

			})
			Context("os.MkdirAll errs", func() {
				It("should return error", func() {

					/* arrange */
					containerFilePath := "dummyContainerFilePath"
					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						containerFilePath: "dummyScopeName",
					}

					fakeOS := new(ios.Fake)
					mkdirAllErr := errors.New("dummyError")
					fakeOS.MkdirAllReturns(mkdirAllErr)

					expectedErr := fmt.Errorf(
						"unable to mkdir for bound file '%v'; error was: %v",
						containerFilePath,
						mkdirAllErr.Error(),
					)

					objectUnderTest := _Files{
						os: fakeOS,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("os.MkdirAll doesn't err", func() {
				It("should call os.OpenFile w/ expected args", func() {
					/* arrange */
					containerFilePath := "dummyContainerFilePath"
					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						containerFilePath: "dummyScopeName",
					}

					providedScratchDirPath := "dummyScratchDirPath"

					fakeOS := new(ios.Fake)

					objectUnderTest := _Files{
						io: new(iio.Fake),
						os: fakeOS,
					}

					/* act */
					objectUnderTest.Interpret(
						new(pkg.FakeHandle),
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						providedScratchDirPath,
					)

					/* assert */
					actualPath, actualFlags, actualPerm := fakeOS.OpenFileArgsForCall(0)
					Expect(actualPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))
					Expect(actualFlags).To(Equal(os.O_RDWR | os.O_CREATE))
					Expect(actualPerm).To(Equal(os.FileMode(0666)))

				})
				Context("os.OpenFile errs", func() {
					It("should return error", func() {

						/* arrange */
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							"dummyContainerFilePath": "dummyScopeName",
						}

						expectedErr := errors.New("dummyError")

						fakeOS := new(ios.Fake)
						fakeOS.OpenFileReturns(nil, expectedErr)

						objectUnderTest := _Files{
							os: fakeOS,
						}

						/* act */
						_, actualErr := objectUnderTest.Interpret(
							new(pkg.FakeHandle),
							map[string]*model.Value{},
							providedSCGContainerCallFiles,
							"dummyScratchDirPath",
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("os.OpenFile doesn't err", func() {
					It("should call io.Copy w/ expected args", func() {
						/* arrange */
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							"dummyContainerFilePath": "dummyScopeName",
						}

						providedScratchDirPath := "dummyScratchDirPath"
						expectedCopyReader := strings.NewReader("")

						fakeIO := new(iio.Fake)

						fakeOS := new(ios.Fake)
						expectedCopyWriter, err := ioutil.TempFile("", "")
						fakeOS.OpenFileReturns(expectedCopyWriter, err)

						objectUnderTest := _Files{
							io: fakeIO,
							os: fakeOS,
						}

						/* act */
						objectUnderTest.Interpret(
							new(pkg.FakeHandle),
							map[string]*model.Value{},
							providedSCGContainerCallFiles,
							providedScratchDirPath,
						)

						/* assert */
						actualCopyWriter, actualCopyReader := fakeIO.CopyArgsForCall(0)
						Expect(actualCopyReader).To(Equal(expectedCopyReader))
						Expect(actualCopyWriter).To(Equal(expectedCopyWriter))
					})
					Context("io.Copy errs", func() {
						It("should return error", func() {

							/* arrange */
							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								"dummyContainerFilePath": "dummyScopeName",
							}

							expectedErr := errors.New("dummyError")

							fakeIO := new(iio.Fake)
							fakeIO.CopyReturns(0, expectedErr)

							objectUnderTest := _Files{
								io: fakeIO,
								os: new(ios.Fake),
							}

							/* act */
							_, actualErr := objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								map[string]*model.Value{},
								providedSCGContainerCallFiles,
								"dummyScratchDirPath",
							)

							/* assert */
							Expect(actualErr).To(Equal(expectedErr))
						})
					})
					Context("io.Copy doesn't err", func() {
						It("should return expected results", func() {

							/* arrange */
							containerFilePath := "dummyContainerFilePath"
							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								containerFilePath: "dummyScopeName",
							}

							providedScratchDirPath := "dummyScratchDirPath"

							expectedDCGContainerCallFiles := map[string]string{
								containerFilePath: filepath.Join(providedScratchDirPath, containerFilePath),
							}

							fakeOS := new(ios.Fake)

							objectUnderTest := _Files{
								io: new(iio.Fake),
								os: fakeOS,
							}

							/* act */
							actualDCGContainerCallFiles, actualErr := objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								map[string]*model.Value{},
								providedSCGContainerCallFiles,
								providedScratchDirPath,
							)

							/* assert */
							Expect(actualDCGContainerCallFiles).To(Equal(expectedDCGContainerCallFiles))
							Expect(actualErr).To(BeNil())
						})
					})
				})
			})
		})
	})
})
