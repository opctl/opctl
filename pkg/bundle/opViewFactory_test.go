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

				fakeFileSystem := new(fs.Fake)
				fakeFileSystem.GetBytesOfFileReturns(nil, expectedError)

				objectUnderTest := newOpViewFactory(
					fakeFileSystem,
					new(format.Fake),
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

				fakeYamlFormat := new(format.Fake)
				fakeYamlFormat.ToReturns(expectedError)

				objectUnderTest := newOpViewFactory(
					new(fs.Fake),
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

			fakeFileSystem := new(fs.Fake)
			fakeFileSystem.GetBytesOfFileReturns(expectedBytes, nil)

			fakeYamlFormat := new(format.Fake)

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
			dummyParams := map[string]*model.Param{
				"dummyName": {
					String: &model.StringParam{
						Constraints: &model.StringConstraints{
							MinLength: 0,
							MaxLength: 1000,
							Pattern:   "dummyPattern",
							Format:    "dummyFormat",
							Enum:      []string{"dummyEnumItem1"},
						},
						Default:     "dummyDefault",
						Description: "dummyDescription",
						IsSecret:    true,
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

			fakeFileSystem := new(fs.Fake)

			fakeYamlFormat := new(format.Fake)
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

				fakeFileSystem := new(fs.Fake)

				fakeYamlFormat := new(format.Fake)
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

				fakeFileSystem := new(fs.Fake)

				fakeYamlFormat := new(format.Fake)
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

					fakeFileSystem := new(fs.Fake)

					fakeYamlFormat := new(format.Fake)
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
					Inputs: map[string]*model.Param{
						"dockerPassword": {
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									MinLength: 1,
								},
								IsSecret: true,
							},
						},
						"dockerUsername": {
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Format: "email",
								},
								IsSecret: true,
							},
						},
						"dockerSocket": {
							Socket: &model.SocketParam{
								Description: "socket for docker daemon",
							},
						},
					},
					Outputs: map[string]*model.Param{
						"dockerConfig": {
							File: &model.FileParam{
								Description: "config for docker CLI",
								IsSecret:    true,
							},
						},
					},
					Run: &model.Scg{
						Container: &model.ScgContainerCall{
							Cmd: []string{
								"docker", "login", "-u", "$(dockerUsername)", "-p", "$(dockerPassword)",
							},
							Files: map[string]*model.ScgBinding{
								"/root/.docker/config.json": {
									Bind: "dockerConfig",
								},
							},
							Image: "docker:1.13",
							Sockets: map[string]*model.ScgBinding{
								"/var/run/docker.sock": {
									Bind: "dockerSocket",
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
					Inputs: map[string]*model.Param{
						"npmRegistryUrl": {
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Format: "uri",
								},
								Default: "https://npmjs.org/",
							},
						},
						"appDir": {
							Dir: &model.DirParam{
								Default:     ".",
								Description: "directory containing the app",
							},
						},
					},
					Outputs: map[string]*model.Param{
						"appDir": {
							Dir: &model.DirParam{
								Description: "directory containing the app",
							},
						},
					},
					Run: &model.Scg{
						Serial: []*model.Scg{
							{
								Op: &model.ScgOpCall{
									Ref: "install-deps",
									Inputs: map[string]string{
										"npmRegistryUrl": "",
										"appDir":         "",
									},
									Outputs: map[string]string{
										"appDir": "",
									},
								},
							},
							{
								Op: &model.ScgOpCall{
									Ref: "debug-api",
									Inputs: map[string]string{
										"appDir": "",
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
					Inputs: map[string]*model.Param{
						"appDir": {
							Dir: &model.DirParam{
								Default:     ".",
								Description: "directory containing the app",
							},
						},
						"apiSocket": {
							Socket: &model.SocketParam{
								Description: "socket for the API",
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
							Dirs: map[string]*model.ScgBinding{
								"/opt/app": {
									Bind: "appDir",
								},
							},
							Image: "node:7.4",
							Sockets: map[string]*model.ScgBinding{
								"api": {
									Bind: "apiSocket",
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
