package dirs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/ios"

	"github.com/golang-utils/dircopier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	dirFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/dir/fakes"
)

var _ = Context("Dirs", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter("")).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		tempDir, err := ioutil.TempFile("", "")
		if nil != err {
			panic(err)
		}
		It("should call dirInterpreter.Interpret w/ expected args", func() {
			/* arrange */

			containerDirPath := "/dummyDir1Path.txt"

			providedSCGContainerCallDirs := map[string]string{
				// implicitly bound
				containerDirPath: "",
			}
			providedOpHandle := new(modelFakes.FakeDataHandle)
			providedScope := map[string]*model.Value{}
			providedScratchDir := "dummyScratchDir"

			fakeDirInterpreter := new(dirFakes.FakeInterpreter)
			// error to trigger immediate return
			fakeDirInterpreter.InterpretReturns(nil, errors.New("dummyError"))
			fakeOS := new(ios.Fake)
			fakeOS.MkdirAllReturns(errors.New("dummyError"))

			objectUnderTest := _interpreter{
				dirInterpreter: fakeDirInterpreter,
				os:             fakeOS,
			}

			/* act */
			objectUnderTest.Interpret(
				providedOpHandle,
				providedScope,
				providedSCGContainerCallDirs,
				providedScratchDir,
			)

			/* assert */
			actualScope,
				actualExpression,
				actualOpHandle := fakeDirInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(fmt.Sprintf("$(%v)", containerDirPath)))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
		})
		Context("dirInterpreter.Interpret errs", func() {
			It("should call os.MkdirAll w/ expected args", func() {
				/* arrange */
				providedScratchDirPath := "providedScratchDirPath"

				providedContainerDirPath := "/dummyDir1Path.txt"
				providedSCGContainerCallDirs := map[string]string{
					// implicitly bound
					providedContainerDirPath: "",
				}

				fakeDirInterpreter := new(dirFakes.FakeInterpreter)
				fakeDirInterpreter.InterpretReturns(nil, errors.New("dummyError"))

				expectedPath := filepath.Join(providedScratchDirPath, providedContainerDirPath)

				fakeOS := new(ios.Fake)
				// error to trigger immediate return
				fakeOS.MkdirAllReturns(errors.New("dummyError"))

				objectUnderTest := _interpreter{
					dirInterpreter: fakeDirInterpreter,
					os:             fakeOS,
				}

				/* act */
				objectUnderTest.Interpret(
					new(modelFakes.FakeDataHandle),
					map[string]*model.Value{},
					providedSCGContainerCallDirs,
					providedScratchDirPath,
				)

				/* assert */
				actualPath,
					actualMode := fakeOS.MkdirAllArgsForCall(0)

				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualMode).To(Equal(os.FileMode(0700)))
			})
		})
		Context("dirInterpreter.Interpret doesn't err", func() {
			Context("value.Dir not prefixed by dataDirPath", func() {
				It("should return expected results", func() {
					/* arrange */
					containerDirPath := "/dummyDir1Path.txt"

					fakeDirInterpreter := new(dirFakes.FakeInterpreter)
					filePath := tempDir.Name()
					fakeDirInterpreter.InterpretReturns(&model.Value{Dir: &filePath}, nil)

					expectedDCGContainerCallDirs := map[string]string{
						containerDirPath: filePath,
					}

					objectUnderTest := _interpreter{
						dirInterpreter: fakeDirInterpreter,
						dataDirPath:    "dummydataDirPath",
					}

					/* act */
					actualDCGContainerCallDirs, actualErr := objectUnderTest.Interpret(
						new(modelFakes.FakeDataHandle),
						map[string]*model.Value{},
						map[string]string{
							// implicitly bound
							containerDirPath: "",
						},
						"dummyScratchDir",
					)

					/* assert */
					Expect(actualDCGContainerCallDirs).To(Equal(expectedDCGContainerCallDirs))
					Expect(actualErr).To(BeNil())

				})
			})
			Context("value.Dir prefixed by dataDirPath", func() {
				It("should call dircopier.OS w/ expected args", func() {
					/* arrange */
					providedScratchDir := "dummyScratchDir"
					containerDirPath := "/dummyDir1Path.txt"

					fakeDirInterpreter := new(dirFakes.FakeInterpreter)
					filePath := tempDir.Name()
					fakeDirInterpreter.InterpretReturns(&model.Value{Dir: &filePath}, nil)

					expectedPath := filepath.Join(providedScratchDir, containerDirPath)

					fakeDirCopier := new(dircopier.Fake)

					// err to trigger immediate return
					fakeDirCopier.OSReturns(errors.New("dummyError"))

					objectUnderTest := _interpreter{
						dirInterpreter: fakeDirInterpreter,
						dirCopier:      fakeDirCopier,
					}

					/* act */
					objectUnderTest.Interpret(
						new(modelFakes.FakeDataHandle),
						map[string]*model.Value{},
						map[string]string{
							// implicitly bound
							containerDirPath: "",
						},
						providedScratchDir,
					)

					/* assert */
					actualSrcPath,
						actualDstPath := fakeDirCopier.OSArgsForCall(0)

					Expect(actualSrcPath).To(Equal(filePath))
					Expect(actualDstPath).To(Equal(expectedPath))

				})
				Context("dircopier.OS errs", func() {
					It("should return expected error", func() {
						/* arrange */
						containerDirPath := "/dummyDir1Path.txt"

						fakeDirInterpreter := new(dirFakes.FakeInterpreter)
						filePath := tempDir.Name()
						fakeDirInterpreter.InterpretReturns(&model.Value{Dir: &filePath}, nil)

						fakeDirCopier := new(dircopier.Fake)

						copyError := fmt.Errorf("dummyCopyError")

						// err to trigger immediate return
						fakeDirCopier.OSReturns(copyError)

						expectedErr := fmt.Errorf(
							"unable to bind %v to %v; error was %v",
							containerDirPath,
							fmt.Sprintf("$(%v)", containerDirPath),
							copyError,
						)

						objectUnderTest := _interpreter{
							dirInterpreter: fakeDirInterpreter,
							dirCopier:      fakeDirCopier,
						}

						/* act */
						_, actualErr := objectUnderTest.Interpret(
							new(modelFakes.FakeDataHandle),
							map[string]*model.Value{},
							map[string]string{
								// implicitly bound
								containerDirPath: "",
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
