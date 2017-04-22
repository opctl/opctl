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
		It("should call localResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedGetReq := &GetReq{
				BasePath: "dummyBasePath",
				PkgRef:   "dummyPkgRef",
			}

			resolvedPkgRef := "dummyPath"

			fakeLocalResolver := new(fakeLocalResolver)
			fakeLocalResolver.ResolveReturns(resolvedPkgRef, true)

			objectUnderTest := _getter{
				manifestUnmarshaller: new(fakeManifestUnmarshaller),
				localResolver:        fakeLocalResolver,
			}

			/* act */
			objectUnderTest.Get(providedGetReq)

			/* assert */
			actualBasePath, actualPkgRef := fakeLocalResolver.ResolveArgsForCall(0)

			Expect(actualBasePath).To(Equal(providedGetReq.BasePath))
			Expect(actualPkgRef).To(Equal(providedGetReq.PkgRef))
		})
		Context("is local pkg", func() {
			It("should call manifestUnmarshaller.Unmarshal w/ expected args", func() {
				/* arrange */
				providedGetReq := &GetReq{
					PkgRef: "dummyPkgRef",
				}

				resolvedPkgRef := "dummyPath"

				fakeLocalResolver := new(fakeLocalResolver)
				fakeLocalResolver.ResolveReturns(resolvedPkgRef, true)

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)

				objectUnderTest := _getter{
					manifestUnmarshaller: fakeManifestUnmarshaller,
					localResolver:        fakeLocalResolver,
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

				fakeLocalResolver := new(fakeLocalResolver)
				fakeLocalResolver.ResolveReturns(resolvedPkgRef, true)

				expectedView := &model.PkgManifest{Name: "dummyName"}
				expectedErr := errors.New("dummyError")

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
				fakeManifestUnmarshaller.UnmarshalReturns(expectedView, expectedErr)

				objectUnderTest := _getter{
					manifestUnmarshaller: fakeManifestUnmarshaller,
					localResolver:        fakeLocalResolver,
				}

				/* act */
				actualView, actualErr := objectUnderTest.Get(providedGetReq)

				/* assert */
				Expect(actualView).To(Equal(expectedView))
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("isn't local pkg", func() {
			It("should call git.PlainClone w/ expected args", func() {

				/* arrange */
				providedGetReq := &GetReq{
					PkgRef: "dummyPkgRef#0.0.0",
				}

				stringParts := strings.Split(providedGetReq.PkgRef, "#")
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
					git:                  fakeGit,
					manifestUnmarshaller: new(fakeManifestUnmarshaller),
					localResolver:        new(fakeLocalResolver),
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

					stringParts := strings.Split(providedGetReq.PkgRef, "#")
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

					expectedError := errors.New("dummyError")

					fakeGit := new(vgit.Fake)
					fakeGit.PlainCloneReturns(expectedError)

					objectUnderTest := _getter{
						fs:                   fakeFS,
						git:                  fakeGit,
						manifestUnmarshaller: new(fakeManifestUnmarshaller),
						localResolver:        new(fakeLocalResolver),
					}

					/* act */
					_, actualError := objectUnderTest.Get(providedGetReq)

					/* assert */
					Expect(fakeFS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("git.PlainClone doesn't error", func() {
				It("should return result of manifestUnmarshaller.Unmarshal", func() {
					/* arrange */
					providedGetReq := &GetReq{
						PkgRef: "dummyPkgRef#0.0.0",
					}

					expectedView := &model.PkgManifest{Name: "dummyName"}
					expectedErr := errors.New("dummyError")

					fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
					fakeManifestUnmarshaller.UnmarshalReturns(expectedView, expectedErr)

					objectUnderTest := _getter{
						git:                  new(vgit.Fake),
						manifestUnmarshaller: fakeManifestUnmarshaller,
						localResolver:        new(fakeLocalResolver),
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
