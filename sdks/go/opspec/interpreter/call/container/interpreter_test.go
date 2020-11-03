package container

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	cmdFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/cmd/fakes"
	dirsFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/dirs/fakes"
	envvarsFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/envvars/fakes"
	filesFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/files/fakes"
	imageFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/image/fakes"
	socketsFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container/sockets/fakes"
	strFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/str/fakes"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter("")).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("calls os.MkdirAll w/ expected scratchdir path & returns error", func() {
			/* arrange */
			dataDirPath := "/dummydataDirPath"
			providedContainerID := "dummyContainerID"

			expectedScratchDirPath := filepath.Join(
				dataDirPath,
				"dcg",
				providedContainerID,
				"fs",
			)
			expectedScratchDirMode := os.FileMode(0700)

			expectedError := errors.New("dummyError")

			fakeOS := new(ios.Fake)
			// error to trigger immediate return
			fakeOS.MkdirAllReturns(expectedError)

			objectUnderTest := _interpreter{
				os:          fakeOS,
				dataDirPath: dataDirPath,
			}

			/* act */
			_, actualError := objectUnderTest.Interpret(
				map[string]*model.Value{},
				&model.ContainerCallSpec{},
				providedContainerID,
				"dummyOpPath",
			)

			/* assert */
			actualScratchDirPath, actualScratchDirMode := fakeOS.MkdirAllArgsForCall(0)
			Expect(actualScratchDirPath).To(Equal(expectedScratchDirPath))
			Expect(actualScratchDirMode).To(Equal(expectedScratchDirMode))
			Expect(actualError).To(Equal(expectedError))
		})

		It("should call cmdFakesInterpret w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}

			providedContainerCallSpec := &model.ContainerCallSpec{
				Cmd: []interface{}{
					"cmd",
				},
			}

			fakeCmdInterpreter := new(cmdFakes.FakeInterpreter)

			objectUnderTest := _interpreter{
				cmdInterpreter:     fakeCmdInterpreter,
				dirsInterpreter:    new(dirsFakes.FakeInterpreter),
				envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
				filesInterpreter:   new(filesFakes.FakeInterpreter),
				imageInterpreter:   new(imageFakes.FakeInterpreter),
				os:                 new(ios.Fake),
				socketsInterpreter: new(socketsFakes.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedContainerCallSpec,
				"dummyContainerID",
				"dummyOpPath",
			)

			/* assert */
			actualScope,
				actualScgContainerCallCmd := fakeCmdInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallCmd).To(Equal(providedContainerCallSpec.Cmd))

		})
		Context("cmdFakesInterpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeCmdInterpreter := new(cmdFakes.FakeInterpreter)
				fakeCmdInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					cmdInterpreter: fakeCmdInterpreter,
					os:             new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("cmdFakesInterpret doesn't error", func() {
			It("should return expected containerCall.Cmd", func() {
				/* arrange */
				expectedContainerCallCmd := []string{
					"cmd",
				}

				fakeCmdInterpreter := new(cmdFakes.FakeInterpreter)
				fakeCmdInterpreter.InterpretReturns(expectedContainerCallCmd, nil)

				objectUnderTest := _interpreter{
					cmdInterpreter:     fakeCmdInterpreter,
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(socketsFakes.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualResult.Cmd).To(Equal(expectedContainerCallCmd))
			})
		})

		It("should call dirsFakes.Interpret w/ expected args", func() {
			/* arrange */
			containerDirBind := "dummyContainerDirBind"

			providedScope := map[string]*model.Value{
				containerDirBind: {String: new(string)},
			}

			dirName := "dummyDirName"
			providedContainerCallSpec := &model.ContainerCallSpec{
				Dirs: map[string]string{
					// implicitly bound
					dirName: "",
				},
			}

			provideddataDirPath := "dummydataDirPath"
			providedContainerID := "dummyContainerID"

			expectedScratchDirPath := filepath.Join(
				provideddataDirPath,
				"dcg",
				providedContainerID,
				"fs",
			)

			fakeDirsInterpreter := new(dirsFakes.FakeInterpreter)

			objectUnderTest := _interpreter{
				cmdInterpreter:     new(cmdFakes.FakeInterpreter),
				dirsInterpreter:    fakeDirsInterpreter,
				envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
				filesInterpreter:   new(filesFakes.FakeInterpreter),
				imageInterpreter:   new(imageFakes.FakeInterpreter),
				os:                 new(ios.Fake),
				dataDirPath:        provideddataDirPath,
				socketsInterpreter: new(socketsFakes.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedContainerCallSpec,
				providedContainerID,
				"dummyOpPath",
			)

			/* assert */
			actualScope, actualScgContainerCallDirs, actualScratchDir := fakeDirsInterpreter.InterpretArgsForCall(0)
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallDirs).To(Equal(providedContainerCallSpec.Dirs))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("dirsFakes.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeDirsInterpreter := new(dirsFakes.FakeInterpreter)
				fakeDirsInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    fakeDirsInterpreter,
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("dirsFakes.Interpret doesn't error", func() {
			It("should return expected containerCall.Dirs", func() {
				/* arrange */
				expectedContainerCallDirs := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeDirsInterpreter := new(dirsFakes.FakeInterpreter)
				fakeDirsInterpreter.InterpretReturns(expectedContainerCallDirs, nil)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    fakeDirsInterpreter,
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(socketsFakes.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualResult.Dirs).To(Equal(expectedContainerCallDirs))
			})
		})

		It("should call envVars.Interpret w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}

			envVarName := "dummyEnvVarName"
			providedContainerCallSpec := &model.ContainerCallSpec{
				EnvVars: map[string]interface{}{
					// implicitly bound
					envVarName: "",
				},
			}

			fakeEnvVarsInterpreter := new(envvarsFakes.FakeInterpreter)

			objectUnderTest := _interpreter{
				cmdInterpreter:     new(cmdFakes.FakeInterpreter),
				dirsInterpreter:    new(dirsFakes.FakeInterpreter),
				envVarsInterpreter: fakeEnvVarsInterpreter,
				filesInterpreter:   new(filesFakes.FakeInterpreter),
				imageInterpreter:   new(imageFakes.FakeInterpreter),
				os:                 new(ios.Fake),
				socketsInterpreter: new(socketsFakes.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedContainerCallSpec,
				"dummyContainerID",
				"dummyOpPath",
			)

			/* assert */
			actualScope,
				actualScgContainerCallEnvVars := fakeEnvVarsInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallEnvVars).To(Equal(providedContainerCallSpec.EnvVars))

		})
		Context("envVars.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeEnvVarsInterpreter := new(envvarsFakes.FakeInterpreter)
				fakeEnvVarsInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: fakeEnvVarsInterpreter,
					os:                 new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("envVars.Interpret doesn't error", func() {
			It("should return expected containerCall.EnvVars", func() {
				/* arrange */
				expectedContainerCallEnvVars := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeEnvVarsInterpreter := new(envvarsFakes.FakeInterpreter)
				fakeEnvVarsInterpreter.InterpretReturns(expectedContainerCallEnvVars, nil)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: fakeEnvVarsInterpreter,
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(socketsFakes.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualResult.EnvVars).To(Equal(expectedContainerCallEnvVars))
			})
		})

		It("should call filesFakes.Interpret w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}

			fileName := "dummyFileName"
			providedContainerCallSpec := &model.ContainerCallSpec{
				Files: map[string]interface{}{
					// implicitly bound
					fileName: "",
				},
			}

			provideddataDirPath := "dummydataDirPath"
			providedContainerID := "dummyContainerID"

			expectedScratchDirPath := filepath.Join(
				provideddataDirPath,
				"dcg",
				providedContainerID,
				"fs",
			)

			fakeFilesInterpreter := new(filesFakes.FakeInterpreter)

			objectUnderTest := _interpreter{
				cmdInterpreter:     new(cmdFakes.FakeInterpreter),
				dirsInterpreter:    new(dirsFakes.FakeInterpreter),
				envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
				filesInterpreter:   fakeFilesInterpreter,
				imageInterpreter:   new(imageFakes.FakeInterpreter),
				os:                 new(ios.Fake),
				dataDirPath:        provideddataDirPath,
				socketsInterpreter: new(socketsFakes.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedContainerCallSpec,
				providedContainerID,
				"dummyOpPath",
			)

			/* assert */
			actualScope, actualScgContainerCallFiles, actualScratchDir := fakeFilesInterpreter.InterpretArgsForCall(0)
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallFiles).To(Equal(providedContainerCallSpec.Files))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("filesFakes.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeFilesInterpreter := new(filesFakes.FakeInterpreter)
				fakeFilesInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   fakeFilesInterpreter,
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("filesFakes.Interpret doesn't error", func() {
			It("should return expected containerCall.Files", func() {
				/* arrange */
				expectedContainerCallFiles := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeFilesInterpreter := new(filesFakes.FakeInterpreter)
				fakeFilesInterpreter.InterpretReturns(expectedContainerCallFiles, nil)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   fakeFilesInterpreter,
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(socketsFakes.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualResult.Files).To(Equal(expectedContainerCallFiles))
			})
		})

		It("should call imageFakes.Interpret w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedDataDirPath := "providedDataDirPath"
			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}
			providedContainerCallSpec := &model.ContainerCallSpec{
				Image: &model.ContainerCallImageSpec{
					Ref: "dummyRef",
					PullCreds: &model.CredsSpec{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				},
			}

			providedContainerID := "providedContainerID"

			expectedScratchDir := filepath.Join(
				providedDataDirPath,
				"dcg",
				providedContainerID,
				"fs",
			)

			fakeImageInterpreter := new(imageFakes.FakeInterpreter)

			objectUnderTest := _interpreter{
				cmdInterpreter:     new(cmdFakes.FakeInterpreter),
				dataDirPath:        providedDataDirPath,
				dirsInterpreter:    new(dirsFakes.FakeInterpreter),
				envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
				filesInterpreter:   new(filesFakes.FakeInterpreter),
				imageInterpreter:   fakeImageInterpreter,
				os:                 new(ios.Fake),
				socketsInterpreter: new(socketsFakes.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedContainerCallSpec,
				providedContainerID,
				"dummyOpPath",
			)

			/* assert */
			actualScope,
				actualScgContainerCallImage,
				actualScratchDir := fakeImageInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallImage).To(Equal(providedContainerCallSpec.Image))
			Expect(actualScratchDir).To(Equal(expectedScratchDir))

		})
		Context("imageFakes.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeImageInterpreter := new(imageFakes.FakeInterpreter)
				fakeImageInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   fakeImageInterpreter,
					os:                 new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("imageFakes.Interpret doesn't error", func() {
			It("should return expected containerCall.Image", func() {
				/* arrange */
				expectedContainerCallImage := &model.ContainerCallImage{
					Ref: new(string),
					PullCreds: &model.Creds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				}

				fakeImageInterpreter := new(imageFakes.FakeInterpreter)
				fakeImageInterpreter.InterpretReturns(expectedContainerCallImage, nil)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   fakeImageInterpreter,
					os:                 new(ios.Fake),
					socketsInterpreter: new(socketsFakes.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualResult.Image).To(Equal(expectedContainerCallImage))
			})
		})

		Context("containerCallSpec.Name truthy", func() {
			It("should call stringInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				containerFileBind := "dummyContainerFileBind"

				providedScope := map[string]*model.Value{
					containerFileBind: {String: new(string)},
				}
				expectedScgContainerCallName := "name"
				providedContainerCallSpec := &model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{},
					Name:  &expectedScgContainerCallName,
				}

				fakeStrInterpreter := new(strFakes.FakeInterpreter)

				fakeStrInterpreter.InterpretReturns(&model.Value{String: new(string)}, nil)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(socketsFakes.FakeInterpreter),
					stringInterpreter:  fakeStrInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedContainerCallSpec,
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				actualScope,
					actualScgContainerCallName := fakeStrInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualScgContainerCallName).To(Equal(*providedContainerCallSpec.Name))

			})
			Context("stringInterpreter.Interpret errors", func() {
				It("should return expected error", func() {
					/* arrange */
					expectedErr := errors.New("dummyError")
					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					fakeStrInterpreter.InterpretReturns(nil, expectedErr)

					objectUnderTest := _interpreter{
						cmdInterpreter:     new(cmdFakes.FakeInterpreter),
						dirsInterpreter:    new(dirsFakes.FakeInterpreter),
						envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
						filesInterpreter:   new(filesFakes.FakeInterpreter),
						imageInterpreter:   new(imageFakes.FakeInterpreter),
						os:                 new(ios.Fake),
						stringInterpreter:  fakeStrInterpreter,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.ContainerCallSpec{
							Name: new(string),
						},
						"dummyContainerID",
						"dummyOpPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("stringInterpreter.Interpret doesn't error", func() {
				It("should return expected containerCall.Name", func() {
					/* arrange */
					expectedContainerCallName := "name"

					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					fakeStrInterpreter.InterpretReturns(&model.Value{String: &expectedContainerCallName}, nil)

					objectUnderTest := _interpreter{
						cmdInterpreter:     new(cmdFakes.FakeInterpreter),
						dirsInterpreter:    new(dirsFakes.FakeInterpreter),
						envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
						filesInterpreter:   new(filesFakes.FakeInterpreter),
						imageInterpreter:   new(imageFakes.FakeInterpreter),
						os:                 new(ios.Fake),
						socketsInterpreter: new(socketsFakes.FakeInterpreter),
						stringInterpreter:  fakeStrInterpreter,
					}

					/* act */
					actualResult, _ := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.ContainerCallSpec{
							Name: new(string),
						},
						"dummyContainerID",
						"dummyOpPath",
					)

					/* assert */
					Expect(*actualResult.Name).To(Equal(expectedContainerCallName))
				})
			})
		})

		Context("containerCallSpec.WorkDir truthy", func() {
			It("should call stringInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				containerFileBind := "dummyContainerFileBind"

				providedScope := map[string]*model.Value{
					containerFileBind: {String: new(string)},
				}
				providedContainerCallSpec := &model.ContainerCallSpec{
					Image:   &model.ContainerCallImageSpec{},
					WorkDir: "workDir",
				}

				fakeStrInterpreter := new(strFakes.FakeInterpreter)
				fakeStrInterpreter.InterpretReturns(&model.Value{String: new(string)}, nil)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(socketsFakes.FakeInterpreter),
					stringInterpreter:  fakeStrInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedScope,
					providedContainerCallSpec,
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				actualScope,
					actualScgContainerCallWorkDir := fakeStrInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualScgContainerCallWorkDir).To(Equal(providedContainerCallSpec.WorkDir))

			})
			Context("stringInterpreter.Interpret errors", func() {
				It("should return expected error", func() {
					/* arrange */
					expectedErr := errors.New("dummyError")
					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					fakeStrInterpreter.InterpretReturns(nil, expectedErr)

					objectUnderTest := _interpreter{
						cmdInterpreter:     new(cmdFakes.FakeInterpreter),
						dirsInterpreter:    new(dirsFakes.FakeInterpreter),
						envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
						filesInterpreter:   new(filesFakes.FakeInterpreter),
						imageInterpreter:   new(imageFakes.FakeInterpreter),
						os:                 new(ios.Fake),
						stringInterpreter:  fakeStrInterpreter,
					}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.ContainerCallSpec{
							WorkDir: "dummyWorkDir",
						},
						"dummyContainerID",
						"dummyOpPath",
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("stringInterpreter.Interpret doesn't error", func() {
				It("should return expected containerCall.WorkDir", func() {
					/* arrange */
					expectedContainerCallWorkDir := "workDir"

					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					fakeStrInterpreter.InterpretReturns(&model.Value{String: &expectedContainerCallWorkDir}, nil)

					objectUnderTest := _interpreter{
						cmdInterpreter:     new(cmdFakes.FakeInterpreter),
						dirsInterpreter:    new(dirsFakes.FakeInterpreter),
						envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
						filesInterpreter:   new(filesFakes.FakeInterpreter),
						imageInterpreter:   new(imageFakes.FakeInterpreter),
						os:                 new(ios.Fake),
						socketsInterpreter: new(socketsFakes.FakeInterpreter),
						stringInterpreter:  fakeStrInterpreter,
					}

					/* act */
					actualResult, _ := objectUnderTest.Interpret(
						map[string]*model.Value{},
						&model.ContainerCallSpec{
							WorkDir: "dummyWorkDir",
						},
						"dummyContainerID",
						"dummyOpPath",
					)

					/* assert */
					Expect(actualResult.WorkDir).To(Equal(expectedContainerCallWorkDir))
				})
			})
		})

		It("should call socketsFakes.Interpret w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}

			envVarName := "dummyEnvVarName"
			providedContainerCallSpec := &model.ContainerCallSpec{
				Sockets: map[string]string{
					// implicitly bound
					envVarName: "",
				},
			}

			provideddataDirPath := "dummydataDirPath"
			providedContainerID := "dummyContainerID"

			expectedScratchDirPath := filepath.Join(
				provideddataDirPath,
				"dcg",
				providedContainerID,
				"fs",
			)

			fakeSocketsInterpreter := new(socketsFakes.FakeInterpreter)

			objectUnderTest := _interpreter{
				cmdInterpreter:     new(cmdFakes.FakeInterpreter),
				dirsInterpreter:    new(dirsFakes.FakeInterpreter),
				envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
				filesInterpreter:   new(filesFakes.FakeInterpreter),
				imageInterpreter:   new(imageFakes.FakeInterpreter),
				os:                 new(ios.Fake),
				dataDirPath:        provideddataDirPath,
				socketsInterpreter: fakeSocketsInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedContainerCallSpec,
				providedContainerID,
				"dummyOpPath",
			)

			/* assert */
			actualScope, actualScgContainerCallSockets, actualScratchDir := fakeSocketsInterpreter.InterpretArgsForCall(0)
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallSockets).To(Equal(providedContainerCallSpec.Sockets))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("socketsFakes.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeSocketsInterpreter := new(socketsFakes.FakeInterpreter)
				fakeSocketsInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: fakeSocketsInterpreter,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("socketsFakes.Interpret doesn't error", func() {
			It("should return expected containerCall.Sockets", func() {
				/* arrange */
				expectedContainerCallSockets := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeSocketsInterpreter := new(socketsFakes.FakeInterpreter)
				fakeSocketsInterpreter.InterpretReturns(expectedContainerCallSockets, nil)

				objectUnderTest := _interpreter{
					cmdInterpreter:     new(cmdFakes.FakeInterpreter),
					dirsInterpreter:    new(dirsFakes.FakeInterpreter),
					envVarsInterpreter: new(envvarsFakes.FakeInterpreter),
					filesInterpreter:   new(filesFakes.FakeInterpreter),
					imageInterpreter:   new(imageFakes.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: fakeSocketsInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"dummyContainerID",
					"dummyOpPath",
				)

				/* assert */
				Expect(actualResult.Sockets).To(Equal(expectedContainerCallSockets))
			})
		})
	})
})
