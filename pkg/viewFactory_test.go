package pkg

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/vioutil"
	"os"
	"path"
	"reflect"
)

var _ = Describe("_viewFactory", func() {
	wd, err := os.Getwd()
	if nil != err {
		panic(err)
	}

	Context("Construct", func() {

		It("should call validate w/ expected inputs", func() {
			/* arrange */
			providedPkgRef := "/dummy/path"

			fakeValidator := new(fakeValidator)

			// err to cause immediate return
			fakeValidator.ValidateReturns([]error{errors.New("dummyError")})

			objectUnderTest := _viewFactory{
				validator: fakeValidator,
			}

			/* act */
			objectUnderTest.Construct(providedPkgRef)

			/* assert */
			Expect(fakeValidator.ValidateArgsForCall(0)).To(Equal(providedPkgRef))
		})
		Context("validator.Validate returns errors", func() {
			It("should return the expected error", func() {
				/* arrange */

				errs := []error{errors.New("dummyErr1"), errors.New("dummyErr2")}
				expectedErr := fmt.Errorf(`
-
  Error(s):
    - %v
    - %v
-`, errs[0], errs[1])

				fakeValidator := new(fakeValidator)

				// err to cause immediate return
				fakeValidator.ValidateReturns(errs)

				objectUnderTest := _viewFactory{
					validator: fakeValidator,
				}

				/* act */
				_, actualError := objectUnderTest.Construct("")

				/* assert */
				Expect(actualError).To(Equal(expectedErr))
			})
		})
		Context("validator.Validate doesn't return errors", func() {
			It("should call ioutil.ReadFile w/expected args", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"

				fakeIOUtil := new(vioutil.Fake)
				// err to cause immediate return
				fakeIOUtil.ReadFileReturns(nil, errors.New("dummyError"))

				objectUnderTest := newViewFactory(
					fakeIOUtil,
					new(fakeValidator),
					new(format.Fake),
				)

				/* act */
				objectUnderTest.Construct(providedPkgRef)

				/* assert */
				Expect(fakeIOUtil.ReadFileArgsForCall(0)).
					To(Equal(path.Join(providedPkgRef, NameOfPkgManifestFile)))

			})
			Context("ioutil.ReadFile returns an error", func() {

				It("should return expected error", func() {

					/* arrange */
					expectedError := errors.New("dummyError")

					fakeIOUtil := new(vioutil.Fake)
					fakeIOUtil.ReadFileReturns(nil, expectedError)

					objectUnderTest := newViewFactory(
						fakeIOUtil,
						new(fakeValidator),
						new(format.Fake),
					)

					/* act */
					_, actualError := objectUnderTest.Construct("/dummy/path")

					/* assert */
					Expect(actualError).To(Equal(expectedError))

				})

			})

			Context("YamlFormat.From returns an error", func() {
				It("should return expected error", func() {

					/* arrange */
					expectedError := errors.New("FromError")

					fakeYamlFormat := new(format.Fake)
					fakeYamlFormat.ToReturns(expectedError)

					objectUnderTest := newViewFactory(
						new(vioutil.Fake),
						new(fakeValidator),
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

				fakeIOUtil := new(vioutil.Fake)
				fakeIOUtil.ReadFileReturns(expectedBytes, nil)

				fakeYamlFormat := new(format.Fake)

				objectUnderTest := newViewFactory(
					fakeIOUtil,
					new(fakeValidator),
					fakeYamlFormat,
				)

				/* act */
				objectUnderTest.Construct("/dummy/path")

				/* assert */
				actualBytes, _ := fakeYamlFormat.ToArgsForCall(0)
				Expect(actualBytes).To(Equal(expectedBytes))

			})

			It("should return expected packageView", func() {

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

				expectedCallGraph := &model.SCG{
					Op: &model.SCGOpCall{
						Pkg: &model.SCGOpCallPkg{
							Ref: "dummyPkgRef",
						},
					},
				}

				expectedPackageView := &model.PackageView{
					Description: "dummyDescription",
					Inputs:      dummyParams,
					Name:        "dummyName",
					Outputs:     dummyParams,
					Run:         expectedCallGraph,
					Version:     "dummyVersion",
				}

				fakeYamlFormat := new(format.Fake)
				fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

					stubbedPackageManifestView := model.PackageManifestView{
						Name:        expectedPackageView.Name,
						Description: expectedPackageView.Description,
						Version:     expectedPackageView.Version,
						Inputs:      dummyParams,
						Outputs:     dummyParams,
						Run:         expectedCallGraph,
					}

					reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedPackageManifestView))
					return
				}

				objectUnderTest := newViewFactory(
					new(vioutil.Fake),
					new(fakeValidator),
					fakeYamlFormat,
				)

				/* act */
				actualPackageView, _ := objectUnderTest.Construct("/dummy/op/path")

				/* assert */
				Expect(actualPackageView).To(Equal(expectedPackageView))

			})

			Context("packageManifestView.Run.Parallel is not empty", func() {
				It("should return expected packageView.Run", func() {

					/* arrange */

					expectedCallGraph := &model.SCG{
						Parallel: []*model.SCG{
							{
								Op: &model.SCGOpCall{
									Pkg: &model.SCGOpCallPkg{
										Ref: "dummyPkgRef",
									},
								},
							},
						},
					}

					fakeYamlFormat := new(format.Fake)
					fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

						stubbedPackageManifestView := model.PackageManifestView{
							Run: expectedCallGraph,
						}

						reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedPackageManifestView))
						return
					}

					objectUnderTest := newViewFactory(
						new(vioutil.Fake),
						new(fakeValidator),
						fakeYamlFormat,
					)

					/* act */
					actualPackageView, _ := objectUnderTest.Construct("/dummy/op/path")

					/* assert */
					Expect(actualPackageView.Run).To(Equal(expectedCallGraph))

				})
			})
			Context("packageManifestView.Run.Parallel is empty", func() {
				It("should return expected packageView.Run", func() {

					/* arrange */
					expectedCallGraph := &model.SCG{
						Op: &model.SCGOpCall{
							Pkg: &model.SCGOpCallPkg{
								Ref: "dummyPkgRef",
							},
						},
					}

					fakeYamlFormat := new(format.Fake)
					fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

						stubbedPackageManifestView := model.PackageManifestView{
							Run: expectedCallGraph,
						}

						reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedPackageManifestView))
						return
					}

					objectUnderTest := newViewFactory(
						new(vioutil.Fake),
						new(fakeValidator),
						fakeYamlFormat,
					)

					/* act */
					actualPackageView, _ := objectUnderTest.Construct("/dummy/op/path")

					/* assert */
					Expect(actualPackageView.Run).To(Equal(expectedCallGraph))

				})
				Context("packageManifestView.Run.Serial is not empty", func() {
					It("should return expected packageView.Run", func() {

						/* arrange */
						expectedCallGraph := &model.SCG{
							Serial: []*model.SCG{
								{
									Op: &model.SCGOpCall{
										Pkg: &model.SCGOpCallPkg{
											Ref: "dummyPkgRef",
										},
									},
								},
							},
						}

						fakeYamlFormat := new(format.Fake)
						fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

							stubbedPackageManifestView := model.PackageManifestView{
								Run: expectedCallGraph,
							}

							reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedPackageManifestView))
							return
						}

						objectUnderTest := newViewFactory(
							new(vioutil.Fake),
							new(fakeValidator),
							fakeYamlFormat,
						)

						/* act */
						actualPackageView, _ := objectUnderTest.Construct("/dummy/op/path")

						/* assert */
						Expect(actualPackageView.Run).To(Equal(expectedCallGraph))

					})
				})

			})
		})
		Context("passed ./testdata/opspec-0.1.3/examples/docker/.opspec/login", func() {
			It("should return expected packageView", func() {

				/* arrange */
				expectedPackageView := &model.PackageView{
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
					Run: &model.SCG{
						Container: &model.SCGContainerCall{
							Cmd: []string{
								"docker", "login", "-u", "$(dockerUsername)", "-p", "$(dockerPassword)",
							},
							Files: map[string]string{
								"/root/.docker/config.json": "dockerConfig",
							},
							Image: &model.SCGContainerCallImage{
								Ref: "docker:1.13",
							},
							Sockets: map[string]string{
								"/var/run/docker.sock": "dockerSocket",
							},
						},
					},
				}

				objectUnderTest := newViewFactory(
					vioutil.New(osfs.New()),
					new(fakeValidator),
					format.NewYamlFormat(),
				)

				/* act */
				actualPackageView, err :=
					objectUnderTest.Construct(
						fmt.Sprintf("%v/../testdata/opspec-0.1.3/examples/docker/.opspec/login", wd))
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualPackageView).To(Equal(expectedPackageView))

			})
		})
		Context("passed ./testdata/opspec-0.1.3/examples/nodejs/.opspec/debug", func() {
			It("should return expected packageView", func() {

				/* arrange */
				expectedPackageView := &model.PackageView{
					Description: "Ensures deps are installed and debugs the node app",
					Name:        "debug",
					Inputs: map[string]*model.Param{
						"npmRegistryUrl": {
							String: &model.StringParam{
								Constraints: &model.StringConstraints{
									Format: "uri",
								},
								Default: "https://registry.npmjs.org/",
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
					Run: &model.SCG{
						Serial: []*model.SCG{
							{
								Op: &model.SCGOpCall{
									Pkg: &model.SCGOpCallPkg{
										Ref: "install-deps",
									},
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
								Op: &model.SCGOpCall{
									Pkg: &model.SCGOpCallPkg{
										Ref: "debug-api",
									},
									Inputs: map[string]string{
										"appDir": "",
									},
								},
							},
						},
					},
				}

				objectUnderTest := newViewFactory(
					vioutil.New(osfs.New()),
					new(fakeValidator),
					format.NewYamlFormat(),
				)

				/* act */
				actualPackageView, err :=
					objectUnderTest.Construct(
						fmt.Sprintf("%v/../testdata/opspec-0.1.3/examples/nodejs/.opspec/debug", wd))
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualPackageView).To(Equal(expectedPackageView))

			})
		})
	})
})
