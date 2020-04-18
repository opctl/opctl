package dirs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/golang-interfaces/ios"

	"github.com/golang-utils/dircopier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
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
				providedScope,
				providedSCGContainerCallDirs,
				providedScratchDir,
			)

			/* assert */
			actualScope,
				actualExpression,
				actualScratchDir,
				actualCreateIfNotExists := fakeDirInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualExpression).To(Equal(fmt.Sprintf("$(%v)", containerDirPath)))
			Expect(actualScratchDir).To(Equal(providedScratchDir))
			Expect(actualCreateIfNotExists).To(BeTrue())
		})
		Context("dirInterpreter.Interpret errs", func() {
			It("should return expected error", func() {
				/* arrange */
				containerDirPath := "/dummyDir1Path.txt"

				fakeDirInterpreter := new(dirFakes.FakeInterpreter)

				interpretError := fmt.Errorf("dummyCopyError")
				fakeDirInterpreter.InterpretReturns(nil, interpretError)

				expectedErr := fmt.Errorf(
					"unable to bind %v to %v; error was %v",
					containerDirPath,
					fmt.Sprintf("$(%v)", containerDirPath),
					interpretError,
				)

				objectUnderTest := _interpreter{
					dirInterpreter: fakeDirInterpreter,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
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
