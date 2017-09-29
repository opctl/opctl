package files

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iio"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Context("Files", func() {
	tempFile, err := ioutil.TempFile("", "")
	if nil != err {
		panic(err)
	}

	tempFileInfo, err := os.Stat(tempFile.Name())
	if nil != err {
		panic(err)
	}
	Context("Interpret", func() {
		Context("bound value is absolute path", func() {
			It("should call expression.EvalToFile w/ expected args", func() {
				/* arrange */

				containerFilePath := "/dummyFile1Path.txt"

				providedSCGContainerCallFiles := map[string]string{
					// implicitly bound
					containerFilePath: "",
				}
				providedPkgHandle := new(pkg.FakeHandle)
				providedScope := map[string]*model.Value{}
				providedScratchDir := "dummyScratchDir"

				fakeExpression := new(expression.Fake)
				// error to trigger immediate return
				fakeExpression.EvalToFileReturns(nil, errors.New("dummyError"))

				objectUnderTest := _Files{
					expression: fakeExpression,
				}

				/* act */
				objectUnderTest.Interpret(
					providedPkgHandle,
					providedScope,
					providedSCGContainerCallFiles,
					providedScratchDir,
				)

				/* assert */
				actualScope,
					actualExpression,
					actualPkgHandle,
					actualScratchDir := fakeExpression.EvalToFileArgsForCall(0)
				Expect(actualScope).To(Equal(providedScope))
				Expect(actualExpression).To(Equal(fmt.Sprintf("$(%v)", containerFilePath)))
				Expect(actualPkgHandle).To(Equal(providedPkgHandle))
				Expect(actualScratchDir).To(Equal(providedScratchDir))
			})
			Context("expression.EvalToFile errs", func() {
				It("should return expected error", func() {
					/* arrange */
					containerFilePath := "/dummyFile1Path.txt"
					providedSCGContainerCallFiles := map[string]string{
						// implicitly bound
						containerFilePath: "",
					}

					getContentErr := fmt.Errorf("dummyError")

					fakeExpression := new(expression.Fake)
					fakeExpression.EvalToFileReturns(nil, getContentErr)

					expectedErr := fmt.Errorf(
						"unable to bind %v to %v; error was %v",
						containerFilePath,
						fmt.Sprintf("$(%v)", containerFilePath),
						getContentErr,
					)

					objectUnderTest := _Files{
						expression: fakeExpression,
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
			Context("expression.EvalToFile doesn't err", func() {
				Context("value.File not prefixed by rootFSPath", func() {
					It("should return expected results", func() {
						/* arrange */
						containerFilePath := "/dummyFile1Path.txt"

						fakeExpression := new(expression.Fake)
						filePath := tempFile.Name()
						fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

						expectedDCGContainerCallFiles := map[string]string{
							containerFilePath: filePath,
						}

						objectUnderTest := _Files{
							expression: fakeExpression,
							rootFSPath: "dummyRootFSPath",
						}

						/* act */
						actualDCGContainerCallFiles, actualErr := objectUnderTest.Interpret(
							new(pkg.FakeHandle),
							map[string]*model.Value{},
							map[string]string{
								// implicitly bound
								containerFilePath: "",
							},
							"dummyScratchDir",
						)

						/* assert */
						Expect(actualDCGContainerCallFiles).To(Equal(expectedDCGContainerCallFiles))
						Expect(actualErr).To(BeNil())

					})
				})
				Context("value.File prefixed by rootFSPath", func() {
					It("should call os.MkdirAll w/ expected args", func() {
						/* arrange */
						containerFilePath := "/parent/child/dummyFilePath.txt"
						providedScratchDirPath := "dummyScratchDirPath"

						fakeExpression := new(expression.Fake)
						filePath := tempFile.Name()
						fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

						fakeOS := new(ios.Fake)

						// err to trigger immediate return
						fakeOS.MkdirAllReturns(errors.New("dummyError"))

						expectedPath := filepath.Join(providedScratchDirPath, filepath.Dir(containerFilePath))

						objectUnderTest := _Files{
							expression: fakeExpression,
							os:         fakeOS,
						}

						/* act */
						objectUnderTest.Interpret(
							new(pkg.FakeHandle),
							map[string]*model.Value{},
							map[string]string{
								// implicitly bound
								containerFilePath: "",
							},
							providedScratchDirPath,
						)

						/* assert */
						actualPath,
							actualFileMode := fakeOS.MkdirAllArgsForCall(0)

						Expect(actualPath).To(Equal(expectedPath))
						Expect(actualFileMode).To(Equal(os.FileMode(0777)))

					})
					Context("os.MkdirAll errs", func() {
						It("should return expected error", func() {
							/* arrange */
							containerFilePath := "/dummyFile1Path.txt"

							fakeExpression := new(expression.Fake)
							filePath := tempFile.Name()
							fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

							fakeOS := new(ios.Fake)

							mkdirAllErr := fmt.Errorf("dummyMkdirAllError")
							fakeOS.MkdirAllReturns(mkdirAllErr)

							expectedErr := fmt.Errorf(
								"unable to bind %v to %v; error was %v",
								containerFilePath,
								fmt.Sprintf("$(%v)", containerFilePath),
								mkdirAllErr,
							)

							objectUnderTest := _Files{
								expression: fakeExpression,
								os:         fakeOS,
							}

							/* act */
							_, actualErr := objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								map[string]*model.Value{},
								map[string]string{
									// implicitly bound
									containerFilePath: "",
								},
								"dummyScratchDirPath",
							)

							/* assert */
							Expect(actualErr).To(Equal(expectedErr))
						})
					})
					Context("os.MkdirAll doesn't err", func() {
						It("should call os.Open w/ expected args", func() {
							/* arrange */
							containerFilePath := "/dummyFile1Path.txt"

							fakeExpression := new(expression.Fake)
							filePath := tempFile.Name()
							fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

							fakeOS := new(ios.Fake)

							// err to trigger immediate return
							fakeOS.OpenReturns(nil, errors.New("dummyError"))

							objectUnderTest := _Files{
								expression: fakeExpression,
								os:         fakeOS,
							}

							/* act */
							objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								map[string]*model.Value{},
								map[string]string{
									// implicitly bound
									containerFilePath: "",
								},
								"dummyScratchDir",
							)

							/* assert */
							actualPath := fakeOS.OpenArgsForCall(0)

							Expect(actualPath).To(Equal(filePath))

						})
						Context("os.Open errs", func() {
							It("should return expected error", func() {
								/* arrange */
								containerFilePath := "/dummyFile1Path.txt"

								fakeExpression := new(expression.Fake)
								filePath := tempFile.Name()
								fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

								fakeOS := new(ios.Fake)

								openError := fmt.Errorf("dummyOpenError")
								fakeOS.OpenReturns(nil, openError)

								expectedErr := fmt.Errorf(
									"unable to bind %v to %v; error was %v",
									containerFilePath,
									fmt.Sprintf("$(%v)", containerFilePath),
									openError,
								)

								objectUnderTest := _Files{
									expression: fakeExpression,
									os:         fakeOS,
								}

								/* act */
								_, actualErr := objectUnderTest.Interpret(
									new(pkg.FakeHandle),
									map[string]*model.Value{},
									map[string]string{
										// implicitly bound
										containerFilePath: "",
									},
									"dummyScratchDirPath",
								)

								/* assert */
								Expect(actualErr).To(Equal(expectedErr))
							})
						})
						Context("os.Open doesn't err", func() {
							It("should call os.Stat w/ expected args", func() {
								/* arrange */
								containerFilePath := "/dummyFile1Path.txt"

								fakeExpression := new(expression.Fake)
								filePath := tempFile.Name()
								fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

								fakeOS := new(ios.Fake)
								fakeOS.OpenReturns(tempFile, nil)

								// err to trigger immediate return
								fakeOS.StatReturns(nil, errors.New("dummyError"))

								objectUnderTest := _Files{
									expression: fakeExpression,
									os:         fakeOS,
								}

								/* act */
								objectUnderTest.Interpret(
									new(pkg.FakeHandle),
									map[string]*model.Value{},
									map[string]string{
										// implicitly bound
										containerFilePath: "",
									},
									"dummyScratchDir",
								)

								/* assert */
								actualPath := fakeOS.StatArgsForCall(0)

								Expect(actualPath).To(Equal(filePath))

							})
							Context("os.Stat errs", func() {
								It("should return expected error", func() {
									/* arrange */
									containerFilePath := "/dummyFile1Path.txt"

									fakeExpression := new(expression.Fake)
									filePath := tempFile.Name()
									fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

									fakeOS := new(ios.Fake)
									fakeOS.OpenReturns(tempFile, nil)

									statErr := fmt.Errorf("dummyStatError")
									fakeOS.StatReturns(nil, statErr)

									expectedErr := fmt.Errorf(
										"unable to bind %v to %v; error was %v",
										containerFilePath,
										fmt.Sprintf("$(%v)", containerFilePath),
										statErr,
									)

									objectUnderTest := _Files{
										expression: fakeExpression,
										os:         fakeOS,
									}

									/* act */
									_, actualErr := objectUnderTest.Interpret(
										new(pkg.FakeHandle),
										map[string]*model.Value{},
										map[string]string{
											// implicitly bound
											containerFilePath: "",
										},
										"dummyScratchDirPath",
									)

									/* assert */
									Expect(actualErr).To(Equal(expectedErr))
								})
							})
							Context("os.Stat doesn't err", func() {
								It("should call os.OpenFile w/ expected args", func() {
									/* arrange */
									containerFilePath := "/dummyFile1Path.txt"
									dummyScratchDir := "dummyScratchDir"

									fakeExpression := new(expression.Fake)
									fakeExpression.EvalToFileReturns(&model.Value{File: new(string)}, nil)

									fakeOS := new(ios.Fake)
									fakeOS.OpenReturns(tempFile, nil)
									fakeOS.StatReturns(tempFileInfo, nil)

									// err to trigger immediate return
									fakeOS.OpenFileReturns(tempFile, errors.New("dummyError"))

									expectedPath := filepath.Join(dummyScratchDir, containerFilePath)

									objectUnderTest := _Files{
										expression: fakeExpression,
										os:         fakeOS,
									}

									/* act */
									objectUnderTest.Interpret(
										new(pkg.FakeHandle),
										map[string]*model.Value{},
										map[string]string{
											// implicitly bound
											containerFilePath: "",
										},
										"dummyScratchDir",
									)

									/* assert */
									actualPath,
										actualFlags,
										actualPerm := fakeOS.OpenFileArgsForCall(0)

									Expect(actualPath).To(Equal(expectedPath))
									Expect(actualFlags).To(Equal(os.O_RDWR | os.O_CREATE))
									Expect(actualPerm).To(Equal(tempFileInfo.Mode()))

								})
								Context("os.OpenFile errs", func() {
									It("should return expected err", func() {
										/* arrange */
										containerFilePath := "/dummyFile1Path.txt"
										providedSCGContainerCallFiles := map[string]string{
											// implicitly bound
											containerFilePath: "",
										}

										fakeExpression := new(expression.Fake)
										filePath := "dummyFilePath"
										fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

										openFileErr := fmt.Errorf("dummyError")

										fakeOS := new(ios.Fake)
										fakeOS.OpenReturns(tempFile, nil)
										fakeOS.StatReturns(tempFileInfo, nil)
										fakeOS.OpenFileReturns(tempFile, openFileErr)

										expectedErr := fmt.Errorf(
											"unable to bind %v to %v; error was %v",
											containerFilePath,
											fmt.Sprintf("$(%v)", containerFilePath),
											openFileErr,
										)

										objectUnderTest := _Files{
											expression: fakeExpression,
											os:         fakeOS,
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
										containerFilePath := "/dummyFile1Path.txt"
										providedSCGContainerCallFiles := map[string]string{
											// implicitly bound
											containerFilePath: "",
										}

										fakeExpression := new(expression.Fake)
										filePath := "dummyFilePath"
										fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

										fakeOS := new(ios.Fake)
										fakeOS.OpenReturns(tempFile, nil)
										fakeOS.StatReturns(tempFileInfo, nil)
										fakeOS.OpenFileReturns(tempFile, nil)

										fakeIO := new(iio.Fake)
										// err to trigger immediate return
										fakeIO.CopyReturns(0, errors.New("dummyErr"))

										objectUnderTest := _Files{
											expression: fakeExpression,
											os:         fakeOS,
											io:         fakeIO,
										}

										/* act */
										objectUnderTest.Interpret(
											new(pkg.FakeHandle),
											map[string]*model.Value{},
											providedSCGContainerCallFiles,
											"dummyScratchDirPath",
										)

										/* assert */
										actualWriter, actualReader := fakeIO.CopyArgsForCall(0)
										Expect(actualWriter).To(Equal(tempFile))
										Expect(actualReader).To(Equal(tempFile))
									})
									Context("io.Copy errs", func() {
										It("should return expected err", func() {
											/* arrange */
											containerFilePath := "/dummyFile1Path.txt"
											providedSCGContainerCallFiles := map[string]string{
												// implicitly bound
												containerFilePath: "",
											}

											fakeExpression := new(expression.Fake)
											filePath := "dummyFilePath"
											fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

											fakeOS := new(ios.Fake)
											fakeOS.OpenReturns(tempFile, nil)
											fakeOS.StatReturns(tempFileInfo, nil)
											fakeOS.OpenFileReturns(tempFile, nil)

											copyErr := fmt.Errorf("dummyError")

											fakeIO := new(iio.Fake)
											fakeIO.CopyReturns(0, copyErr)

											expectedErr := fmt.Errorf(
												"unable to bind %v to %v; error was %v",
												containerFilePath,
												fmt.Sprintf("$(%v)", containerFilePath),
												copyErr,
											)

											objectUnderTest := _Files{
												expression: fakeExpression,
												os:         fakeOS,
												io:         fakeIO,
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
								})
							})
						})
					})
				})
			})
		})
	})
})
