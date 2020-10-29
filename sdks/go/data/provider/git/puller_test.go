package git

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/provider/git/internal"
	. "github.com/opctl/opctl/sdks/go/data/provider/git/internal/fakes"
	"github.com/opctl/opctl/sdks/go/model"
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
			Context("git.PlainClone errors", func() {
				Context("err.Error() returns git.ErrRepositoryAlreadyExists", func() {
					It("shouldn't error", func() {

						/* arrange */
						providedPath := os.TempDir()

						fakeRefParser := new(FakeRefParser)
						fakeRefParser.ParseReturns(&internal.Ref{
							// some public repo that's relatively small
							Name:    "github.com/opspec-pkgs/_.op.create",
							Version: "3.2.0",
						}, nil)

						objectUnderTest := _puller{
							os:        new(ios.Fake),
							refParser: fakeRefParser,
						}

						/* act */
						firstErr := objectUnderTest.Pull(
							context.Background(),
							providedPath,
							"dummyRef",
							nil,
						)
						if nil != firstErr {
							panic(firstErr)
						}

						actualError := objectUnderTest.Pull(
							context.Background(),
							providedPath,
							"dummyRef",
							nil,
						)

						/* assert */
						Expect(actualError).To(BeNil())
					})
				})
				Context("err.Error() returns transport.ErrAuthenticationRequired error", func() {
					It("should return expected error", func() {

						/* arrange */
						providedPath := os.TempDir()

						ref := &internal.Ref{
							// use some private repo
							Name:    "github.com/Remitly/infra-ops",
							Version: "9.1.6",
						}

						fakeRefParser := new(FakeRefParser)
						fakeRefParser.ParseReturns(ref, nil)

						expectedError := model.ErrDataProviderAuthentication{}

						objectUnderTest := _puller{
							os:        new(ios.Fake),
							refParser: fakeRefParser,
						}

						/* act */
						actualError := objectUnderTest.Pull(
							context.Background(),
							providedPath,
							"dummyRef",
							nil,
						)

						/* assert */
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("err.Error() returns transport.ErrAuthorizationFailed error", func() {
					It("should return expected error", func() {

						/* arrange */
						providedPath := os.TempDir()

						ref := &internal.Ref{
							// use gitlab cuz github returns 404 not 403
							Name:    "gitlab.com/joetesterperson1/private",
							Version: "0.0.0",
						}

						fakeRefParser := new(FakeRefParser)
						fakeRefParser.ParseReturns(ref, nil)

						expectedError := model.ErrDataProviderAuthorization{}

						objectUnderTest := _puller{
							os:        new(ios.Fake),
							refParser: fakeRefParser,
						}

						/* act */
						actualError := objectUnderTest.Pull(
							context.Background(),
							providedPath,
							"dummyDataRef",
							&model.Creds{
								Username: "joetesterperson",
								Password: "MWgQpun9TWUx2iFQctyJ",
							},
						)

						/* assert */
						Expect(actualError).To(Equal(expectedError))
					})
				})
				Context("err.Error() returns other error", func() {
					It("should return error", func() {

						/* arrange */
						providedPath := os.TempDir()
						ref := &internal.Ref{
							Name:    "dummyDataRef",
							Version: "0.0.0",
						}

						expectedMsg := fmt.Sprintf(`Get "https://%s/info/refs?service=git-upload-pack": dial tcp: lookup dummyDataRef on 127.0.0.11:53: no such host`, ref.Name)

						fakeRefParser := new(FakeRefParser)
						fakeRefParser.ParseReturns(ref, nil)

						objectUnderTest := _puller{
							os:        new(ios.Fake),
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
						Expect(actualError.Error()).To(Equal(expectedMsg))
					})
				})
			})
			Context("git.PlainClone doesn't error", func() {
				It("should remove pkg '.git' sub dir & return errors", func() {

					/* arrange */
					providedPath := os.TempDir()

					ref := &internal.Ref{
						// some public repo that's relatively small
						Name:    "github.com/opspec-pkgs/_.op.create",
						Version: "3.3.1",
					}

					fakeRefParser := new(FakeRefParser)
					fakeRefParser.ParseReturns(ref, nil)

					expectedPath := filepath.Join(ref.ToPath(providedPath), ".git")

					fakeOS := new(ios.Fake)
					expectedError := errors.New("dummyError")
					fakeOS.RemoveAllReturns(expectedError)

					objectUnderTest := _puller{
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
					Expect(actualError).To(Equal(expectedError))
					Expect(fakeOS.RemoveAllArgsForCall(0)).To(Equal(expectedPath))
				})
			})
		})
	})
})
