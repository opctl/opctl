package core

import (
	"github.com/golang-utils/dircopier"
	"github.com/golang-utils/filecopier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/interpolate"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strconv"
)

var _ = Context("dcgFactory", func() {
	Context("containerCall", func() {
		Context("scg.Cmd not empty", func() {
			It("should call interpolate w/ expected args for each scg.Cmd entry", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Data{
					"name1": {String: &providedString1},
				}

				providedSCGContainerCall := &model.SCGContainerCall{
					Cmd: []string{
						"dummy1",
						"dummy2",
					},
				}

				fakeInterpolate := new(interpolate.Fake)

				objectUnderTest := _dcgFactory{
					interpolate: fakeInterpolate,
				}

				/* act */
				objectUnderTest.Construct(
					providedCurrentScope,
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					"dummyPkgRef",
				)

				/* assert */
				for expectedCmdIndex, expectedCmdEntry := range providedSCGContainerCall.Cmd {
					actualCmdEntry, actualScope := fakeInterpolate.InterpolateArgsForCall(expectedCmdIndex)
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

				fakeInterpolate := new(interpolate.Fake)
				fakeInterpolate.InterpolateReturnsOnCall(0, expectedCmd[0])
				fakeInterpolate.InterpolateReturnsOnCall(1, expectedCmd[1])

				objectUnderTest := _dcgFactory{
					interpolate: fakeInterpolate,
				}

				/* act */
				actualDCGContainerCall, _ := objectUnderTest.Construct(
					map[string]*model.Data{},
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					"dummyPkgRef",
				)

				/* assert */
				Expect(actualDCGContainerCall.Cmd).To(Equal(expectedCmd))
			})
		})
		Context("scg.Dirs not empty", func() {
			It("should return expected dcg.Dirs", func() {

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

				expectedDir1Path := "/dummyFile1Path.txt"
				expectedDirs := map[string]string{
					expectedDir1Path: filepath.Join(expectedScratchDirPath, expectedDir1Path),
				}

				providedSCGContainerCall := &model.SCGContainerCall{
					Dirs: map[string]string{
						// implicitly bound
						expectedDir1Path: "",
					},
				}

				objectUnderTest := _dcgFactory{
					dirCopier:  new(dircopier.Fake),
					rootFSPath: rootFSPath,
				}

				/* act */
				actualDCGContainerCall, _ := objectUnderTest.Construct(
					map[string]*model.Data{},
					providedSCGContainerCall,
					providedContainerId,
					providedRootOpId,
					"dummyPkgRef",
				)

				/* assert */
				Expect(actualDCGContainerCall.Dirs).To(Equal(expectedDirs))
			})
		})
		Context("scg.EnvVars not empty", func() {
			It("should return expected dcg.EnvVars", func() {

				/* arrange */
				providedCurrentScopeRef1 := "dummyScopeRef1"
				providedCurrentScopeRef1String := "dummyScopeRef1String"
				providedCurrentScopeRef2 := "dummyScopeRef2"
				providedCurrentScopeRef2Number := float64(2.3)
				providedCurrentScope := map[string]*model.Data{
					providedCurrentScopeRef1: {String: &providedCurrentScopeRef1String},
					providedCurrentScopeRef2: {Number: &providedCurrentScopeRef2Number},
				}

				expectedEnvVars := map[string]string{
					providedCurrentScopeRef1: providedCurrentScopeRef1String,
					providedCurrentScopeRef2: strconv.FormatFloat(providedCurrentScopeRef2Number, 'f', -1, 64),
				}

				providedSCGContainerCall := &model.SCGContainerCall{
					EnvVars: map[string]string{
						// implicitly bound to string
						providedCurrentScopeRef1: "",
						// implicitly bound to number
						providedCurrentScopeRef2: "",
					},
				}

				objectUnderTest := _dcgFactory{}

				/* act */
				actualDCGContainerCall, _ := objectUnderTest.Construct(
					providedCurrentScope,
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					"dummyPkgRef",
				)

				/* assert */
				Expect(actualDCGContainerCall.EnvVars).To(Equal(expectedEnvVars))
			})
		})
		Context("scg.Files not empty", func() {
			It("should return expected dcg.Files", func() {

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

				expectedFile1Path := "/dummyFile1Path.txt"
				expectedFiles := map[string]string{
					expectedFile1Path: filepath.Join(expectedScratchDirPath, expectedFile1Path),
				}

				providedSCGContainerCall := &model.SCGContainerCall{
					Files: map[string]string{
						// implicitly bound
						expectedFile1Path: "",
					},
				}

				objectUnderTest := _dcgFactory{
					fileCopier: new(filecopier.Fake),
					rootFSPath: rootFSPath,
				}

				/* act */
				actualDCGContainerCall, _ := objectUnderTest.Construct(
					map[string]*model.Data{},
					providedSCGContainerCall,
					providedContainerId,
					providedRootOpId,
					"dummyPkgRef",
				)

				/* assert */
				Expect(actualDCGContainerCall.Files).To(Equal(expectedFiles))
			})
		})
		Context("scg.Image not empty", func() {
			It("should call interpolate w/ expected args for scg.Image.PullCreds.Username/Password", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Data{
					"name1": {String: &providedString1},
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

				fakeInterpolate := new(interpolate.Fake)

				objectUnderTest := _dcgFactory{
					interpolate: fakeInterpolate,
				}

				/* act */
				objectUnderTest.Construct(
					providedCurrentScope,
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					"dummyPkgRef",
				)

				/* assert */
				actualUsername, actualUsernameScope := fakeInterpolate.InterpolateArgsForCall(0)
				Expect(actualUsername).To(Equal(providedSCGContainerCall.Image.PullCreds.Username))
				Expect(actualUsernameScope).To(Equal(providedCurrentScope))

				actualPassword, actualPasswordScope := fakeInterpolate.InterpolateArgsForCall(1)
				Expect(actualPassword).To(Equal(providedSCGContainerCall.Image.PullCreds.Password))
				Expect(actualPasswordScope).To(Equal(providedCurrentScope))
			})
			It("should return expected dcg.Image", func() {

				/* arrange */
				providedSCGContainerCall := &model.SCGContainerCall{
					Image: &model.SCGContainerCallImage{
						Ref:       "dummyImageRef",
						PullCreds: &model.SCGPullCreds{},
					},
				}

				fakeInterpolate := new(interpolate.Fake)
				expectedUsername := "expectedUsername"
				fakeInterpolate.InterpolateReturnsOnCall(0, expectedUsername)

				expectedPassword := "expectedPassword"
				fakeInterpolate.InterpolateReturnsOnCall(1, expectedPassword)

				expectedImage := &model.DCGContainerCallImage{
					Ref: providedSCGContainerCall.Image.Ref,
					PullCreds: &model.DCGPullCreds{
						Username: expectedUsername,
						Password: expectedPassword,
					},
				}

				objectUnderTest := _dcgFactory{
					interpolate: fakeInterpolate,
				}

				/* act */
				actualDCGContainerCall, _ := objectUnderTest.Construct(
					map[string]*model.Data{},
					providedSCGContainerCall,
					"dummyContainerId",
					"dummyRootOpId",
					"dummyPkgRef",
				)

				/* assert */
				Expect(actualDCGContainerCall.Image).To(Equal(expectedImage))
			})
		})
	})
})
