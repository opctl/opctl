package bundle

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/opspec-io/sdk-golang/util/fs"
	"os"
	"reflect"
)

var _ = Describe("_opViewFactory", func() {

	Context("Construct", func() {

		Context("when FileSystem.GetBytesOfFile returns an error", func() {

			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("GetBytesOfFileError")

				fakeFileSystem := new(fs.FakeFileSystem)
				fakeFileSystem.GetBytesOfFileReturns(nil, expectedError)

				objectUnderTest := newOpViewFactory(
					fakeFileSystem,
					new(format.FakeFormat),
				)

				/* act */
				_, actualError := objectUnderTest.Construct("/dummy/path")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})

		})

		Context("when YamlFormat.From returns an error", func() {
			It("should be returned", func() {

				/* arrange */
				expectedError := errors.New("FromError")

				fakeYamlFormat := new(format.FakeFormat)
				fakeYamlFormat.ToReturns(expectedError)

				objectUnderTest := newOpViewFactory(
					new(fs.FakeFileSystem),
					fakeYamlFormat,
				)

				/* act */
				_, actualError := objectUnderTest.Construct("/dummy/path")

				/* assert */
				Expect(actualError).To(Equal(expectedError))

			})
		})

		It("should call YamlFormat.From with expected bytes", func() {

			/* arrange */
			expectedBytes := []byte{0, 8, 10}

			fakeFileSystem := new(fs.FakeFileSystem)
			fakeFileSystem.GetBytesOfFileReturns(expectedBytes, nil)

			fakeYamlFormat := new(format.FakeFormat)

			objectUnderTest := newOpViewFactory(
				fakeFileSystem,
				fakeYamlFormat,
			)

			/* act */
			objectUnderTest.Construct("/dummy/path")

			/* assert */
			actualBytes, _ := fakeYamlFormat.ToArgsForCall(0)
			Expect(actualBytes).To(Equal(expectedBytes))

		})

		It("should return expected opView", func() {

			/* arrange */
			dummyParams := []*model.Param{
				{
					String: &model.StringParam{
						Default:     "dummyDefault",
						Description: "dummyDescription",
						Constraints: &model.StringConstraints{
							Length: &model.StringLengthConstraint{
								Min:         0,
								Max:         1000,
								Description: "dummyStringLengthConstraintDescription",
							},
							Patterns: []*model.StringPatternConstraint{
								{
									Regex:       ".*",
									Description: "dummyStringPatternConstraintDescription",
								},
							},
						},
						Name:     "dummyName",
						IsSecret: true,
					},
				},
			}

			expectedCallGraph := &model.Scg{
				Op: &model.ScgOpCall{
					Ref: "dummyOpRef",
				},
			}

			expectedOpView := model.OpView{
				Description: "dummyDescription",
				Inputs:      dummyParams,
				Name:        "dummyName",
				Outputs:     dummyParams,
				Run:         expectedCallGraph,
				Version:     "dummyVersion",
			}

			fakeFileSystem := new(fs.FakeFileSystem)

			fakeYamlFormat := new(format.FakeFormat)
			fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

				stubbedOpManifest := model.OpManifest{
					Manifest: model.Manifest{
						Name:        expectedOpView.Name,
						Description: expectedOpView.Description,
						Version:     expectedOpView.Version,
					},
					Inputs:  dummyParams,
					Outputs: dummyParams,
					Run:     expectedCallGraph,
				}

				reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
				return
			}

			objectUnderTest := newOpViewFactory(
				fakeFileSystem,
				fakeYamlFormat,
			)

			/* act */
			actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

			/* assert */
			Expect(actualOpView).To(Equal(expectedOpView))

		})

		Context("when opManifest.Run.Parallel is not empty", func() {
			It("should return expected opView.Run", func() {

				/* arrange */

				expectedCallGraph := &model.Scg{
					Parallel: []*model.Scg{
						{
							Op: &model.ScgOpCall{
								Ref: "dummyRef",
							},
						},
					},
				}

				fakeFileSystem := new(fs.FakeFileSystem)

				fakeYamlFormat := new(format.FakeFormat)
				fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

					stubbedOpManifest := model.OpManifest{
						Run: expectedCallGraph,
					}

					reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
					return
				}

				objectUnderTest := newOpViewFactory(
					fakeFileSystem,
					fakeYamlFormat,
				)

				/* act */
				actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

				/* assert */
				Expect(actualOpView.Run).To(Equal(expectedCallGraph))

			})
		})
		Context("when opManifest.Run.Parallel is empty", func() {
			It("should return expected opView.Run", func() {

				/* arrange */
				expectedCallGraph := &model.Scg{
					Op: &model.ScgOpCall{
						Ref: "dummyOpRef",
					},
				}

				fakeFileSystem := new(fs.FakeFileSystem)

				fakeYamlFormat := new(format.FakeFormat)
				fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

					stubbedOpManifest := model.OpManifest{
						Run: expectedCallGraph,
					}

					reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
					return
				}

				objectUnderTest := newOpViewFactory(
					fakeFileSystem,
					fakeYamlFormat,
				)

				/* act */
				actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

				/* assert */
				Expect(actualOpView.Run).To(Equal(expectedCallGraph))

			})
			Context("when opManifest.Run.Serial is not empty", func() {
				It("should return expected opView.Run", func() {

					/* arrange */
					expectedCallGraph := &model.Scg{
						Serial: []*model.Scg{
							{
								Op: &model.ScgOpCall{
									Ref: "dummyRef",
								},
							},
						},
					}

					fakeFileSystem := new(fs.FakeFileSystem)

					fakeYamlFormat := new(format.FakeFormat)
					fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

						stubbedOpManifest := model.OpManifest{
							Run: expectedCallGraph,
						}

						reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
						return
					}

					objectUnderTest := newOpViewFactory(
						fakeFileSystem,
						fakeYamlFormat,
					)

					/* act */
					actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

					/* assert */
					Expect(actualOpView.Run).To(Equal(expectedCallGraph))

				})
			})

		})
		Context("when passed ./testdata/opspec-0.1.3/examples/docker/.opspec/login", func() {
			wd, err := os.Getwd()
			if nil != err {
				panic(err)
			}
			It("should return expected opView", func() {

				/* arrange */
				expectedOpView := model.OpView{
					Description: "Logs in to a docker registry",
					Name:        "login",
					Inputs: []*model.Param{
						{
							String: &model.StringParam{
								Name:        "DOCKER_PASSWORD",
								Description: "Password for docker registry",
								IsSecret:    true,
							},
						},
						{
							String: &model.StringParam{
								Name:        "DOCKER_USERNAME",
								Description: "Username for docker registry",
								IsSecret:    true,
							},
						},
					},
					Outputs: []*model.Param{
						{
							File: &model.FileParam{
								Name:        "DOCKER_CONFIG",
								Description: "Can be used for caching",
								IsSecret:    true,
							},
						},
					},
					Run: &model.Scg{
						Container: &model.ScgContainerCall{
							Cmd: []string{
								"sh",
								"-c",
								`# start dockerd
dockerd &>/dev/null &

echo sleeping for 2 sec to allow dockerd enough time to startup
sleep 2

docker login -u "$(DOCKER_USERNAME)" -p "$(DOCKER_PASSWORD)"
`,
							},
							Files: map[string]*model.ScgContainerFile{
								"/root/.docker/config.json": {
									Bind: "DOCKER_CONFIG",
								},
							},
							Image: "docker:1.12-dind",
						},
					},
				}

				objectUnderTest := newOpViewFactory(
					fs.NewFileSystem(),
					format.NewYamlFormat(),
				)

				/* act */
				actualOpView, err :=
					objectUnderTest.Construct(
						fmt.Sprintf("%v/../../testdata/opspec-0.1.3/examples/docker/.opspec/login", wd))
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualOpView).To(Equal(expectedOpView))

			})
		})
		Context("when passed ./testdata/opspec-0.1.3/examples/nodejs/.opspec/debug", func() {
			wd, err := os.Getwd()
			if nil != err {
				panic(err)
			}
			It("should return expected opView", func() {

				/* arrange */
				expectedOpView := model.OpView{
					Description: "Ensures deps are installed and debugs the node app",
					Name:        "debug",
					Inputs: []*model.Param{
						{
							String: &model.StringParam{
								Name:     "NPM_CONFIG_REGISTRY",
								IsSecret: true,
							},
						},
						{
							Dir: &model.DirParam{
								Name:        "APP_DIR",
								Description: "Directory containing the app",
							},
						},
					},
					Outputs: []*model.Param{
						{
							Dir: &model.DirParam{
								Name:        "APP_DIR",
								Description: "Directory containing the app (returned to support caching)",
							},
						},
					},
					Run: &model.Scg{
						Serial: []*model.Scg{
							{
								Op: &model.ScgOpCall{
									Ref: "install-deps",
									Inputs: map[string]string{
										"NPM_CONFIG_REGISTRY": "",
										"APP_DIR":             "",
									},
									Outputs: map[string]string{
										"APP_DIR": "",
									},
								},
							},
							{
								Op: &model.ScgOpCall{
									Ref: "debug-api",
									Inputs: map[string]string{
										"APP_DIR": "",
									},
								},
							},
						},
					},
				}

				objectUnderTest := newOpViewFactory(
					fs.NewFileSystem(),
					format.NewYamlFormat(),
				)

				/* act */
				actualOpView, err :=
					objectUnderTest.Construct(
						fmt.Sprintf("%v/../../testdata/opspec-0.1.3/examples/nodejs/.opspec/debug", wd))
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualOpView).To(Equal(expectedOpView))

			})
		})
		Context("when passed ./testdata/opspec-0.1.3/examples/nodejs/.opspec/test-acceptance", func() {
			wd, err := os.Getwd()
			if nil != err {
				panic(err)
			}
			It("should return expected opView", func() {

				/* arrange */
				expectedOpView := model.OpView{
					Description: "Runs acceptance tests",
					Name:        "test-acceptance",
					Inputs: []*model.Param{
						{
							Dir: &model.DirParam{
								Name:        "APP_DIR",
								Description: "Directory containing the app",
							},
						},
						{
							NetSocket: &model.NetSocketParam{
								Name:        "API_SOCKET",
								Description: "Network socket for the API under test",
								Constraints: &model.NetSocketConstraints{
									PortNumber: &model.PortNumberNetSocketConstraint{
										Number: 80,
									},
								},
							},
						},
					},
					Run: &model.Scg{
						Container: &model.ScgContainerCall{
							Cmd: []string{
								"node_modules/.bin/mocha",
								"--recursive",
								"--reporter=spec",
								"tests/unit",
							},
							Dirs: map[string]*model.ScgContainerDir{
								"/opt/app": {
									Bind: "APP_DIR",
								},
							},
							Image: "node:7.3",
							Net: []*model.ScgContainerNetEntry{
								{
									Bind: "API_SOCKET",
								},
							},
							WorkDir: "/opt/app",
						},
					},
				}

				objectUnderTest := newOpViewFactory(
					fs.NewFileSystem(),
					format.NewYamlFormat(),
				)

				/* act */
				actualOpView, err :=
					objectUnderTest.Construct(
						fmt.Sprintf("%v/../../testdata/opspec-0.1.3/examples/nodejs/.opspec/test-acceptance", wd))
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualOpView).To(Equal(expectedOpView))

			})
		})
	})
})
