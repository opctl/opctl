package container

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
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
			providedRootOpID := "dummyRootOpID"

			expectedScratchDirPath := filepath.Join(
				dataDirPath,
				"dcg",
				providedRootOpID,
				"containers",
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
				&model.SCGContainerCall{},
				providedContainerID,
				providedRootOpID,
				new(modelFakes.FakeDataHandle),
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

			providedSCGContainerCall := &model.SCGContainerCall{
				Cmd: []interface{}{
					"cmd",
				},
			}

			providedOpHandle := new(modelFakes.FakeDataHandle)

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
				providedSCGContainerCall,
				"dummyContainerID",
				"dummyRootOpID",
				providedOpHandle,
			)

			/* assert */
			actualScope,
				actualScgContainerCallCmd,
				actualOpHandle := fakeCmdInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallCmd).To(Equal(providedSCGContainerCall.Cmd))
			Expect(actualOpHandle).To(Equal(providedOpHandle))

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("cmdFakesInterpret doesn't error", func() {
			It("should return expected dcgContainerCall.Cmd", func() {
				/* arrange */
				expectedDCGContainerCallCmd := []string{
					"cmd",
				}

				fakeCmdInterpreter := new(cmdFakes.FakeInterpreter)
				fakeCmdInterpreter.InterpretReturns(expectedDCGContainerCallCmd, nil)

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualResult.Cmd).To(Equal(expectedDCGContainerCallCmd))
			})
		})

		It("should call dirsFakes.Interpret w/ expected args", func() {
			/* arrange */
			containerDirBind := "dummyContainerDirBind"

			providedScope := map[string]*model.Value{
				containerDirBind: {String: new(string)},
			}

			dirName := "dummyDirName"
			providedSCGContainerCall := &model.SCGContainerCall{
				Dirs: map[string]string{
					// implicitly bound
					dirName: "",
				},
			}

			provideddataDirPath := "dummydataDirPath"
			providedContainerID := "dummyContainerID"
			providedRootOpID := "dummyRootOpID"
			providedOpHandle := new(modelFakes.FakeDataHandle)

			expectedScratchDirPath := filepath.Join(
				provideddataDirPath,
				"dcg",
				providedRootOpID,
				"containers",
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
				providedSCGContainerCall,
				providedContainerID,
				providedRootOpID,
				providedOpHandle,
			)

			/* assert */
			actualOpHandle, actualScope, actualScgContainerCallDirs, actualScratchDir := fakeDirsInterpreter.InterpretArgsForCall(0)
			Expect(actualOpHandle).To(Equal(providedOpHandle))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallDirs).To(Equal(providedSCGContainerCall.Dirs))
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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("dirsFakes.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.Dirs", func() {
				/* arrange */
				expectedDCGContainerCallDirs := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeDirsInterpreter := new(dirsFakes.FakeInterpreter)
				fakeDirsInterpreter.InterpretReturns(expectedDCGContainerCallDirs, nil)

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualResult.Dirs).To(Equal(expectedDCGContainerCallDirs))
			})
		})

		It("should call envVars.Interpret w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}

			envVarName := "dummyEnvVarName"
			providedSCGContainerCall := &model.SCGContainerCall{
				EnvVars: map[string]interface{}{
					// implicitly bound
					envVarName: "",
				},
			}

			providedOpHandle := new(modelFakes.FakeDataHandle)

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
				providedSCGContainerCall,
				"dummyContainerID",
				"dummyRootOpID",
				providedOpHandle,
			)

			/* assert */
			actualScope,
				actualScgContainerCallEnvVars,
				actualOpHandle := fakeEnvVarsInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallEnvVars).To(Equal(providedSCGContainerCall.EnvVars))
			Expect(actualOpHandle).To(Equal(providedOpHandle))

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("envVars.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.EnvVars", func() {
				/* arrange */
				expectedDCGContainerCallEnvVars := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeEnvVarsInterpreter := new(envvarsFakes.FakeInterpreter)
				fakeEnvVarsInterpreter.InterpretReturns(expectedDCGContainerCallEnvVars, nil)

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualResult.EnvVars).To(Equal(expectedDCGContainerCallEnvVars))
			})
		})

		It("should call filesFakes.Interpret w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}

			fileName := "dummyFileName"
			providedSCGContainerCall := &model.SCGContainerCall{
				Files: map[string]interface{}{
					// implicitly bound
					fileName: "",
				},
			}

			provideddataDirPath := "dummydataDirPath"
			providedContainerID := "dummyContainerID"
			providedRootOpID := "dummyRootOpID"
			providedOpHandle := new(modelFakes.FakeDataHandle)

			expectedScratchDirPath := filepath.Join(
				provideddataDirPath,
				"dcg",
				providedRootOpID,
				"containers",
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
				providedSCGContainerCall,
				providedContainerID,
				providedRootOpID,
				providedOpHandle,
			)

			/* assert */
			actualOpHandle, actualScope, actualScgContainerCallFiles, actualScratchDir := fakeFilesInterpreter.InterpretArgsForCall(0)
			Expect(actualOpHandle).To(Equal(providedOpHandle))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallFiles).To(Equal(providedSCGContainerCall.Files))
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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("filesFakes.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.Files", func() {
				/* arrange */
				expectedDCGContainerCallFiles := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeFilesInterpreter := new(filesFakes.FakeInterpreter)
				fakeFilesInterpreter.InterpretReturns(expectedDCGContainerCallFiles, nil)

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualResult.Files).To(Equal(expectedDCGContainerCallFiles))
			})
		})

		It("should call imageFakes.Interpret w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}
			providedSCGContainerCall := &model.SCGContainerCall{
				Image: &model.SCGContainerCallImage{
					Ref: new(string),
					PullCreds: &model.SCGPullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				},
			}

			providedOpHandle := new(modelFakes.FakeDataHandle)

			fakeImageInterpreter := new(imageFakes.FakeInterpreter)

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
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				"dummyContainerID",
				"dummyRootOpID",
				providedOpHandle,
			)

			/* assert */
			actualScope,
				actualScgContainerCallImage,
				actualOpHandle := fakeImageInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallImage).To(Equal(providedSCGContainerCall.Image))
			Expect(actualOpHandle).To(Equal(providedOpHandle))

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("imageFakes.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.Image", func() {
				/* arrange */
				expectedDCGContainerCallImage := &model.DCGContainerCallImage{
					Ref: new(string),
					PullCreds: &model.PullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				}

				fakeImageInterpreter := new(imageFakes.FakeInterpreter)
				fakeImageInterpreter.InterpretReturns(expectedDCGContainerCallImage, nil)

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualResult.Image).To(Equal(expectedDCGContainerCallImage))
			})
		})

		Context("scgContainerCall.Name truthy", func() {
			It("should call stringInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				containerFileBind := "dummyContainerFileBind"

				providedScope := map[string]*model.Value{
					containerFileBind: {String: new(string)},
				}
				expectedScgContainerCallName := "name"
				providedSCGContainerCall := &model.SCGContainerCall{
					Image: &model.SCGContainerCallImage{},
					Name:  &expectedScgContainerCallName,
				}

				providedOpHandle := new(modelFakes.FakeDataHandle)

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
					providedSCGContainerCall,
					"dummyContainerID",
					"dummyRootOpID",
					providedOpHandle,
				)

				/* assert */
				actualScope,
					actualScgContainerCallName,
					actualOpHandle := fakeStrInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualScgContainerCallName).To(Equal(*providedSCGContainerCall.Name))
				Expect(actualOpHandle).To(Equal(providedOpHandle))

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
						&model.SCGContainerCall{
							Name: new(string),
						},
						"dummyContainerID",
						"dummyRootOpID",
						new(modelFakes.FakeDataHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("stringInterpreter.Interpret doesn't error", func() {
				It("should return expected dcgContainerCall.Name", func() {
					/* arrange */
					expectedDCGContainerCallName := "name"

					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					fakeStrInterpreter.InterpretReturns(&model.Value{String: &expectedDCGContainerCallName}, nil)

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
						&model.SCGContainerCall{
							Name: new(string),
						},
						"dummyContainerID",
						"dummyRootOpID",
						new(modelFakes.FakeDataHandle),
					)

					/* assert */
					Expect(*actualResult.Name).To(Equal(expectedDCGContainerCallName))
				})
			})
		})

		Context("scgContainerCall.WorkDir truthy", func() {
			It("should call stringInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				containerFileBind := "dummyContainerFileBind"

				providedScope := map[string]*model.Value{
					containerFileBind: {String: new(string)},
				}
				providedSCGContainerCall := &model.SCGContainerCall{
					Image:   &model.SCGContainerCallImage{},
					WorkDir: "workDir",
				}

				providedOpHandle := new(modelFakes.FakeDataHandle)

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
					providedSCGContainerCall,
					"dummyContainerID",
					"dummyRootOpID",
					providedOpHandle,
				)

				/* assert */
				actualScope,
					actualScgContainerCallWorkDir,
					actualOpHandle := fakeStrInterpreter.InterpretArgsForCall(0)

				Expect(actualScope).To(Equal(providedScope))
				Expect(actualScgContainerCallWorkDir).To(Equal(providedSCGContainerCall.WorkDir))
				Expect(actualOpHandle).To(Equal(providedOpHandle))

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
						&model.SCGContainerCall{
							WorkDir: "dummyWorkDir",
						},
						"dummyContainerID",
						"dummyRootOpID",
						new(modelFakes.FakeDataHandle),
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("stringInterpreter.Interpret doesn't error", func() {
				It("should return expected dcgContainerCall.WorkDir", func() {
					/* arrange */
					expectedDCGContainerCallWorkDir := "workDir"

					fakeStrInterpreter := new(strFakes.FakeInterpreter)
					fakeStrInterpreter.InterpretReturns(&model.Value{String: &expectedDCGContainerCallWorkDir}, nil)

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
						&model.SCGContainerCall{
							WorkDir: "dummyWorkDir",
						},
						"dummyContainerID",
						"dummyRootOpID",
						new(modelFakes.FakeDataHandle),
					)

					/* assert */
					Expect(actualResult.WorkDir).To(Equal(expectedDCGContainerCallWorkDir))
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
			providedSCGContainerCall := &model.SCGContainerCall{
				Sockets: map[string]string{
					// implicitly bound
					envVarName: "",
				},
			}

			provideddataDirPath := "dummydataDirPath"
			providedContainerID := "dummyContainerID"
			providedRootOpID := "dummyRootOpID"

			expectedScratchDirPath := filepath.Join(
				provideddataDirPath,
				"dcg",
				providedRootOpID,
				"containers",
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
				providedSCGContainerCall,
				providedContainerID,
				providedRootOpID,
				new(modelFakes.FakeDataHandle),
			)

			/* assert */
			actualScope, actualScgContainerCallSockets, actualScratchDir := fakeSocketsInterpreter.InterpretArgsForCall(0)
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallSockets).To(Equal(providedSCGContainerCall.Sockets))
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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("socketsFakes.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.Sockets", func() {
				/* arrange */
				expectedDCGContainerCallSockets := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeSocketsInterpreter := new(socketsFakes.FakeInterpreter)
				fakeSocketsInterpreter.InterpretReturns(expectedDCGContainerCallSockets, nil)

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
					&model.SCGContainerCall{},
					"dummyContainerID",
					"dummyRootOpID",
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualResult.Sockets).To(Equal(expectedDCGContainerCallSockets))
			})
		})
	})
})
