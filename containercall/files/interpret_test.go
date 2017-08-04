package files

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iio"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
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
			It("should call pkgHandle.GetContent w/ expected args", func() {
				/* arrange */

				containerFilePath := "/dummyFile1Path.txt"

				providedSCGContainerCallFiles := map[string]string{
					// implicitly bound
					containerFilePath: "",
				}

				providedParentOpPkgHandle := new(pkg.FakeHandle)
				// error to trigger immediate return
				providedParentOpPkgHandle.GetContentReturns(nil, errors.New("dummyError"))

				objectUnderTest := _Files{}

				/* act */
				objectUnderTest.Interpret(
					providedParentOpPkgHandle,
					map[string]*model.Value{},
					providedSCGContainerCallFiles,
					"dummyScratchDir",
				)

				/* assert */
				actualContext, actualContentPath := providedParentOpPkgHandle.GetContentArgsForCall(0)
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

					providedParentOpPkgHandle := new(pkg.FakeHandle)
					providedParentOpPkgHandle.GetContentReturns(nil, getContentErr)

					expectedErr := fmt.Errorf(
						"Unable to bind file '%v' to pkg content '%v'; error was: %v",
						containerFilePath,
						containerFilePath,
						getContentErr,
					)

					objectUnderTest := _Files{}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						providedParentOpPkgHandle,
						map[string]*model.Value{},
						providedSCGContainerCallFiles,
						"dummyScratchDirPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("pkgHandle.GetContent doesn't err", func() {
				It("should call os.Open w/ expected args", func() {

				})
				Context("os.Open errs", func() {
					It("should return expected err", func() {
						/* arrange */
						containerFilePath := "/dummyFile1Path.txt"
						providedSCGContainerCallFiles := map[string]string{
							// implicitly bound
							containerFilePath: "",
						}

						openErr := fmt.Errorf("dummyError")

						fakeOS := new(ios.Fake)
						fakeOS.OpenReturns(nil, openErr)

						expectedErr := fmt.Errorf(
							"Unable to bind file '%v' to pkg content '%v'; error was: %v",
							containerFilePath,
							containerFilePath,
							openErr,
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
				Context("os.Open doesn't err", func() {
					It("should call io.Copy w/ expected args", func() {

					})
					Context("io.Copy errs", func() {
						It("should return expected err", func() {
							/* arrange */
							containerFilePath := "/dummyFile1Path.txt"
							providedSCGContainerCallFiles := map[string]string{
								// implicitly bound
								containerFilePath: "",
							}

							copyErr := fmt.Errorf("dummyError")

							fakeIO := new(iio.Fake)
							fakeIO.CopyReturns(0, copyErr)

							expectedErr := fmt.Errorf(
								"Unable to bind file '%v' to pkg content '%v'; error was: %v",
								containerFilePath,
								containerFilePath,
								copyErr,
							)

							objectUnderTest := _Files{
								os: new(ios.Fake),
								io: fakeIO,
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
						new(pkg.FakeHandle),
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

					expectedErr := fmt.Errorf(
						"Unable to bind file '%v' to '%v'; '%v' not a file, number, object, or string",
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
							new(pkg.FakeHandle),
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
								new(pkg.FakeHandle),
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
						new(pkg.FakeHandle),
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
							new(pkg.FakeHandle),
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
								new(pkg.FakeHandle),
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
			Context("value.Object not nil", func() {
				It("should call json.Marshal w/ expected args", func() {
					/* arrange */
					scopeName := "dummyScopeName"

					providedScope := map[string]*model.Value{
						scopeName: {Object: map[string]interface{}{}},
					}

					providedSCGContainerCallFiles := map[string]string{
						// explicitly bound
						"dummyContainerFilePath": scopeName,
					}

					providedScratchDirPath := "dummyScratchDirPath"

					fakeJSON := new(ijson.Fake)
					// error to trigger immediate return
					fakeJSON.MarshalReturns(nil, errors.New("dummyError"))

					objectUnderTest := _Files{
						json: fakeJSON,
					}

					/* act */
					objectUnderTest.Interpret(
						new(pkg.FakeHandle),
						providedScope,
						providedSCGContainerCallFiles,
						providedScratchDirPath,
					)

					/* assert */
					Expect(fakeJSON.MarshalArgsForCall(0)).To(Equal(providedScope[scopeName].Object))
				})
				Context("json.Marshal errs", func() {

					It("should return expected error", func() {
						/* arrange */
						scopeName := "dummyScopeName"

						providedScope := map[string]*model.Value{
							scopeName: {Object: map[string]interface{}{}},
						}

						containerFilePath := "dummyContainerFilePath"
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							containerFilePath: scopeName,
						}

						providedScratchDirPath := "dummyScratchDirPath"

						marshalErr := errors.New("dummyError")
						fakeJSON := new(ijson.Fake)
						// error to trigger immediate return
						fakeJSON.MarshalReturns(nil, marshalErr)

						expectedErr := fmt.Errorf(
							"Unable to bind file '%v' to %v; error was: %v",
							containerFilePath,
							scopeName,
							marshalErr.Error(),
						)

						objectUnderTest := _Files{
							json: fakeJSON,
						}

						/* act */
						_, actualErr := objectUnderTest.Interpret(
							new(pkg.FakeHandle),
							providedScope,
							providedSCGContainerCallFiles,
							providedScratchDirPath,
						)

						/* assert */
						Expect(actualErr.Error()).To(Equal(expectedErr.Error()))
					})
				})
				Context("json.Marshal doesn't err", func() {
					It("should call os.MkdirAll w/ expected args", func() {
						/* arrange */
						scopeName := "dummyScopeName"

						providedScope := map[string]*model.Value{
							scopeName: {Object: map[string]interface{}{}},
						}

						containerFilePath := "dummyContainerFilePath"
						providedSCGContainerCallFiles := map[string]string{
							// explicitly bound
							containerFilePath: scopeName,
						}

						providedScratchDirPath := "dummyScratchDirPath"

						fakeOS := new(ios.Fake)

						objectUnderTest := _Files{
							io:   new(iio.Fake),
							json: new(ijson.Fake),
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
						actualPath, actualFileMode := fakeOS.MkdirAllArgsForCall(0)
						Expect(actualPath).To(Equal(filepath.Dir(filepath.Join(providedScratchDirPath, containerFilePath))))
						Expect(actualFileMode).To(Equal(os.FileMode(0700)))

					})
					Context("os.MkdirAll errs", func() {
						It("should return error", func() {

							/* arrange */
							scopeName := "dummyScopeName"

							providedScope := map[string]*model.Value{
								scopeName: {Object: map[string]interface{}{}},
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
								json: new(ijson.Fake),
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
						It("should call os.Create w/ expected args", func() {
							/* arrange */
							scopeName := "dummyScopeName"

							providedScope := map[string]*model.Value{
								scopeName: {Object: map[string]interface{}{}},
							}

							containerFilePath := "dummyContainerFilePath"
							providedSCGContainerCallFiles := map[string]string{
								// explicitly bound
								containerFilePath: scopeName,
							}

							providedScratchDirPath := "dummyScratchDirPath"

							fakeOS := new(ios.Fake)

							objectUnderTest := _Files{
								io:   new(iio.Fake),
								json: new(ijson.Fake),
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
							actualPath := fakeOS.CreateArgsForCall(0)
							Expect(actualPath).To(Equal(filepath.Join(providedScratchDirPath, containerFilePath)))

						})
						Context("os.Create errs", func() {
							It("should return error", func() {

								/* arrange */
								scopeName := "dummyScopeName"

								providedScope := map[string]*model.Value{
									scopeName: {Object: map[string]interface{}{}},
								}

								providedSCGContainerCallFiles := map[string]string{
									// explicitly bound
									"dummyContainerFilePath": scopeName,
								}

								expectedErr := errors.New("dummyError")

								fakeOS := new(ios.Fake)
								fakeOS.CreateReturns(nil, expectedErr)

								objectUnderTest := _Files{
									json: new(ijson.Fake),
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
						Context("os.Create doesn't err", func() {
							It("should call io.Copy w/ expected args", func() {
								/* arrange */
								scopeName := "dummyScopeName"

								providedScope := map[string]*model.Value{
									scopeName: {Object: map[string]interface{}{}},
								}

								providedSCGContainerCallFiles := map[string]string{
									// explicitly bound
									"dummyContainerFilePath": scopeName,
								}

								marshalledObject := []byte("dummyMarshalledObject")
								fakeJSON := new(ijson.Fake)
								fakeJSON.MarshalReturns(marshalledObject, nil)

								providedScratchDirPath := "dummyScratchDirPath"
								expectedCopyReader := bytes.NewReader(marshalledObject)

								fakeIO := new(iio.Fake)

								fakeOS := new(ios.Fake)
								expectedCopyWriter, err := ioutil.TempFile("", "")
								fakeOS.CreateReturns(expectedCopyWriter, err)

								objectUnderTest := _Files{
									io:   fakeIO,
									json: fakeJSON,
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
										scopeName: {Object: map[string]interface{}{}},
									}

									providedSCGContainerCallFiles := map[string]string{
										// explicitly bound
										"dummyContainerFilePath": scopeName,
									}

									expectedErr := errors.New("dummyError")

									fakeIO := new(iio.Fake)
									fakeIO.CopyReturns(0, expectedErr)

									objectUnderTest := _Files{
										io:   fakeIO,
										json: new(ijson.Fake),
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
										scopeName: {Object: map[string]interface{}{}},
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
										io:   new(iio.Fake),
										json: new(ijson.Fake),
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

					expectedErr := fmt.Errorf(
						"Unable to bind file '%v' to '%v'; '%v' not a file, number, object, or string",
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
						new(pkg.FakeHandle),
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
							new(pkg.FakeHandle),
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
								new(pkg.FakeHandle),
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
						new(pkg.FakeHandle),
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
							new(pkg.FakeHandle),
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
