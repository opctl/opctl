package files

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/filecopier"
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
	Context("Interpret", func() {
		It("should call expression.EvalToFile w/ expected args", func() {
			/* arrange */

			containerFilePath := "/dummyFile1Path.txt"

			providedSCGContainerCallFiles := map[string]interface{}{
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
				providedSCGContainerCallFiles := map[string]interface{}{
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
						map[string]interface{}{
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
						map[string]interface{}{
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
							map[string]interface{}{
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
					It("should call filecopier.OS w/ expected args", func() {
						/* arrange */
						providedScratchDir := "dummyScratchDir"
						containerFilePath := "/dummyFile1Path.txt"

						fakeExpression := new(expression.Fake)
						filePath := tempFile.Name()
						fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

						expectedPath := filepath.Join(providedScratchDir, containerFilePath)

						fakeFileCopier := new(filecopier.Fake)

						// err to trigger immediate return
						fakeFileCopier.OSReturns(errors.New("dummyError"))

						objectUnderTest := _Files{
							expression: fakeExpression,
							fileCopier: fakeFileCopier,
							os:         new(ios.Fake),
						}

						/* act */
						objectUnderTest.Interpret(
							new(pkg.FakeHandle),
							map[string]*model.Value{},
							map[string]interface{}{
								// implicitly bound
								containerFilePath: "",
							},
							providedScratchDir,
						)

						/* assert */
						actualSrcPath,
							actualDstPath := fakeFileCopier.OSArgsForCall(0)

						Expect(actualSrcPath).To(Equal(filePath))
						Expect(actualDstPath).To(Equal(expectedPath))

					})
					Context("filecopier.OS errs", func() {
						It("should return expected error", func() {
							/* arrange */
							containerFilePath := "/dummyFile1Path.txt"

							fakeExpression := new(expression.Fake)
							filePath := tempFile.Name()
							fakeExpression.EvalToFileReturns(&model.Value{File: &filePath}, nil)

							fakeFileCopier := new(filecopier.Fake)

							copyError := fmt.Errorf("dummyCopyError")

							// err to trigger immediate return
							fakeFileCopier.OSReturns(copyError)

							expectedErr := fmt.Errorf(
								"unable to bind %v to %v; error was %v",
								containerFilePath,
								fmt.Sprintf("$(%v)", containerFilePath),
								copyError,
							)

							objectUnderTest := _Files{
								expression: fakeExpression,
								fileCopier: fakeFileCopier,
								os:         new(ios.Fake),
							}

							/* act */
							_, actualErr := objectUnderTest.Interpret(
								new(pkg.FakeHandle),
								map[string]*model.Value{},
								map[string]interface{}{
									// implicitly bound
									containerFilePath: "",
								},
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
