package pkg

import (
	"errors"
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/virtual-go/fs"
	"path"
	"strings"
)

var _ = Describe("Getter", func() {
	Context("newGetter()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(newGetter(nil, nil)).Should(Not(BeNil()))
		})
	})
	Context("Get", func() {
		It("should call fs.Stat w/ expected args", func() {
			/* arrange */
			providedGetReq := &GetReq{
				Path:   "/dummyPath",
				PkgRef: "dummyPkgRef",
			}

			expectedPkgPath := path.Join(providedGetReq.Path, ".opspec", providedGetReq.PkgRef)

			fakeFS := new(fs.Fake)
			fakeFS.StatReturns(nil, nil)

			objectUnderTest := _getter{
				fs:          fakeFS,
				viewFactory: new(fakeViewFactory),
			}

			/* act */
			objectUnderTest.Get(providedGetReq)

			/* assert */
			Expect(fakeFS.StatArgsForCall(0)).To(Equal(expectedPkgPath))
		})
		Context("is embedded pkg", func() {
			It("should call viewFactory.Construct w/ expected args", func() {
				/* arrange */
				providedGetReq := &GetReq{
					Path:   "/dummyPath",
					PkgRef: "dummyPkgRef",
				}

				expectedPkgPath := path.Join(providedGetReq.Path, ".opspec", providedGetReq.PkgRef)

				fakeFS := new(fs.Fake)
				fakeFS.StatReturns(nil, nil)

				fakeViewFactory := new(fakeViewFactory)

				objectUnderTest := _getter{
					fs:          fakeFS,
					viewFactory: fakeViewFactory,
				}

				/* act */
				objectUnderTest.Get(providedGetReq)

				/* assert */
				Expect(fakeViewFactory.ConstructArgsForCall(0)).To(Equal(expectedPkgPath))
			})
			It("should return result of viewFactory.Construct", func() {
				/* arrange */
				providedGetReq := &GetReq{
					Path:   "/dummyPath",
					PkgRef: "dummyPkgRef",
				}

				expectedView := &model.PackageView{
					Description: "dummyDescription",
					Inputs:      map[string]*model.Param{},
					Outputs:     map[string]*model.Param{},
					Name:        "dummyName",
					Run: &model.SCG{
						Op: &model.SCGOpCall{
							Pkg: &model.SCGOpCallPkg{
								Ref: "dummyPkgRef",
							},
						},
					},
					Version: "",
				}
				expectedErr := errors.New("dummyError")

				fakeFS := new(fs.Fake)
				fakeFS.StatReturns(nil, nil)

				fakeViewFactory := new(fakeViewFactory)
				fakeViewFactory.ConstructReturns(expectedView, expectedErr)

				objectUnderTest := _getter{
					fs:          fakeFS,
					viewFactory: fakeViewFactory,
				}

				/* act */
				actualView, actualErr := objectUnderTest.Get(providedGetReq)

				/* assert */
				Expect(actualView).To(Equal(expectedView))
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("isn't embedded pkg", func() {
			Context("is cached", func() {
				It("should call viewFactory.Construct w/ expected args", func() {
					/* arrange */
					providedGetReq := &GetReq{
						Path:   "/dummyPath",
						PkgRef: "dummyPkgRef#000",
					}

					stringParts := strings.Split(providedGetReq.PkgRef, "#")
					repoName := stringParts[0]
					repoRefName := stringParts[1]

					expectedPkgPath := path.Join(
						appdatapath.New().PerUser(),
						".opspec",
						"cache",
						"pkgs",
						repoName,
						repoRefName,
					)

					fakeFS := new(fs.Fake)
					fakeFS.StatReturnsOnCall(0, nil, errors.New(""))
					fakeFS.StatReturnsOnCall(1, nil, nil)

					fakeViewFactory := new(fakeViewFactory)

					objectUnderTest := _getter{
						fs:          fakeFS,
						viewFactory: fakeViewFactory,
					}

					/* act */
					objectUnderTest.Get(providedGetReq)

					/* assert */
					Expect(fakeViewFactory.ConstructArgsForCall(0)).To(Equal(expectedPkgPath))
				})
				It("should return result of viewFactory.Construct", func() {
					/* arrange */
					providedGetReq := &GetReq{
						Path:   "/dummyPath",
						PkgRef: "dummyPkgRef#000",
					}

					expectedView := &model.PackageView{
						Description: "dummyDescription",
						Inputs:      map[string]*model.Param{},
						Outputs:     map[string]*model.Param{},
						Name:        "dummyName",
						Run: &model.SCG{
							Op: &model.SCGOpCall{
								Pkg: &model.SCGOpCallPkg{
									Ref: "dummyPkgRef",
								},
							},
						},
						Version: "",
					}
					expectedErr := errors.New("dummyError")

					fakeFS := new(fs.Fake)
					fakeFS.StatReturnsOnCall(0, nil, errors.New(""))
					fakeFS.StatReturnsOnCall(1, nil, nil)

					fakeViewFactory := new(fakeViewFactory)
					fakeViewFactory.ConstructReturns(expectedView, expectedErr)

					objectUnderTest := _getter{
						fs:          fakeFS,
						viewFactory: fakeViewFactory,
					}

					/* act */
					actualView, actualErr := objectUnderTest.Get(providedGetReq)

					/* assert */
					Expect(actualView).To(Equal(expectedView))
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
		})
	})
})
