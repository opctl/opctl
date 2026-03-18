package git

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("_git", func() {
	Context("TryResolve", func() {
		Context("repo exists but completion marker doesn't", func() {
			Context("invalid git ref", func() {

			})
			It("should return error", func() {
				wd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				opRef := filepath.Join(wd, "../testdata/testop")

				objectUnderTest := New(filepath.Dir(opRef), nil)

				/* act */
				actualHandle, actualErr := objectUnderTest.TryResolve(
					context.Background(),
					opRef,
				)

				/* assert */
				Expect(actualErr).To(MatchError("invalid git ref: missing version"))
				Expect(actualHandle).To(BeNil())
			})
			Context("clone errors", func() {
				It("should return err", func() {
					dataDir, err := os.MkdirTemp("", "")
					if err != nil {
						panic(err)
					}
					objectUnderTest := New(dataDir, nil)

					/* act */
					_, actualErr := objectUnderTest.TryResolve(
						context.Background(),
						"not/exists",
					)

					/* assert */
					Expect(actualErr).To(MatchError("invalid git ref: missing version"))
				})
				Context("Clone doesn't error", func() {
					It("should return expected result", func() {
						/* arrange */
						// some public repo that's relatively small
						providedRef := "github.com/opspec-pkgs/_.op.create#3.3.1"
						basePath, err := os.MkdirTemp("", "")
						if err != nil {
							panic(err)
						}
						objectUnderTest := New(basePath, nil)
						expectedHandle := newHandle(filepath.Join(basePath, providedRef), providedRef)

						/* act */
						actualHandle, actualError := objectUnderTest.TryResolve(
							context.Background(),
							providedRef,
						)

						/* assert */
						Expect(actualHandle).To(Equal(expectedHandle))
						Expect(actualError).To(BeNil())
					})
				})
			})
		})
		Context("remote hash matches cached marker", func() {
			It("should return handle without re-cloning", func() {
				/* arrange */
				providedRef := "github.com/opspec-pkgs/_.op.create#3.3.1"
				basePath, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}
				objectUnderTest := New(basePath, nil)

				// first call — populates cache and writes the hash marker
				_, err = objectUnderTest.TryResolve(context.Background(), providedRef)
				Expect(err).To(BeNil())

				// record marker mtime before second call
				repoPath := filepath.Join(basePath, providedRef)
				markerPath := filepath.Join(basePath, fmt.Sprintf(".%x", sha1.Sum([]byte(repoPath))))
				markerInfoBefore, err := os.Stat(markerPath)
				Expect(err).To(BeNil())

				/* act — second call with same hash should hit cache */
				actualHandle, actualErr := objectUnderTest.TryResolve(context.Background(), providedRef)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualHandle).NotTo(BeNil())

				// marker must not have been rewritten (mtime unchanged = no re-clone)
				markerInfoAfter, err := os.Stat(markerPath)
				Expect(err).To(BeNil())
				Expect(markerInfoAfter.ModTime()).To(Equal(markerInfoBefore.ModTime()))
			})
		})
		Context("cached marker contains stale hash (tag re-pointed)", func() {
			It("should re-clone and update the marker", func() {
				/* arrange */
				providedRef := "github.com/opspec-pkgs/_.op.create#3.3.1"
				basePath, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}
				objectUnderTest := New(basePath, nil)

				// first call — populates cache
				_, err = objectUnderTest.TryResolve(context.Background(), providedRef)
				Expect(err).To(BeNil())

				// overwrite the marker with a fake hash to simulate a re-tagged commit
				repoPath := filepath.Join(basePath, providedRef)
				markerPath := filepath.Join(basePath, fmt.Sprintf(".%x", sha1.Sum([]byte(repoPath))))
				err = os.WriteFile(markerPath, []byte("0000000000000000000000000000000000000000"), 0755)
				Expect(err).To(BeNil())

				/* act — should detect hash mismatch and re-clone */
				actualHandle, actualErr := objectUnderTest.TryResolve(context.Background(), providedRef)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualHandle).NotTo(BeNil())

				// marker must now contain the real hash, not our fake one
				updatedHash, err := os.ReadFile(markerPath)
				Expect(err).To(BeNil())
				Expect(string(updatedHash)).NotTo(Equal("0000000000000000000000000000000000000000"))
				Expect(len(string(updatedHash))).To(Equal(40)) // SHA-1 hex length
			})
		})
		Context("called in parallel w/ same pkg ref", func() {
			It("should return expected result", func() {
				/* arrange */
				// some public repo that's relatively small
				providedRef := "github.com/opspec-pkgs/_.op.create#3.3.1"

				basePath, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}

				objectUnderTest := New(basePath, nil)

				expectedResult := newHandle(filepath.Join(basePath, providedRef), providedRef)

				var (
					actualResult1,
					actualResult2 model.DataHandle
				)
				var (
					actualErr1,
					actualErr2 error
				)

				/* act */
				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					actualResult1, actualErr1 = objectUnderTest.TryResolve(
						context.Background(),
						providedRef,
					)
					wg.Done()
				}()

				wg.Add(1)
				go func() {
					actualResult2, actualErr2 = objectUnderTest.TryResolve(
						context.Background(),
						providedRef,
					)
					wg.Done()
				}()

				wg.Wait()

				/* assert */
				Expect(actualErr1).To(BeNil())
				Expect(actualErr2).To(BeNil())
				Expect(actualResult1.Ref()).To(Equal(expectedResult.Ref()))
				Expect(actualResult2.Ref()).To(Equal(expectedResult.Ref()))
			})
		})
		Context("called in parallel w/ different pkg ref", func() {
			It("should return expected result", func() {
				/* arrange */
				// some public repo that's relatively small
				providedRef1 := "github.com/opspec-pkgs/_.op.create#3.3.1"
				providedRef2 := "github.com/opspec-pkgs/_.op.create#3.0.0"

				basePath, err := os.MkdirTemp("", "")
				if err != nil {
					panic(err)
				}

				objectUnderTest := New(basePath, nil)

				expectedResult1 := newHandle(filepath.Join(basePath, providedRef1), providedRef1)
				expectedResult2 := newHandle(filepath.Join(basePath, providedRef2), providedRef2)

				var (
					actualResult1,
					actualResult2 model.DataHandle
				)
				var (
					actualErr1,
					actualErr2 error
				)

				/* act */
				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					actualResult1, actualErr1 = objectUnderTest.TryResolve(
						context.Background(),
						providedRef1,
					)
					wg.Done()
				}()

				wg.Add(1)
				go func() {
					actualResult2, actualErr2 = objectUnderTest.TryResolve(
						context.Background(),
						providedRef2,
					)
					wg.Done()
				}()

				wg.Wait()

				/* assert */
				Expect(actualErr1).To(BeNil())
				Expect(actualResult1.Ref()).To(Equal(expectedResult1.Ref()))

				Expect(actualErr2).To(BeNil())
				Expect(actualResult2.Ref()).To(Equal(expectedResult2.Ref()))
			})
		})
	})
})
