package pkg

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"time"
)

var _ = Context("gitProvider", func() {
	Context("TryResolve", func() {
		Context("called serially", func() {

			It("should call localFSProvider.TryResolve w/ expected args", func() {
				/* arrange */
				providedCtx := context.Background()
				providedPkgRef := "dummyPkgRef"

				fakeLocalFSProvider := new(FakeProvider)
				// err to trigger immediate return
				fakeLocalFSProvider.TryResolveReturns(nil, errors.New("dummyError"))

				objectUnderTest := gitProvider{
					localFSProvider: fakeLocalFSProvider,
				}

				/* act */
				objectUnderTest.TryResolve(
					providedCtx,
					providedPkgRef,
				)

				/* assert */
				actualCtx,
					actualPkgRef := fakeLocalFSProvider.TryResolveArgsForCall(0)

				Expect(actualCtx).To(Equal(providedCtx))
				Expect(actualPkgRef).To(Equal(providedPkgRef))
			})
			Context("localFSProvider.TryResolve errors", func() {
				It("should return err", func() {
					/* arrange */
					expectedErr := errors.New("dummyError")

					fakeLocalFSProvider := new(FakeProvider)
					// err to trigger immediate return
					fakeLocalFSProvider.TryResolveReturns(nil, expectedErr)

					objectUnderTest := gitProvider{
						localFSProvider: fakeLocalFSProvider,
					}

					/* act */
					_, actualError := objectUnderTest.TryResolve(
						context.Background(),
						"dummyPkgRef",
					)

					/* assert */
					Expect(actualError).To(Equal(expectedErr))
				})
			})
			Context("localFSProvider.TryResolve doesn't err", func() {
				Context("localFSProvider.TryResolve returns handle", func() {
					It("should return handle", func() {
						/* arrange */
						expectedHandle := new(FakeHandle)

						fakeLocalFSProvider := new(FakeProvider)
						// err to trigger immediate return
						fakeLocalFSProvider.TryResolveReturns(expectedHandle, nil)

						objectUnderTest := gitProvider{
							localFSProvider: fakeLocalFSProvider,
						}

						/* act */
						actualHandle, actualError := objectUnderTest.TryResolve(
							context.Background(),
							"dummyPkgRef",
						)

						/* assert */
						Expect(actualHandle).To(Equal(expectedHandle))
						Expect(actualError).To(BeNil())
					})
				})
				Context("FSProvider.TryResolve doesn't return a handle", func() {
					It("should call puller.Pull w/ expected args", func() {
						/* arrange */
						providedCtx := context.Background()
						providedPkgRef := "dummyPkgRef"
						basePath := "dummyBasePath"
						pullCreds := &model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"}

						fakePuller := new(fakePuller)
						// err to trigger immediate return
						fakePuller.PullReturns(errors.New("dummyError"))

						objectUnderTest := gitProvider{
							basePath:        basePath,
							localFSProvider: new(FakeProvider),
							pullCreds:       pullCreds,
							puller:          fakePuller,
						}

						/* act */
						objectUnderTest.TryResolve(
							providedCtx,
							providedPkgRef,
						)

						/* assert */
						actualCtx,
							actualBasePath,
							actualPkgRef,
							actualPullCreds := fakePuller.PullArgsForCall(0)
						Expect(actualCtx).To(Equal(providedCtx))
						Expect(actualBasePath).To(Equal(basePath))
						Expect(actualPkgRef).To(Equal(providedPkgRef))
						Expect(actualPullCreds).To(Equal(pullCreds))
					})
					Context("puller.Pull errors", func() {
						It("should return err", func() {
							/* arrange */
							expectedErr := errors.New("dummyError")

							fakePuller := new(fakePuller)
							// err to trigger immediate return
							fakePuller.PullReturns(expectedErr)

							objectUnderTest := gitProvider{
								localFSProvider: new(FakeProvider),
								puller:          fakePuller,
							}

							/* act */
							_, actualError := objectUnderTest.TryResolve(
								context.Background(),
								"dummyPkgRef",
							)

							/* assert */
							Expect(actualError).To(Equal(expectedErr))
						})
					})
					Context("puller.Pull doesn't error", func() {
						It("should return expected result", func() {
							/* arrange */
							providedPkgRef := "dummyPkgRef"
							basePath := "dummyBasePath"

							objectUnderTest := gitProvider{
								basePath:        basePath,
								localFSProvider: new(FakeProvider),
								puller:          new(fakePuller),
							}

							/* act */
							actualHandle, actualError := objectUnderTest.TryResolve(
								context.Background(),
								providedPkgRef,
							)

							/* assert */
							Expect(actualHandle).To(Equal(newGitHandle(filepath.Join(basePath, providedPkgRef), providedPkgRef)))
							Expect(actualError).To(BeNil())
						})
					})
				})
			})
		})
		Context("called in parallel w/ same pkg ref", func() {
			It("should not call localFSProvider.TryResolve & return same result", func() {
				/* arrange */
				providedPkgRef := "dummyPkgRef"
				expectedErr := errors.New("dummyError")

				fakeLocalFSProvider := new(FakeProvider)
				// err to trigger immediate return
				fakeLocalFSProvider.TryResolveStub = func(ctx context.Context, pkgRef string) (model.PkgHandle, error) {
					// ensure go routine has time to overlap
					<-time.After(100 * time.Millisecond)

					// error to trigger immediate return
					return nil, expectedErr
				}

				objectUnderTest := gitProvider{
					localFSProvider: fakeLocalFSProvider,
				}

				/* act */
				_, actualErr := objectUnderTest.TryResolve(
					context.Background(),
					providedPkgRef,
				)

				/* assert */
				Expect(fakeLocalFSProvider.TryResolveCallCount()).To(Equal(1))
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
		Context("called in parallel w/ different pkg ref", func() {
			It("should call localFSProvider.TryResolve w/ expected args", func() {
				/* arrange */
				providedCtx0 := context.Background()
				providedPkgRef0 := "dummyPkgRef0"

				providedCtx1 := context.Background()
				providedPkgRef1 := "dummyPkgRef1"

				fakeLocalFSProvider := new(FakeProvider)
				// err to trigger immediate return
				fakeLocalFSProvider.TryResolveReturns(nil, errors.New("dummyError"))

				objectUnderTest := gitProvider{
					localFSProvider: fakeLocalFSProvider,
				}

				/* act */
				objectUnderTest.TryResolve(
					providedCtx0,
					providedPkgRef0,
				)

				objectUnderTest.TryResolve(
					providedCtx1,
					providedPkgRef1,
				)

				/* assert */
				actualCtx0,
					actualPkgRef0 := fakeLocalFSProvider.TryResolveArgsForCall(0)
				Expect(actualCtx0).To(Equal(providedCtx0))
				Expect(actualPkgRef0).To(Equal(providedPkgRef0))

				actualCtx1,
					actualPkgRef1 := fakeLocalFSProvider.TryResolveArgsForCall(1)
				Expect(actualCtx1).To(Equal(providedCtx1))
				Expect(actualPkgRef1).To(Equal(providedPkgRef1))
			})
		})
	})
})
