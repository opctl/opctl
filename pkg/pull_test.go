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
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path"
	"strings"
)

var _ = Describe("Pkg", func() {
	Context("Pull", func() {
		Context("pkgRef invalid format", func() {
			It("should return expected result", func() {
				/* arrange */
				providedPkgRef := "notValid"
				expectedErr := fmt.Errorf(
					"Invalid remote pkgRef: '%v'. Valid remote pkgRef's are of the form: 'host/path#semver",
					providedPkgRef,
				)

				objectUnderTest := pkg{}

				/* act */
				actualError := objectUnderTest.Pull(providedPkgRef, nil)

				/* assert */
				Expect(actualError).To(Equal(expectedErr))
			})
		})
		Context("pkgRef valid format", func() {
			It("should call git.PlainClone w/ expected args", func() {

				/* arrange */
				providedPkgRef := "dummyPkgRef#0.0.0"
				providedOpts := &PullOpts{
					Username: "dummyUsername",
					Password: "dummyPassword",
				}

				stringParts := strings.Split(providedPkgRef, "#")
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
					Auth:          http.NewBasicAuth(providedOpts.Username, providedOpts.Password),
					URL:           fmt.Sprintf("https://%v", repoName),
					ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", repoRefName)),
					Depth:         1,
					Progress:      os.Stdout,
				}

				fakeGit := new(vgit.Fake)

				objectUnderTest := pkg{
					git: fakeGit,
				}

				/* act */
				objectUnderTest.Pull(providedPkgRef, providedOpts)

				/* assert */
				actualPath,
					actualIsBare,
					actualCloneOptions := fakeGit.PlainCloneArgsForCall(0)

				Expect(actualPath).To(Equal(expectedPath))
				Expect(actualIsBare).To(Equal(expectedIsBare))
				Expect(actualCloneOptions).To(Equal(expectedCloneOptions))
			})
			Context("git.PlainClone errors", func() {
				Context("err.Error() returns git.ErrRepositoryAlreadyExists", func() {
					It("shouldn't call fs.RemoveAll or error", func() {

						/* arrange */
						providedPkgRef := "dummyPkgRef#0.0.0"

						fakeFS := new(fs.Fake)

						fakeGit := new(vgit.Fake)
						fakeGit.PlainCloneReturns(git.ErrRepositoryAlreadyExists)

						objectUnderTest := pkg{
							fs:  fakeFS,
							git: fakeGit,
						}

						/* act */
						actualError := objectUnderTest.Pull(providedPkgRef, nil)

						/* assert */
						Expect(actualError).To(BeNil())
						Expect(fakeFS.RemoveAllCallCount()).To(Equal(0))
					})
				})
				Context("err.Error() doesn't return git.ErrRepositoryAlreadyExists", func() {
					It("should call fs.RemoveAll w/ expected args & return error", func() {

						/* arrange */
						providedPkgRef := "dummyPkgRef#0.0.0"

						stringParts := strings.Split(providedPkgRef, "#")
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

						objectUnderTest := pkg{
							fs:  fakeFS,
							git: fakeGit,
						}

						/* act */
						actualError := objectUnderTest.Pull(providedPkgRef, nil)

						/* assert */
						Expect(fakeFS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
						Expect(actualError).To(Equal(expectedError))
					})
				})
			})
			Context("git.PlainClone doesn't error", func() {
				It("shouldn't err", func() {
					/* arrange */
					providedPkgRef := "dummyPkgRef#0.0.0"

					expectedView := &model.PkgManifest{Name: "dummyName"}
					expectedErr := errors.New("dummyError")

					fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
					fakeManifestUnmarshaller.UnmarshalReturns(expectedView, expectedErr)

					objectUnderTest := pkg{
						git: new(vgit.Fake),
					}

					/* act */
					actualErr := objectUnderTest.Pull(providedPkgRef, nil)

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
