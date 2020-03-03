package git

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	igit "github.com/golang-interfaces/gopkg.in-src-d-go-git.v4"
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/provider/git/internal"
	. "github.com/opctl/opctl/sdks/go/data/provider/git/internal/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var _ = Context("puller", func() {
	Context("Pull", func() {
		It("should call refParser.Parse w/ expected args", func() {
			/* arrange */
			providedDataRef := "dummyDataRef"

			fakeRefParser := new(FakeRefParser)
			// error to trigger immediate return
			fakeRefParser.ParseReturns(nil, errors.New("dummyError"))

			objectUnderTest := _puller{
				refParser: fakeRefParser,
			}

			/* act */
			objectUnderTest.Pull(
				context.Background(),
				"dummyPath",
				providedDataRef,
				nil,
			)

			/* assert */
			Expect(fakeRefParser.ParseArgsForCall(0)).To(Equal(providedDataRef))
		})
		Context("refParser.Parse errs", func() {
			It("should return error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")

				fakeRefParser := new(FakeRefParser)
				fakeRefParser.ParseReturns(nil, expectedError)

				objectUnderTest := _puller{
					refParser: fakeRefParser,
				}

				/* act */
				actualError := objectUnderTest.Pull(
					context.Background(),
					"dummyPath",
					"dummyDataRef",
					nil,
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("refParser.Parse doesn't err", func() {
			It("should call git.PlainClone w/ expected args", func() {

				/* arrange */
				providedCtx := context.Background()
				providedPath := "dummyPath"
				providedPullCreds := &model.PullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				}

				ref := &internal.Ref{
					Name:    "dummyDataRef",
					Version: "0.0.0",
				}

				fakeRefParser := new(FakeRefParser)
				fakeRefParser.ParseReturns(ref, nil)

				expectedCloneOptions := &git.CloneOptions{
					Auth: &http.BasicAuth{
						Username: providedPullCreds.Username,
						Password: providedPullCreds.Password,
					},
					URL:           fmt.Sprintf("https://%v", ref.Name),
					ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%v", ref.Version)),
					Depth:         1,
					Progress:      os.Stdout,
				}

				fakeGit := new(igit.Fake)

				objectUnderTest := _puller{
					git:       fakeGit,
					os:        new(ios.Fake),
					refParser: fakeRefParser,
				}

				/* act */
				objectUnderTest.Pull(
					providedCtx,
					providedPath,
					"dummyDataRef",
					providedPullCreds,
				)

				/* assert */
				actualPath,
					actualIsBare,
					actualCloneOptions := fakeGit.PlainCloneArgsForCall(0)

				Expect(actualPath).To(Equal(ref.ToPath(providedPath)))
				Expect(actualIsBare).To(BeFalse())
				Expect(actualCloneOptions).To(Equal(expectedCloneOptions))
			})
			Context("git.PlainClone errors", func() {
				Context("err.Error() returns git.ErrRepositoryAlreadyExists", func() {
					It("shouldn't call fs.RemoveAll or error", func() {

						/* arrange */
						fakeRefParser := new(FakeRefParser)
						fakeRefParser.ParseReturns(&internal.Ref{}, nil)

						fakeOS := new(ios.Fake)

						fakeGit := new(igit.Fake)
						fakeGit.PlainCloneReturns(nil, git.ErrRepositoryAlreadyExists)

						objectUnderTest := _puller{
							git:       fakeGit,
							os:        fakeOS,
							refParser: fakeRefParser,
						}

						/* act */
						actualError := objectUnderTest.Pull(
							context.Background(),
							"dummyPath",
							"dummyDataRef",
							nil,
						)

						/* assert */
						Expect(actualError).To(BeNil())
						Expect(fakeOS.RemoveAllCallCount()).To(Equal(0))
					})
				})
				Context("err.Error() returns transport.ErrAuthenticationRequired error", func() {
					It("should call fs.RemoveAll w/ expected args & return expected error", func() {

						/* arrange */
						providedPath := "dummyPath"

						ref := &internal.Ref{
							Name:    "dummyDataRef",
							Version: "0.0.0",
						}

						fakeRefParser := new(FakeRefParser)
						fakeRefParser.ParseReturns(ref, nil)

						expectedPath := filepath.Join(
							providedPath,
							fmt.Sprintf("%v#%v", ref.Name, ref.Version),
						)

						fakeOS := new(ios.Fake)
						expectedError := model.ErrDataProviderAuthentication{}

						fakeGit := new(igit.Fake)
						fakeGit.PlainCloneReturns(nil, transport.ErrAuthenticationRequired)

						objectUnderTest := _puller{
							git:       fakeGit,
							os:        fakeOS,
							refParser: fakeRefParser,
						}

						/* act */
						actualError := objectUnderTest.Pull(
							context.Background(),
							providedPath,
							"dummyDataRef",
							nil,
						)

						/* assert */
						Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("err.Error() returns transport.ErrAuthorizationFailed error", func() {
					It("should call fs.RemoveAll w/ expected args & return expected error", func() {

						/* arrange */
						providedPath := "dummyPath"

						ref := &internal.Ref{
							Name:    "dummyDataRef",
							Version: "0.0.0",
						}

						fakeRefParser := new(FakeRefParser)
						fakeRefParser.ParseReturns(ref, nil)

						expectedPath := filepath.Join(
							providedPath,
							fmt.Sprintf("%v#%v", ref.Name, ref.Version),
						)

						fakeOS := new(ios.Fake)
						expectedError := model.ErrDataProviderAuthorization{}

						fakeGit := new(igit.Fake)
						fakeGit.PlainCloneReturns(nil, transport.ErrAuthorizationFailed)

						objectUnderTest := _puller{
							git:       fakeGit,
							os:        fakeOS,
							refParser: fakeRefParser,
						}

						/* act */
						actualError := objectUnderTest.Pull(
							context.Background(),
							providedPath,
							"dummyDataRef",
							nil,
						)

						/* assert */
						Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("err.Error() returns other error", func() {
					It("should call fs.RemoveAll w/ expected args & return error", func() {

						/* arrange */
						providedPath := "dummypath"
						ref := &internal.Ref{
							Name:    "dummyDataRef",
							Version: "0.0.0",
						}

						fakeRefParser := new(FakeRefParser)
						fakeRefParser.ParseReturns(ref, nil)

						fakeOS := new(ios.Fake)

						expectedError := errors.New("dummyError")

						fakeGit := new(igit.Fake)
						fakeGit.PlainCloneReturns(nil, expectedError)

						objectUnderTest := _puller{
							git:       fakeGit,
							os:        fakeOS,
							refParser: fakeRefParser,
						}

						/* act */
						actualError := objectUnderTest.Pull(
							context.Background(),
							providedPath,
							"dummyDataRef",
							nil,
						)

						/* assert */
						Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(ref.ToPath(providedPath)))
						Expect(actualError).To(Equal(expectedError))
					})
				})
			})
			Context("git.PlainClone doesn't error", func() {
				It("shouldn't err", func() {
					/* arrange */
					fakeRefParser := new(FakeRefParser)
					fakeRefParser.ParseReturns(&internal.Ref{}, nil)

					objectUnderTest := _puller{
						git:       new(igit.Fake),
						os:        new(ios.Fake),
						refParser: fakeRefParser,
					}

					/* act */
					actualErr := objectUnderTest.Pull(
						context.Background(),
						"dummyPath",
						"dummyDataRef",
						nil,
					)

					/* assert */
					Expect(actualErr).To(BeNil())
				})

				It("should remove pkg '.git' sub dir & return errors", func() {

					/* arrange */
					providedPath := "dummypath"

					ref := &internal.Ref{
						Name:    "dummyDataRef",
						Version: "0.0.0",
					}

					fakeRefParser := new(FakeRefParser)
					fakeRefParser.ParseReturns(ref, nil)

					expectedPath := filepath.Join(ref.ToPath(providedPath), ".git")

					fakeOS := new(ios.Fake)
					expectedError := errors.New("dummyError")
					fakeOS.RemoveAllReturns(expectedError)

					objectUnderTest := _puller{
						git:       new(igit.Fake),
						os:        fakeOS,
						refParser: fakeRefParser,
					}

					/* act */
					actualError := objectUnderTest.Pull(
						context.Background(),
						providedPath,
						"dummyDataRef",
						nil,
					)

					/* assert */
					Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
					Expect(actualError).To(Equal(expectedError))
				})
			})
		})
	})
})
