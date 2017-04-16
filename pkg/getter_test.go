package pkg

import (
	"errors"
	"fmt"
	"github.com/appdataspec/sdk-golang/pkg/appdatapath"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/vgit"
	"github.com/virtual-go/fs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path"
	"strings"
)

var _ = Describe("Getter", func() {
	Context("newGetter()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(newGetter(nil, nil, nil)).Should(Not(BeNil()))
		})
	})
	Context("Get", func() {
		It("should call refResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedGetReq := &GetReq{
				PkgRef: "dummyPkgRef",
			}

			resolvedPkgRef := "dummyPath"

			fakeRefResolver := new(fakeRefResolver)
			fakeRefResolver.ResolveReturns(resolvedPkgRef)

			objectUnderTest := _getter{
				fs:                   new(fs.Fake),
				manifestUnmarshaller: new(fakeManifestUnmarshaller),
				refResolver:          fakeRefResolver,
			}

			/* act */
			objectUnderTest.Get(providedGetReq)

			/* assert */
			Expect(fakeRefResolver.ResolveArgsForCall(0)).To(Equal(providedGetReq.PkgRef))
		})
		It("should call fs.Stat w/ expected args", func() {
			/* arrange */
			providedGetReq := &GetReq{
				PkgRef: "dummyPkgRef",
			}

			resolvedPkgRef := "dummyPath"

			fakeRefResolver := new(fakeRefResolver)
			fakeRefResolver.ResolveReturns(resolvedPkgRef)

			fakeFS := new(fs.Fake)
			fakeFS.StatReturns(nil, nil)

			objectUnderTest := _getter{
				fs:                   fakeFS,
				manifestUnmarshaller: new(fakeManifestUnmarshaller),
				refResolver:          fakeRefResolver,
			}

			/* act */
			objectUnderTest.Get(providedGetReq)

			/* assert */
			Expect(fakeFS.StatArgsForCall(0)).To(Equal(resolvedPkgRef))
		})
		Context("is embedded pkg", func() {
			It("should call manifestUnmarshaller.Unmarshal w/ expected args", func() {
				/* arrange */
				providedGetReq := &GetReq{
					PkgRef: "dummyPkgRef",
				}

				resolvedPkgRef := "dummyPath"

				fakeRefResolver := new(fakeRefResolver)
				fakeRefResolver.ResolveReturns(resolvedPkgRef)

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)

				objectUnderTest := _getter{
					fs:                   new(fs.Fake),
					manifestUnmarshaller: fakeManifestUnmarshaller,
					refResolver:          fakeRefResolver,
				}

				/* act */
				objectUnderTest.Get(providedGetReq)

				/* assert */
				Expect(fakeManifestUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(resolvedPkgRef))
			})
			It("should return result of manifestUnmarshaller.Unmarshal", func() {
				/* arrange */
				providedGetReq := &GetReq{
					PkgRef: "dummyPkgRef",
				}

				resolvedPkgRef := "dummyPath"

				fakeRefResolver := new(fakeRefResolver)
				fakeRefResolver.ResolveReturns(resolvedPkgRef)

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

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
				fakeManifestUnmarshaller.UnmarshalReturns(expectedView, expectedErr)

				objectUnderTest := _getter{
					fs:                   new(fs.Fake),
					manifestUnmarshaller: fakeManifestUnmarshaller,
					refResolver:          fakeRefResolver,
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
						PkgRef: "dummyPkgRef",
					}

					resolvedPkgRef := "dummyResolvedPkgRef#0.0.0"

					fakeRefResolver := new(fakeRefResolver)
					fakeRefResolver.ResolveReturns(resolvedPkgRef)

					stringParts := strings.Split(resolvedPkgRef, "#")
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
						refResolver:          fakeRefResolver,
					}

					/* act */
					objectUnderTest.Get(providedGetReq)

					/* assert */
					Expect(fakeManifestUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(expectedPkgPath))
				})
				It("should return result of manifestUnmarshaller.Unmarshal", func() {
					/* arrange */
					providedGetReq := &GetReq{
						PkgRef: "dummyPkgRef",
					}

					resolvedPkgRef := "dummyResolvedPkgRef#0.0.0"

					fakeRefResolver := new(fakeRefResolver)
					fakeRefResolver.ResolveReturns(resolvedPkgRef)

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
						refResolver:          fakeRefResolver,
					}

					/* act */
					actualView, actualErr := objectUnderTest.Get(providedGetReq)

					/* assert */
					Expect(actualView).To(Equal(expectedView))
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("isn't cached", func() {
				It("should call git.PlainClone w/ expected args", func() {

					/* arrange */
					providedGetReq := &GetReq{
						PkgRef: "dummyPkgRef#0.0.0",
					}

					resolvedPkgRef := "dummyResolvedPkgRef#0.0.0"

					fakeRefResolver := new(fakeRefResolver)
					fakeRefResolver.ResolveReturns(resolvedPkgRef)

					stringParts := strings.Split(resolvedPkgRef, "#")
					repoName := stringParts[0]
					repoRefName := stringParts[1]

					expectedPath := path.Join(
						appdatapath.New().PerUser(),
						"opspec",
						"cache",
						"pkgs",
						repoName,
						repoRefName,
					)

					expectedIsBare := false

					expectedCloneOptions := &git.CloneOptions{
						URL:           fmt.Sprintf("https://%v", repoName),
						ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", repoRefName)),
						Depth:         1,
						Progress:      os.Stdout,
					}

					fakeFS := new(fs.Fake)
					fakeFS.StatReturnsOnCall(0, nil, errors.New(""))
					fakeFS.StatReturnsOnCall(1, nil, errors.New(""))

					fakeGit := new(vgit.Fake)

					objectUnderTest := _getter{
						fs:                   fakeFS,
						git:                  fakeGit,
						manifestUnmarshaller: new(fakeManifestUnmarshaller),
						refResolver:          fakeRefResolver,
					}

					/* act */
					objectUnderTest.Get(providedGetReq)

					/* assert */
					actualPath,
						actualIsBare,
						actualCloneOptions := fakeGit.PlainCloneArgsForCall(0)

					Expect(actualPath).To(Equal(expectedPath))
					Expect(actualIsBare).To(Equal(expectedIsBare))
					Expect(actualCloneOptions).To(Equal(expectedCloneOptions))
				})
				Context("git.PlainClone errors", func() {
					It("should call fs.RemoveAll w/ expected args & return error", func() {

						/* arrange */
						providedGetReq := &GetReq{
							PkgRef: "dummyPkgRef#0.0.0",
						}

						resolvedPkgRef := "dummyResolvedPkgRef#0.0.0"

						fakeRefResolver := new(fakeRefResolver)
						fakeRefResolver.ResolveReturns(resolvedPkgRef)

						stringParts := strings.Split(resolvedPkgRef, "#")
						repoName := stringParts[0]
						repoRefName := stringParts[1]

						expectedPath := path.Join(
							appdatapath.New().PerUser(),
							"opspec",
							"cache",
							"pkgs",
							repoName,
							repoRefName,
						)

						fakeFS := new(fs.Fake)
						fakeFS.StatReturnsOnCall(0, nil, errors.New(""))
						fakeFS.StatReturnsOnCall(1, nil, errors.New(""))

						expectedError := errors.New("dummyError")

						fakeGit := new(vgit.Fake)
						fakeGit.PlainCloneReturns(expectedError)

						objectUnderTest := _getter{
							fs:                   fakeFS,
							git:                  fakeGit,
							manifestUnmarshaller: new(fakeManifestUnmarshaller),
							refResolver:          fakeRefResolver,
						}

						/* act */
						_, actualError := objectUnderTest.Get(providedGetReq)

						/* assert */
						Expect(fakeFS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
						Expect(actualError).To(Equal(expectedError))
					})
				})
			})
		})
	})
})
