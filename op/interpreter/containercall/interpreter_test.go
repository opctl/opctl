package containercall

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/dirs"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/envvars"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/files"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/image"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall/sockets"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
	"os"
	"path/filepath"
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
			rootFSPath := "/dummyRootFSPath"
			providedContainerId := "dummyContainerId"
			providedRootOpId := "dummyRootOpId"

			expectedScratchDirPath := filepath.Join(
				rootFSPath,
				"dcg",
				providedRootOpId,
				"containers",
				providedContainerId,
				"fs",
			)
			expectedScratchDirMode := os.FileMode(0700)

			expectedError := errors.New("dummyError")

			fakeOS := new(ios.Fake)
			// error to trigger immediate return
			fakeOS.MkdirAllReturns(expectedError)

			objectUnderTest := _interpreter{
				os:         fakeOS,
				rootFSPath: rootFSPath,
			}

			/* act */
			_, actualError := objectUnderTest.Interpret(
				map[string]*model.Value{},
				&model.SCGContainerCall{},
				providedContainerId,
				providedRootOpId,
				new(data.FakeHandle),
			)

			/* assert */
			actualScratchDirPath, actualScratchDirMode := fakeOS.MkdirAllArgsForCall(0)
			Expect(actualScratchDirPath).To(Equal(expectedScratchDirPath))
			Expect(actualScratchDirMode).To(Equal(expectedScratchDirMode))
			Expect(actualError).To(Equal(expectedError))
		})

		Context("container.Cmd not empty", func() {
			It("should call expression.EvalToString w/ expected args for each container.Cmd entry", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Value{
					"name1": {String: &providedString1},
				}
				providedOpDirHandle := new(data.FakeHandle)

				providedSCGContainerCall := &model.SCGContainerCall{
					Cmd: []interface{}{
						"dummy1",
						"dummy2",
					},
				}

				fakeExpression := new(expression.Fake)
				fakeExpression.EvalToStringReturns(&model.Value{String: new(string)}, nil)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   new(files.FakeInterpreter),
					imageInterpreter:   new(image.FakeInterpreter),
					expression:         fakeExpression,
					os:                 new(ios.Fake),
					socketsInterpreter: new(sockets.FakeInterpreter),
				}

				/* act */
				objectUnderTest.Interpret(
					providedCurrentScope,
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					providedOpDirHandle,
				)

				/* assert */
				for expectedCmdIndex, expectedCmdEntry := range providedSCGContainerCall.Cmd {
					actualScope,
						actualCmdEntry,
						actualOpDirHandle := fakeExpression.EvalToStringArgsForCall(expectedCmdIndex)
					Expect(actualCmdEntry).To(Equal(expectedCmdEntry))
					Expect(actualScope).To(Equal(providedCurrentScope))
					Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
				}
			})
			It("should return expected dcg.Cmd", func() {
				/* arrange */
				expectedCmd := []string{
					"dummyCmdEntry1",
					"dummyCmdEntry2",
				}

				providedSCGContainerCall := &model.SCGContainerCall{
					Cmd: []interface{}{
						"dummy1",
						"dummy2",
					},
				}

				fakeExpression := new(expression.Fake)
				fakeExpression.EvalToStringReturnsOnCall(0, &model.Value{String: &expectedCmd[0]}, nil)
				fakeExpression.EvalToStringReturnsOnCall(1, &model.Value{String: &expectedCmd[1]}, nil)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   new(files.FakeInterpreter),
					imageInterpreter:   new(image.FakeInterpreter),
					expression:         fakeExpression,
					os:                 new(ios.Fake),
					socketsInterpreter: new(sockets.FakeInterpreter),
				}

				/* act */
				actualDCGContainerCall, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualDCGContainerCall.Cmd).To(Equal(expectedCmd))
			})
		})

		It("should call dirs w/ expected args", func() {
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

			providedRootFSPath := "dummyRootFSPath"
			providedContainerId := "dummyContainerId"
			providedRootOpId := "dummyRootOpId"
			providedOpDirHandle := new(data.FakeHandle)

			expectedScratchDirPath := filepath.Join(
				providedRootFSPath,
				"dcg",
				providedRootOpId,
				"containers",
				providedContainerId,
				"fs",
			)

			fakeDirsInterpreter := new(dirs.FakeInterpreter)

			objectUnderTest := _interpreter{
				dirsInterpreter:    fakeDirsInterpreter,
				envVarsInterpreter: new(envvars.FakeInterpreter),
				filesInterpreter:   new(files.FakeInterpreter),
				imageInterpreter:   new(image.FakeInterpreter),
				os:                 new(ios.Fake),
				rootFSPath:         providedRootFSPath,
				socketsInterpreter: new(sockets.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				providedContainerId,
				providedRootOpId,
				providedOpDirHandle,
			)

			/* assert */
			actualOpDirHandle, actualScope, actualScgContainerCallDirs, actualScratchDir := fakeDirsInterpreter.InterpretArgsForCall(0)
			Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallDirs).To(Equal(providedSCGContainerCall.Dirs))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("dirs.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeDirsInterpreter := new(dirs.FakeInterpreter)
				fakeDirsInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					dirsInterpreter:    fakeDirsInterpreter,
					envVarsInterpreter: new(envvars.FakeInterpreter),
					imageInterpreter:   new(image.FakeInterpreter),
					os:                 new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("dirs.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.Dirs", func() {
				/* arrange */
				expectedDCGContainerCallDirs := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeDirsInterpreter := new(dirs.FakeInterpreter)
				fakeDirsInterpreter.InterpretReturns(expectedDCGContainerCallDirs, nil)

				objectUnderTest := _interpreter{
					dirsInterpreter:    fakeDirsInterpreter,
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   new(files.FakeInterpreter),
					imageInterpreter:   new(image.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(sockets.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualResult.Dirs).To(Equal(expectedDCGContainerCallDirs))
			})
		})

		It("should call envVars w/ expected args", func() {
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

			providedOpDirHandle := new(data.FakeHandle)

			fakeEnvVarsInterpreter := new(envvars.FakeInterpreter)

			objectUnderTest := _interpreter{
				dirsInterpreter:    new(dirs.FakeInterpreter),
				envVarsInterpreter: fakeEnvVarsInterpreter,
				filesInterpreter:   new(files.FakeInterpreter),
				imageInterpreter:   new(image.FakeInterpreter),
				os:                 new(ios.Fake),
				socketsInterpreter: new(sockets.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				"dummyContainerId",
				"dummyRootOpId",
				providedOpDirHandle,
			)

			/* assert */
			actualScope,
				actualScgContainerCallEnvVars,
				actualOpDirHandle := fakeEnvVarsInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallEnvVars).To(Equal(providedSCGContainerCall.EnvVars))
			Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))

		})
		Context("envVars.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeEnvVarsInterpreter := new(envvars.FakeInterpreter)
				fakeEnvVarsInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: fakeEnvVarsInterpreter,
					os:                 new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
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

				fakeEnvVarsInterpreter := new(envvars.FakeInterpreter)
				fakeEnvVarsInterpreter.InterpretReturns(expectedDCGContainerCallEnvVars, nil)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: fakeEnvVarsInterpreter,
					filesInterpreter:   new(files.FakeInterpreter),
					imageInterpreter:   new(image.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(sockets.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualResult.EnvVars).To(Equal(expectedDCGContainerCallEnvVars))
			})
		})

		It("should call files w/ expected args", func() {
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

			providedRootFSPath := "dummyRootFSPath"
			providedContainerId := "dummyContainerId"
			providedRootOpId := "dummyRootOpId"
			providedOpDirHandle := new(data.FakeHandle)

			expectedScratchDirPath := filepath.Join(
				providedRootFSPath,
				"dcg",
				providedRootOpId,
				"containers",
				providedContainerId,
				"fs",
			)

			fakeFilesInterpreter := new(files.FakeInterpreter)

			objectUnderTest := _interpreter{
				dirsInterpreter:    new(dirs.FakeInterpreter),
				envVarsInterpreter: new(envvars.FakeInterpreter),
				filesInterpreter:   fakeFilesInterpreter,
				imageInterpreter:   new(image.FakeInterpreter),
				os:                 new(ios.Fake),
				rootFSPath:         providedRootFSPath,
				socketsInterpreter: new(sockets.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				providedContainerId,
				providedRootOpId,
				providedOpDirHandle,
			)

			/* assert */
			actualOpDirHandle, actualScope, actualScgContainerCallFiles, actualScratchDir := fakeFilesInterpreter.InterpretArgsForCall(0)
			Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallFiles).To(Equal(providedSCGContainerCall.Files))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("files.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeFilesInterpreter := new(files.FakeInterpreter)
				fakeFilesInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   fakeFilesInterpreter,
					imageInterpreter:   new(image.FakeInterpreter),
					os:                 new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("files.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.Files", func() {
				/* arrange */
				expectedDCGContainerCallFiles := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeFilesInterpreter := new(files.FakeInterpreter)
				fakeFilesInterpreter.InterpretReturns(expectedDCGContainerCallFiles, nil)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   fakeFilesInterpreter,
					imageInterpreter:   new(image.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: new(sockets.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualResult.Files).To(Equal(expectedDCGContainerCallFiles))
			})
		})

		It("should call image w/ expected args", func() {
			/* arrange */
			containerFileBind := "dummyContainerFileBind"

			providedScope := map[string]*model.Value{
				containerFileBind: {String: new(string)},
			}
			providedSCGContainerCall := &model.SCGContainerCall{
				Image: &model.SCGContainerCallImage{
					Ref: "dummyImageRef",
					PullCreds: &model.SCGPullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				},
			}

			providedOpDirHandle := new(data.FakeHandle)

			fakeImageInterpreter := new(image.FakeInterpreter)

			objectUnderTest := _interpreter{
				dirsInterpreter:    new(dirs.FakeInterpreter),
				envVarsInterpreter: new(envvars.FakeInterpreter),
				filesInterpreter:   new(files.FakeInterpreter),
				imageInterpreter:   fakeImageInterpreter,
				os:                 new(ios.Fake),
				socketsInterpreter: new(sockets.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				"dummyContainerId",
				"dummyRootOpId",
				providedOpDirHandle,
			)

			/* assert */
			actualScope,
				actualScgContainerCallImage,
				actualOpDirHandle := fakeImageInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallImage).To(Equal(providedSCGContainerCall.Image))
			Expect(actualOpDirHandle).To(Equal(providedOpDirHandle))

		})
		Context("image.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeImageInterpreter := new(image.FakeInterpreter)
				fakeImageInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   new(files.FakeInterpreter),
					imageInterpreter:   fakeImageInterpreter,
					os:                 new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("image.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.Image", func() {
				/* arrange */
				expectedDCGContainerCallImage := &model.DCGContainerCallImage{
					Ref: "dummyImageRef",
					PullCreds: &model.DCGPullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				}

				fakeImageInterpreter := new(image.FakeInterpreter)
				fakeImageInterpreter.InterpretReturns(expectedDCGContainerCallImage, nil)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   new(files.FakeInterpreter),
					imageInterpreter:   fakeImageInterpreter,
					os:                 new(ios.Fake),
					socketsInterpreter: new(sockets.FakeInterpreter),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualResult.Image).To(Equal(expectedDCGContainerCallImage))
			})
		})

		It("should call sockets w/ expected args", func() {
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

			providedRootFSPath := "dummyRootFSPath"
			providedContainerId := "dummyContainerId"
			providedRootOpId := "dummyRootOpId"

			expectedScratchDirPath := filepath.Join(
				providedRootFSPath,
				"dcg",
				providedRootOpId,
				"containers",
				providedContainerId,
				"fs",
			)

			fakeSocketsInterpreter := new(sockets.FakeInterpreter)

			objectUnderTest := _interpreter{
				dirsInterpreter:    new(dirs.FakeInterpreter),
				envVarsInterpreter: new(envvars.FakeInterpreter),
				filesInterpreter:   new(files.FakeInterpreter),
				imageInterpreter:   new(image.FakeInterpreter),
				os:                 new(ios.Fake),
				rootFSPath:         providedRootFSPath,
				socketsInterpreter: fakeSocketsInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				providedContainerId,
				providedRootOpId,
				new(data.FakeHandle),
			)

			/* assert */
			actualScope, actualScgContainerCallSockets, actualScratchDir := fakeSocketsInterpreter.InterpretArgsForCall(0)
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallSockets).To(Equal(providedSCGContainerCall.Sockets))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("sockets.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeSocketsInterpreter := new(sockets.FakeInterpreter)
				fakeSocketsInterpreter.InterpretReturns(nil, expectedErr)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   new(files.FakeInterpreter),
					imageInterpreter:   new(image.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: fakeSocketsInterpreter,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("sockets.Interpret doesn't error", func() {
			It("should return expected dcgContainerCall.Sockets", func() {
				/* arrange */
				expectedDCGContainerCallSockets := map[string]string{
					"dummyName": "dummyValue",
				}

				fakeSocketsInterpreter := new(sockets.FakeInterpreter)
				fakeSocketsInterpreter.InterpretReturns(expectedDCGContainerCallSockets, nil)

				objectUnderTest := _interpreter{
					dirsInterpreter:    new(dirs.FakeInterpreter),
					envVarsInterpreter: new(envvars.FakeInterpreter),
					filesInterpreter:   new(files.FakeInterpreter),
					imageInterpreter:   new(image.FakeInterpreter),
					os:                 new(ios.Fake),
					socketsInterpreter: fakeSocketsInterpreter,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualResult.Sockets).To(Equal(expectedDCGContainerCallSockets))
			})
		})
	})
})
