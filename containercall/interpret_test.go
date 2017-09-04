package containercall

import (
	"errors"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/containercall/dirs"
	"github.com/opspec-io/sdk-golang/containercall/envvars"
	"github.com/opspec-io/sdk-golang/containercall/files"
	"github.com/opspec-io/sdk-golang/containercall/image"
	"github.com/opspec-io/sdk-golang/containercall/sockets"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	stringPkg "github.com/opspec-io/sdk-golang/string"
	"os"
	"path/filepath"
)

var _ = Context("ContainerCall", func() {
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

			objectUnderTest := _ContainerCall{
				os:         fakeOS,
				rootFSPath: rootFSPath,
			}

			/* act */
			_, actualError := objectUnderTest.Interpret(
				map[string]*model.Value{},
				&model.SCGContainerCall{},
				providedContainerId,
				providedRootOpId,
				new(pkg.FakeHandle),
			)

			/* assert */
			actualScratchDirPath, actualScratchDirMode := fakeOS.MkdirAllArgsForCall(0)
			Expect(actualScratchDirPath).To(Equal(expectedScratchDirPath))
			Expect(actualScratchDirMode).To(Equal(expectedScratchDirMode))
			Expect(actualError).To(Equal(expectedError))
		})

		Context("container.Cmd not empty", func() {
			It("should call string.Interpret w/ expected args for each container.Cmd entry", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Value{
					"name1": {String: &providedString1},
				}

				providedSCGContainerCall := &model.SCGContainerCall{
					Cmd: []string{
						"dummy1",
						"dummy2",
					},
				}

				fakeString := new(stringPkg.Fake)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: new(envvars.Fake),
					files:   new(files.Fake),
					image:   new(image.Fake),
					string:  fakeString,
					os:      new(ios.Fake),
					sockets: new(sockets.Fake),
				}

				/* act */
				objectUnderTest.Interpret(
					providedCurrentScope,
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
				)

				/* assert */
				for expectedCmdIndex, expectedCmdEntry := range providedSCGContainerCall.Cmd {
					actualScope, actualCmdEntry := fakeString.InterpretArgsForCall(expectedCmdIndex)
					Expect(actualCmdEntry).To(Equal(expectedCmdEntry))
					Expect(actualScope).To(Equal(providedCurrentScope))
				}
			})
			It("should return expected dcg.Cmd", func() {
				/* arrange */
				expectedCmd := []string{
					"dummyCmdEntry1",
					"dummyCmdEntry2",
				}

				providedSCGContainerCall := &model.SCGContainerCall{
					Cmd: []string{
						"dummy1",
						"dummy2",
					},
				}

				fakeString := new(stringPkg.Fake)
				fakeString.InterpretReturnsOnCall(0, expectedCmd[0], nil)
				fakeString.InterpretReturnsOnCall(1, expectedCmd[1], nil)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: new(envvars.Fake),
					files:   new(files.Fake),
					image:   new(image.Fake),
					string:  fakeString,
					os:      new(ios.Fake),
					sockets: new(sockets.Fake),
				}

				/* act */
				actualDCGContainerCall, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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
			providedPkgHandle := new(pkg.FakeHandle)

			expectedScratchDirPath := filepath.Join(
				providedRootFSPath,
				"dcg",
				providedRootOpId,
				"containers",
				providedContainerId,
				"fs",
			)

			fakeDirs := new(dirs.Fake)

			objectUnderTest := _ContainerCall{
				dirs:       fakeDirs,
				envVars:    new(envvars.Fake),
				files:      new(files.Fake),
				image:      new(image.Fake),
				os:         new(ios.Fake),
				rootFSPath: providedRootFSPath,
				sockets:    new(sockets.Fake),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				providedContainerId,
				providedRootOpId,
				providedPkgHandle,
			)

			/* assert */
			actualPkgHandle, actualScope, actualScgContainerCallDirs, actualScratchDir := fakeDirs.InterpretArgsForCall(0)
			Expect(actualPkgHandle).To(Equal(providedPkgHandle))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallDirs).To(Equal(providedSCGContainerCall.Dirs))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("dirs.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeDirs := new(dirs.Fake)
				fakeDirs.InterpretReturns(nil, expectedErr)

				objectUnderTest := _ContainerCall{
					dirs:    fakeDirs,
					envVars: new(envvars.Fake),
					image:   new(image.Fake),
					os:      new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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

				fakeDirs := new(dirs.Fake)
				fakeDirs.InterpretReturns(expectedDCGContainerCallDirs, nil)

				objectUnderTest := _ContainerCall{
					dirs:    fakeDirs,
					envVars: new(envvars.Fake),
					files:   new(files.Fake),
					image:   new(image.Fake),
					os:      new(ios.Fake),
					sockets: new(sockets.Fake),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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
				EnvVars: map[string]string{
					// implicitly bound
					envVarName: "",
				},
			}

			fakeEnvVars := new(envvars.Fake)

			objectUnderTest := _ContainerCall{
				dirs:    new(dirs.Fake),
				envVars: fakeEnvVars,
				files:   new(files.Fake),
				image:   new(image.Fake),
				os:      new(ios.Fake),
				sockets: new(sockets.Fake),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				"dummyContainerId",
				"dummyRootOpId",
				new(pkg.FakeHandle),
			)

			/* assert */
			actualScope, actualScgContainerCallEnvVars := fakeEnvVars.InterpretArgsForCall(0)
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallEnvVars).To(Equal(providedSCGContainerCall.EnvVars))
		})
		Context("envVars.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeEnvVars := new(envvars.Fake)
				fakeEnvVars.InterpretReturns(nil, expectedErr)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: fakeEnvVars,
					os:      new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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

				fakeEnvVars := new(envvars.Fake)
				fakeEnvVars.InterpretReturns(expectedDCGContainerCallEnvVars, nil)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: fakeEnvVars,
					files:   new(files.Fake),
					image:   new(image.Fake),
					os:      new(ios.Fake),
					sockets: new(sockets.Fake),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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
				Files: map[string]string{
					// implicitly bound
					fileName: "",
				},
			}

			providedRootFSPath := "dummyRootFSPath"
			providedContainerId := "dummyContainerId"
			providedRootOpId := "dummyRootOpId"
			providedPkgHandle := new(pkg.FakeHandle)

			expectedScratchDirPath := filepath.Join(
				providedRootFSPath,
				"dcg",
				providedRootOpId,
				"containers",
				providedContainerId,
				"fs",
			)

			fakeFiles := new(files.Fake)

			objectUnderTest := _ContainerCall{
				dirs:       new(dirs.Fake),
				envVars:    new(envvars.Fake),
				files:      fakeFiles,
				image:      new(image.Fake),
				os:         new(ios.Fake),
				rootFSPath: providedRootFSPath,
				sockets:    new(sockets.Fake),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				providedContainerId,
				providedRootOpId,
				providedPkgHandle,
			)

			/* assert */
			actualPkgHandle, actualScope, actualScgContainerCallFiles, actualScratchDir := fakeFiles.InterpretArgsForCall(0)
			Expect(actualPkgHandle).To(Equal(providedPkgHandle))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallFiles).To(Equal(providedSCGContainerCall.Files))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("files.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeFiles := new(files.Fake)
				fakeFiles.InterpretReturns(nil, expectedErr)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: new(envvars.Fake),
					files:   fakeFiles,
					image:   new(image.Fake),
					os:      new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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

				fakeFiles := new(files.Fake)
				fakeFiles.InterpretReturns(expectedDCGContainerCallFiles, nil)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: new(envvars.Fake),
					files:   fakeFiles,
					image:   new(image.Fake),
					os:      new(ios.Fake),
					sockets: new(sockets.Fake),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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

			fakeImage := new(image.Fake)

			objectUnderTest := _ContainerCall{
				dirs:    new(dirs.Fake),
				envVars: new(envvars.Fake),
				files:   new(files.Fake),
				image:   fakeImage,
				os:      new(ios.Fake),
				sockets: new(sockets.Fake),
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				"dummyContainerId",
				"dummyRootOpId",
				new(pkg.FakeHandle),
			)

			/* assert */
			actualScope, actualScgContainerCallImage := fakeImage.InterpretArgsForCall(0)
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallImage).To(Equal(providedSCGContainerCall.Image))
		})
		Context("image.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeImage := new(image.Fake)
				fakeImage.InterpretReturns(nil, expectedErr)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: new(envvars.Fake),
					files:   new(files.Fake),
					image:   fakeImage,
					os:      new(ios.Fake),
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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

				fakeImage := new(image.Fake)
				fakeImage.InterpretReturns(expectedDCGContainerCallImage, nil)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: new(envvars.Fake),
					files:   new(files.Fake),
					image:   fakeImage,
					os:      new(ios.Fake),
					sockets: new(sockets.Fake),
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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

			fakeSockets := new(sockets.Fake)

			objectUnderTest := _ContainerCall{
				dirs:       new(dirs.Fake),
				envVars:    new(envvars.Fake),
				files:      new(files.Fake),
				image:      new(image.Fake),
				os:         new(ios.Fake),
				rootFSPath: providedRootFSPath,
				sockets:    fakeSockets,
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedSCGContainerCall,
				providedContainerId,
				providedRootOpId,
				new(pkg.FakeHandle),
			)

			/* assert */
			actualScope, actualScgContainerCallSockets, actualScratchDir := fakeSockets.InterpretArgsForCall(0)
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualScgContainerCallSockets).To(Equal(providedSCGContainerCall.Sockets))
			Expect(actualScratchDir).To(Equal(expectedScratchDirPath))
		})
		Context("sockets.Interpret errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedErr := errors.New("dummyError")
				fakeSockets := new(sockets.Fake)
				fakeSockets.InterpretReturns(nil, expectedErr)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: new(envvars.Fake),
					files:   new(files.Fake),
					image:   new(image.Fake),
					os:      new(ios.Fake),
					sockets: fakeSockets,
				}

				/* act */
				_, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
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

				fakeSockets := new(sockets.Fake)
				fakeSockets.InterpretReturns(expectedDCGContainerCallSockets, nil)

				objectUnderTest := _ContainerCall{
					dirs:    new(dirs.Fake),
					envVars: new(envvars.Fake),
					files:   new(files.Fake),
					image:   new(image.Fake),
					os:      new(ios.Fake),
					sockets: fakeSockets,
				}

				/* act */
				actualResult, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					&model.SCGContainerCall{},
					"dummyContainerId",
					"dummyRootOpId",
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(actualResult.Sockets).To(Equal(expectedDCGContainerCallSockets))
			})
		})
	})
})
