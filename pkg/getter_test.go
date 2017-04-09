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

			expectedPkgPath := path.Join(providedGetReq.Path, RepoDirName, providedGetReq.PkgRef)

			fakeFS := new(fs.Fake)
			fakeFS.StatReturns(nil, nil)

			objectUnderTest := _getter{
				fs:                   fakeFS,
				manifestUnmarshaller: new(fakeManifestUnmarshaller),
			}

			/* act */
			objectUnderTest.Get(providedGetReq)

			/* assert */
			Expect(fakeFS.StatArgsForCall(0)).To(Equal(expectedPkgPath))
		})
		Context("is embedded pkg", func() {
			It("should call manifestUnmarshaller.Unmarshal w/ expected args", func() {
				/* arrange */
				providedGetReq := &GetReq{
					Path:   "/dummyPath",
					PkgRef: "dummyPkgRef",
				}

				expectedPkgPath := path.Join(providedGetReq.Path, RepoDirName, providedGetReq.PkgRef)

				fakeFS := new(fs.Fake)
				fakeFS.StatReturns(nil, nil)

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)

				objectUnderTest := _getter{
					fs:                   fakeFS,
					manifestUnmarshaller: fakeManifestUnmarshaller,
				}

				/* act */
				objectUnderTest.Get(providedGetReq)

				/* assert */
				Expect(fakeManifestUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(expectedPkgPath))
			})
			It("should return result of manifestUnmarshaller.Unmarshal", func() {
				/* arrange */
				providedGetReq := &GetReq{
					Path:   "/dummyPath",
					PkgRef: "dummyPkgRef",
				}

				expectedView := &model.PkgManifest{
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

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
				fakeManifestUnmarshaller.UnmarshalReturns(expectedView, expectedErr)

				objectUnderTest := _getter{
					fs:                   fakeFS,
					manifestUnmarshaller: fakeManifestUnmarshaller,
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
				It("should call manifestUnmarshaller.Unmarshal w/ expected args", func() {
					/* arrange */
					providedGetReq := &GetReq{
						Path:   "/dummyPath",
						PkgRef: "dummyPkgRef#0.0.0",
					}

					stringParts := strings.Split(providedGetReq.PkgRef, "#")
					repoName := stringParts[0]
					repoRefName := stringParts[1]

					expectedPkgPath := path.Join(
						appdatapath.New().PerUser(),
						"opspec",
						"cache",
						"pkgs",
						repoName,
						repoRefName,
					)

					fakeFS := new(fs.Fake)
					fakeFS.StatReturnsOnCall(0, nil, errors.New(""))
					fakeFS.StatReturnsOnCall(1, nil, nil)

					fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)

					objectUnderTest := _getter{
						fs:                   fakeFS,
						manifestUnmarshaller: fakeManifestUnmarshaller,
					}

					/* act */
					objectUnderTest.Get(providedGetReq)

					/* assert */
					Expect(fakeManifestUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(expectedPkgPath))
				})
				It("should return result of manifestUnmarshaller.Unmarshal", func() {
					/* arrange */
					providedGetReq := &GetReq{
						Path:   "/dummyPath",
						PkgRef: "dummyPkgRef#0.0.0",
					}

					expectedView := &model.PkgManifest{
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

					fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
					fakeManifestUnmarshaller.UnmarshalReturns(expectedView, expectedErr)

					objectUnderTest := _getter{
						fs:                   fakeFS,
						manifestUnmarshaller: fakeManifestUnmarshaller,
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
