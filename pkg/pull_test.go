package pkg

import (
	"errors"
	"fmt"
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
	"path/filepath"
)

var _ = Describe("Pkg", func() {
	Context("Pull", func() {
		It("should call git.PlainClone w/ expected args", func() {

			/* arrange */
			providedPath := "dummyPath"
			providedPkgRef := &PkgRef{
				FullyQualifiedName: "dummyPkgRef",
				Version:            "0.0.0",
			}
			providedOpts := &PullOpts{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			expectedPath := filepath.Join(
				providedPath,
				providedPkgRef.FullyQualifiedName,
				providedPkgRef.Version,
			)

			expectedIsBare := false

			expectedCloneOptions := &git.CloneOptions{
				Auth:          http.NewBasicAuth(providedOpts.Username, providedOpts.Password),
				URL:           fmt.Sprintf("https://%v", providedPkgRef.FullyQualifiedName),
				ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", providedPkgRef.Version)),
				Depth:         1,
				Progress:      os.Stdout,
			}

			fakeGit := new(igit.Fake)

			objectUnderTest := _Pkg{
				git: fakeGit,
			}

			/* act */
			objectUnderTest.Pull(providedPath, providedPkgRef, providedOpts)

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
					providedPkgRef := &PkgRef{
						FullyQualifiedName: "dummyPkgRef",
						Version:            "0.0.0",
					}

					fakeOS := new(ios.Fake)

					fakeGit := new(igit.Fake)
					fakeGit.PlainCloneReturns(nil, git.ErrRepositoryAlreadyExists)

					objectUnderTest := _Pkg{
						os:  fakeOS,
						git: fakeGit,
					}

					/* act */
					actualError := objectUnderTest.Pull("dummyPath", providedPkgRef, nil)

					/* assert */
					Expect(actualError).To(BeNil())
					Expect(fakeOS.RemoveAllCallCount()).To(Equal(0))
				})
			})
			Context("err.Error() returns transport.ErrAuthorizationRequired error", func() {
				It("should call fs.RemoveAll w/ expected args & return expected error", func() {

					/* arrange */
					providedPath := "dummyPath"
					providedPkgRef := &PkgRef{
						FullyQualifiedName: "dummyPkgRef",
						Version:            "0.0.0",
					}

					expectedPath := filepath.Join(
						providedPath,
						providedPkgRef.FullyQualifiedName,
						providedPkgRef.Version,
					)

					fakeOS := new(ios.Fake)

					expectedError := ErrAuthenticationFailed{}

					fakeGit := new(igit.Fake)
					fakeGit.PlainCloneReturns(nil, transport.ErrAuthorizationRequired)

					objectUnderTest := _Pkg{
						os:  fakeOS,
						git: fakeGit,
					}

					/* act */
					actualError := objectUnderTest.Pull(providedPath, providedPkgRef, nil)

					/* assert */
					Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("err.Error() returns other error", func() {
				It("should call fs.RemoveAll w/ expected args & return error", func() {

					/* arrange */
					providedPath := "dummypath"
					providedPkgRef := &PkgRef{
						FullyQualifiedName: "dummyPkgRef",
						Version:            "0.0.0",
					}

					expectedPath := filepath.Join(
						providedPath,
						providedPkgRef.FullyQualifiedName,
						providedPkgRef.Version,
					)

					fakeOS := new(ios.Fake)

					expectedError := errors.New("dummyError")

					fakeGit := new(igit.Fake)
					fakeGit.PlainCloneReturns(nil, expectedError)

					objectUnderTest := _Pkg{
						os:  fakeOS,
						git: fakeGit,
					}

					/* act */
					actualError := objectUnderTest.Pull(providedPath, providedPkgRef, nil)

					/* assert */
					Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
					Expect(actualError).To(Equal(expectedError))
				})
			})
		})
		Context("git.PlainClone doesn't error", func() {
			It("shouldn't err", func() {
				/* arrange */
				providedPkgRef := &PkgRef{
					FullyQualifiedName: "dummyPkgRef",
					Version:            "0.0.0",
				}

				expectedView := &model.PkgManifest{Name: "dummyName"}
				expectedErr := errors.New("dummyError")

				fakeManifestUnmarshaller := new(fakeManifestUnmarshaller)
				fakeManifestUnmarshaller.UnmarshalReturns(expectedView, expectedErr)

				objectUnderTest := _Pkg{
					git: new(igit.Fake),
				}

				/* act */
				actualErr := objectUnderTest.Pull("dummyPath", providedPkgRef, nil)

				/* assert */
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
