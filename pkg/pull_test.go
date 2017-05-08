package pkg

import (
	"errors"
	"fmt"
	"github.com/appdataspec/sdk-golang/appdatapath"
	"github.com/golang-interfaces/gopkg.in-src-d-go-git.v4"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"path"
	"strings"
)

var _ = Describe("Pkg", func() {
	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		panic(err)
	}
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
					perUserAppDataPath,
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

				fakeGit := new(igit.Fake)

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

						fakeOS := new(ios.Fake)

						fakeGit := new(igit.Fake)
						fakeGit.PlainCloneReturns(nil, git.ErrRepositoryAlreadyExists)

						objectUnderTest := pkg{
							os:  fakeOS,
							git: fakeGit,
						}

						/* act */
						actualError := objectUnderTest.Pull(providedPkgRef, nil)

						/* assert */
						Expect(actualError).To(BeNil())
						Expect(fakeOS.RemoveAllCallCount()).To(Equal(0))
					})
				})
				Context("err.Error() returns transport.ErrAuthorizationRequired error", func() {
					It("should call fs.RemoveAll w/ expected args & return expected error", func() {

						/* arrange */
						providedPkgRef := "dummyPkgRef#0.0.0"

						stringParts := strings.Split(providedPkgRef, "#")
						repoName := stringParts[0]
						repoRefName := stringParts[1]

						expectedPath := path.Join(
							perUserAppDataPath,
							"opspec",
							"cache",
							"pkgs",
							repoName,
							repoRefName,
						)

						fakeOS := new(ios.Fake)

						expectedError := ErrAuthenticationFailed{}

						fakeGit := new(igit.Fake)
						fakeGit.PlainCloneReturns(nil, transport.ErrAuthorizationRequired)

						objectUnderTest := pkg{
							os:  fakeOS,
							git: fakeGit,
						}

						/* act */
						actualError := objectUnderTest.Pull(providedPkgRef, nil)

						/* assert */
						Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("err.Error() returns other error", func() {
					It("should call fs.RemoveAll w/ expected args & return error", func() {

						/* arrange */
						providedPkgRef := "dummyPkgRef#0.0.0"

						stringParts := strings.Split(providedPkgRef, "#")
						repoName := stringParts[0]
						repoRefName := stringParts[1]

						expectedPath := path.Join(
							perUserAppDataPath,
							"opspec",
							"cache",
							"pkgs",
							repoName,
							repoRefName,
						)

						fakeOS := new(ios.Fake)

						expectedError := errors.New("dummyError")

						fakeGit := new(igit.Fake)
						fakeGit.PlainCloneReturns(nil, expectedError)

						objectUnderTest := pkg{
							os:  fakeOS,
							git: fakeGit,
						}

						/* act */
						actualError := objectUnderTest.Pull(providedPkgRef, nil)

						/* assert */
						Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
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
						git: new(igit.Fake),
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
