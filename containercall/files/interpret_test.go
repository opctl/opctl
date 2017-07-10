package files

import (
	"fmt"
	"github.com/golang-interfaces/iio"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var _ = Context("Files", func() {
	Context("Interpret", func() {
		Context("bound value is absolute path", func() {
			It("should call fileCopier.OS w/ expected args", func() {
				/* arrange */

				containerFilePath := "/dummyFile1Path.txt"

				providedSCGContainerCallFiles := map[string]string{
					// implicitly bound
					containerFilePath: "",
				}

				providedPkgPath := "dummyPkgPath"

				providedScratchDirPath := "dummyScratchDirPath"

				fakeFileCopier := new(filecopier.Fake)
				// error to trigger immediate return
				fakeFileCopier.OSReturns(errors.New("dummyError"))

				objectUnderTest := _Files{
					io:         new(iio.Fake),
					fileCopier: fakeFileCopier,
				}

				/* act */
				objectUnderTest.Interpret(
					providedPkgPath,
					map[string]*model.Value{},
					providedSCGContainerCallFiles,
					providedScratchDirPath,
				)

				/* assert */

				actualSrcPath, actualDstPath := fakeFileCopier.OSArgsForCall(0)
				Expect(actualSrcPath).To(Equal(filepath.Join(providedPkgPath, containerFilePath)))
				Expect(actualDstPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))
			})
			Context("fileCopier.OS errs", func() {

				It("should return expected error", func() {
					/* arrange */
					containerFilePath := "/dummyFile1Path.txt"
					providedSCGContainerCallFiles := map[string]string{
						// implicitly bound
						containerFilePath: "",
					}

					fakeFileCopier := new(filecopier.Fake)
					openError := fmt.Errorf("dummyError")
					fakeFileCopier.OSReturns(openError)

					expectedErr := fmt.Errorf(
						"Unable to bind file '%v' to pkg content '%v'; error was: %v",
						containerFilePath,
						containerFilePath,
						openError,
					)

					objectUnderTest := _Files{
						fileCopier: fakeFileCopier,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						"dummyPkgPath",
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("fileCopier.OS doesn't err", func() {
				It("should return expected results", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"

					containerFilePath := "/dummyFile1Path.txt"
					expectedDCGContainerCallFiles := map[string]string{
						containerFilePath: filepath.Join(providedScratchDir, containerFilePath),
					}

					providedSCGContainerCallFiles := map[string]string{
						// implicitly bound
						containerFilePath: "",
					}

					objectUnderTest := _Files{
						fileCopier: new(filecopier.Fake),
					}

					/* act */
					actualDCGContainerCallFiles, actualErr := objectUnderTest.Interpret(
						"dummyPkgPath",
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						providedScratchDir,
					)

					/* assert */
					Expect(actualDCGContainerCallFiles).To(Equal(expectedDCGContainerCallFiles))
					Expect(actualErr).To(BeNil())
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
						"Unable to bind file '%v' to '%v'; '%v' null",
						containerFilePath,
						scopeName,
						scopeName,
					)

					objectUnderTest := _Files{
						os: new(ios.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						"dummyPkgPath",
						providedScope,
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("value.Socket not nil", func() {
				It("should return expected error", func() {
					/* arrange */
					scopeName := "dummyScopeName"

					providedScope := map[string]*model.Value{
						scopeName: {Socket: new(string)},
					}

					containerFilePath := "dummyContainerFilePath"
					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						containerFilePath: scopeName,
					}

					expectedErr := fmt.Errorf("Unable to bind file '%v' to '%v'; '%v' not a file, number, or string", containerFilePath, scopeName, scopeName)

					objectUnderTest := _Files{
						os: new(ios.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						"dummyPkgPath",
						providedScope,
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("value.Dir not nil", func() {
				It("should return expected error", func() {
					/* arrange */
					scopeName := "dummyScopeName"

					providedScope := map[string]*model.Value{
						scopeName: {Dir: new(string)},
					}

					containerFilePath := "dummyContainerFilePath"
					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						containerFilePath: scopeName,
					}

					expectedErr := fmt.Errorf("Unable to bind file '%v' to '%v'; '%v' not a file, number, or string", containerFilePath, scopeName, scopeName)

					objectUnderTest := _Files{
						os: new(ios.Fake),
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						"dummyPkgPath",
						providedScope,
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("value.File not nil", func() {
				Context("value.File prefixed by rootFSPath", func() {
					It("should call fileCopier.OS w/ expected args", func() {
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

						fakeFileCopier := new(filecopier.Fake)
						// err to trigger immediate return
						fakeFileCopier.OSReturns(errors.New("dummyError"))

						objectUnderTest := _Files{
							fileCopier: fakeFileCopier,
							rootFSPath: providedRootFSPath,
						}

						/* act */
						objectUnderTest.Interpret(
							"dummyPkgPath",
							providedScope,
							providedSCGContainerCallFiles,
							providedScratchDirPath,
						)

						/* assert */
						actualSrcPath, actualDstPath := fakeFileCopier.OSArgsForCall(0)
						Expect(actualSrcPath).To(Equal(scopeValue))
						Expect(actualDstPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))
					})
					Context("fileCopier.OS errs", func() {
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

							fakeFileCopier := new(filecopier.Fake)
							copyErr := errors.New("dummyError")
							fakeFileCopier.OSReturns(copyErr)

							expectedErr := fmt.Errorf(
								"Unable to bind file '%v' to '%v'; error was: %v",
								containerFilePath,
								scopeName,
								copyErr,
							)

							objectUnderTest := _Files{
								fileCopier: fakeFileCopier,
							}

							/* act */
							_, actualErr := objectUnderTest.Interpret(
								"dummyPkgPath",
								providedScope,
								providedSCGContainerCallFiles,
								"dummyScratchDirPath",
							)

							/* assert */
							Expect(actualErr).To(Equal(expectedErr))
						})
					})
					Context("fileCopier.OS doesn't err", func() {
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

							providedScratchDirPath := "dummyScratchDirPath"

							expectedDCGContainerCallFiles := map[string]string{
								containerFilePath: filepath.Join(providedScratchDirPath, containerFilePath),
							}

							objectUnderTest := _Files{
								fileCopier: new(filecopier.Fake),
							}

							/* act */
							actualDCGContainerCallFiles, actualErr := objectUnderTest.Interpret(
								"dummyPkgPath",
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
							"dummyPkgPath",
							providedScope,
							providedSCGContainerCallFiles,
							"dummyScratchDirPath",
						)

						/* assert */
						Expect(actualDCGContainerCallFiles).To(Equal(expectedDCGContainerCallFiles))
						Expect(actualErr).To(BeNil())
					})
				})
			})
			Context("value.Number not nil", func() {
				It("should call os.MkdirAll w/ expected args", func() {
					/* arrange */
					scopeName := "dummyScopeName"

					providedScope := map[string]*model.Value{
						scopeName: {Number: new(float64)},
					}

					containerFilePath := "dummyContainerFilePath"
					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						containerFilePath: scopeName,
					}

					providedScratchDirPath := "dummyScratchDirPath"

					fakeOS := new(ios.Fake)

					objectUnderTest := _Files{
						io: new(iio.Fake),
						os: fakeOS,
					}

					/* act */
					objectUnderTest.Interpret(
						"dummyPkgPath",
						providedScope,
						providedSCGContainerCallFiles,
						providedScratchDirPath,
					)

					/* assert */
					actualPath, actualFileMode := fakeOS.MkdirAllArgsForCall(0)
					Expect(actualPath).To(Equal(filepath.Dir(filepath.Join(providedScratchDirPath, containerFilePath))))
					Expect(actualFileMode).To(Equal(os.FileMode(0700)))

				})
				Context("os.MkdirAll errs", func() {
					It("should return error", func() {

						/* arrange */
						scopeName := "dummyScopeName"

						providedScope := map[string]*model.Value{
							scopeName: {Number: new(float64)},
						}

						containerFilePath := "dummyContainerFilePath"
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							containerFilePath: scopeName,
						}

						expectedErr := errors.New("dummyError")

						fakeOS := new(ios.Fake)
						fakeOS.MkdirAllReturns(expectedErr)

						objectUnderTest := _Files{
							os: fakeOS,
						}

						/* act */
						_, actualErr := objectUnderTest.Interpret(
							"dummyPkgPath",
							providedScope,
							providedSCGContainerCallFiles,
							"dummyScratchDirPath",
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("os.MkdirAll doesn't err", func() {
					It("should call os.Create w/ expected args", func() {
						/* arrange */
						scopeName := "dummyScopeName"

						providedScope := map[string]*model.Value{
							scopeName: {Number: new(float64)},
						}

						containerFilePath := "dummyContainerFilePath"
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							containerFilePath: scopeName,
						}

						providedScratchDirPath := "dummyScratchDirPath"

						fakeOS := new(ios.Fake)

						objectUnderTest := _Files{
							io: new(iio.Fake),
							os: fakeOS,
						}

						/* act */
						objectUnderTest.Interpret(
							"dummyPkgPath",
							providedScope,
							providedSCGContainerCallFiles,
							providedScratchDirPath,
						)

						/* assert */
						actualPath := fakeOS.CreateArgsForCall(0)
						Expect(actualPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))

					})
					Context("os.Create errs", func() {
						It("should return error", func() {

							/* arrange */
							scopeName := "dummyScopeName"

							providedScope := map[string]*model.Value{
								scopeName: {Number: new(float64)},
							}

							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								"dummyContainerFilePath": scopeName,
							}

							expectedErr := errors.New("dummyError")

							fakeOS := new(ios.Fake)
							fakeOS.CreateReturns(nil, expectedErr)

							objectUnderTest := _Files{
								os: fakeOS,
							}

							/* act */
							_, actualErr := objectUnderTest.Interpret(
								"dummyPkgPath",
								providedScope,
								providedSCGContainerCallFiles,
								"dummyScratchDirPath",
							)

							/* assert */
							Expect(actualErr).To(Equal(expectedErr))
						})
					})
					Context("os.Create doesn't err", func() {
						It("should call io.Copy w/ expected args", func() {
							/* arrange */
							scopeName := "dummyScopeName"

							numberValue := 33.3
							providedScope := map[string]*model.Value{
								scopeName: {Number: &numberValue},
							}

							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								"dummyContainerFilePath": scopeName,
							}

							providedScratchDirPath := "dummyScratchDirPath"
							expectedCopyReader := strings.NewReader(strconv.FormatFloat(numberValue, 'f', -1, 64))

							fakeIO := new(iio.Fake)

							fakeOS := new(ios.Fake)
							expectedCopyWriter, err := ioutil.TempFile("", "")
							fakeOS.CreateReturns(expectedCopyWriter, err)

							objectUnderTest := _Files{
								io: fakeIO,
								os: fakeOS,
							}

							/* act */
							objectUnderTest.Interpret(
								"dummyPkgPath",
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
									scopeName: {Number: new(float64)},
								}

								providedSCGContainerCallFiles := map[string]string{
									// explicitly bound
									"dummyContainerFilePath": scopeName,
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
									"dummyPkgPath",
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
									scopeName: {Number: new(float64)},
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
									io: new(iio.Fake),
									os: fakeOS,
								}

								/* act */
								actualDCGContainerCallFiles, actualErr := objectUnderTest.Interpret(
									"dummyPkgPath",
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
			Context("value.String not nil", func() {
				It("should call os.MkdirAll w/ expected args", func() {
					/* arrange */
					scopeName := "dummyScopeName"

					providedScope := map[string]*model.Value{
						scopeName: {String: new(string)},
					}

					containerFilePath := "dummyContainerFilePath"
					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						containerFilePath: scopeName,
					}

					providedScratchDirPath := "dummyScratchDirPath"

					fakeOS := new(ios.Fake)

					objectUnderTest := _Files{
						io: new(iio.Fake),
						os: fakeOS,
					}

					/* act */
					objectUnderTest.Interpret(
						"dummyPkgPath",
						providedScope,
						providedSCGContainerCallFiles,
						providedScratchDirPath,
					)

					/* assert */
					actualPath, actualFileMode := fakeOS.MkdirAllArgsForCall(0)
					Expect(actualPath).To(Equal(filepath.Dir(filepath.Join(providedScratchDirPath, containerFilePath))))
					Expect(actualFileMode).To(Equal(os.FileMode(0700)))

				})
				Context("os.MkdirAll errs", func() {
					It("should return error", func() {

						/* arrange */
						scopeName := "dummyScopeName"

						providedScope := map[string]*model.Value{
							scopeName: {String: new(string)},
						}

						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							"dummyContainerFilePath": scopeName,
						}

						expectedErr := errors.New("dummyError")

						fakeOS := new(ios.Fake)
						fakeOS.MkdirAllReturns(expectedErr)

						objectUnderTest := _Files{
							os: fakeOS,
						}

						/* act */
						_, actualErr := objectUnderTest.Interpret(
							"dummyPkgPath",
							providedScope,
							providedSCGContainerCallFiles,
							"dummyScratchDirPath",
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("os.MkdirAll doesn't err", func() {
					It("should call os.Create w/ expected args", func() {
						/* arrange */
						scopeName := "dummyScopeName"

						providedScope := map[string]*model.Value{
							scopeName: {String: new(string)},
						}

						containerFilePath := "dummyContainerFilePath"
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							containerFilePath: scopeName,
						}

						providedScratchDirPath := "dummyScratchDirPath"

						fakeOS := new(ios.Fake)

						objectUnderTest := _Files{
							io: new(iio.Fake),
							os: fakeOS,
						}

						/* act */
						objectUnderTest.Interpret(
							"dummyPkgPath",
							providedScope,
							providedSCGContainerCallFiles,
							providedScratchDirPath,
						)

						/* assert */
						actualPath := fakeOS.CreateArgsForCall(0)
						Expect(actualPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))

					})
					Context("os.Create errs", func() {
						It("should return error", func() {

							/* arrange */
							scopeName := "dummyScopeName"

							providedScope := map[string]*model.Value{
								scopeName: {String: new(string)},
							}

							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								"dummyContainerFilePath": scopeName,
							}

							expectedErr := errors.New("dummyError")

							fakeOS := new(ios.Fake)
							fakeOS.CreateReturns(nil, expectedErr)

							objectUnderTest := _Files{
								os: fakeOS,
							}

							/* act */
							_, actualErr := objectUnderTest.Interpret(
								"dummyPkgPath",
								providedScope,
								providedSCGContainerCallFiles,
								"dummyScratchDirPath",
							)

							/* assert */
							Expect(actualErr).To(Equal(expectedErr))
						})
					})
					Context("os.Create doesn't err", func() {
						It("should call io.Copy w/ expected args", func() {
							/* arrange */
							scopeName := "dummyScopeName"

							stringValue := "dummyValue"
							providedScope := map[string]*model.Value{
								scopeName: {String: &stringValue},
							}

							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								"dummyContainerFilePath": scopeName,
							}

							providedScratchDirPath := "dummyScratchDirPath"
							expectedCopyReader := strings.NewReader(stringValue)

							fakeIO := new(iio.Fake)

							fakeOS := new(ios.Fake)
							expectedCopyWriter, err := ioutil.TempFile("", "")
							fakeOS.CreateReturns(expectedCopyWriter, err)

							objectUnderTest := _Files{
								io: fakeIO,
								os: fakeOS,
							}

							/* act */
							objectUnderTest.Interpret(
								"dummyPkgPath",
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
								scopeName := "dummyContainerFileBind"

								providedScope := map[string]*model.Value{
									scopeName: {String: new(string)},
								}

								providedSCGContainerCallFiles := map[string]string{
									// explicitly bound
									"dummyContainerFilePath": scopeName,
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
									"dummyPkgPath",
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
									scopeName: {String: new(string)},
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
									io: new(iio.Fake),
									os: fakeOS,
								}

								/* act */
								actualDCGContainerCallFiles, actualErr := objectUnderTest.Interpret(
									"dummyPkgPath",
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
					"dummyPkgPath",
					map[string]*model.Value{},
					providedSCGContainerCallFiles,
					providedScratchDirPath,
				)

				/* assert */
				actualPath, actualFileMode := fakeOS.MkdirAllArgsForCall(0)
				Expect(actualPath).To(Equal(filepath.Dir(filepath.Join(providedScratchDirPath, containerFilePath))))
				Expect(actualFileMode).To(Equal(os.FileMode(0700)))

			})
			Context("os.MkdirAll errs", func() {
				It("should return error", func() {

					/* arrange */
					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						"dummyContainerFilePath": "dummyScopeName",
					}

					expectedErr := errors.New("dummyError")

					fakeOS := new(ios.Fake)
					fakeOS.MkdirAllReturns(expectedErr)

					objectUnderTest := _Files{
						os: fakeOS,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						"dummyPkgPath",
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("os.MkdirAll doesn't err", func() {
				It("should call os.Create w/ expected args", func() {
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
						"dummyPkgPath",
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						providedScratchDirPath,
					)

					/* assert */
					actualPath := fakeOS.CreateArgsForCall(0)
					Expect(actualPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))

				})
				Context("os.Create errs", func() {
					It("should return error", func() {

						/* arrange */
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							"dummyContainerFilePath": "dummyScopeName",
						}

						expectedErr := errors.New("dummyError")

						fakeOS := new(ios.Fake)
						fakeOS.CreateReturns(nil, expectedErr)

						objectUnderTest := _Files{
							os: fakeOS,
						}

						/* act */
						_, actualErr := objectUnderTest.Interpret(
							"dummyPkgPath",
							map[string]*model.Value{},
							providedSCGContainerCallFiles,
							"dummyScratchDirPath",
						)

						/* assert */
						Expect(actualErr).To(Equal(expectedErr))
					})
				})
				Context("os.Create doesn't err", func() {
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
						fakeOS.CreateReturns(expectedCopyWriter, err)

						objectUnderTest := _Files{
							io: fakeIO,
							os: fakeOS,
						}

						/* act */
						objectUnderTest.Interpret(
							"dummyPkgPath",
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
								"dummyPkgPath",
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
								"dummyPkgPath",
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
