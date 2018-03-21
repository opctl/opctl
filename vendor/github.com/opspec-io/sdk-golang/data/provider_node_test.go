package data

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/node/api/client"
	"github.com/pkg/errors"
)

var _ = Context("fsProvider", func() {
	Context("TryResolve", func() {
		It("should call nodeClient.ListPkgContents w/ expected args", func() {
			/* arrange */
			providedDataRef := "dummyDataRef"

			fakeNodeClient := new(client.Fake)

			providedPullCreds := &model.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			objectUnderTest := nodeProvider{
				nodeClient: fakeNodeClient,
				pullCreds:  providedPullCreds,
			}

			/* act */
			objectUnderTest.TryResolve(
				context.Background(),
				providedDataRef,
			)

			/* assert */
			actualContext,
				actualReq := fakeNodeClient.ListPkgContentsArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualReq).To(Equal(model.ListPkgContentsReq{
				PkgRef:    providedDataRef,
				PullCreds: providedPullCreds,
			}))
		})
		Context("nodeClient.ListPkgContentd errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeNodeClient := new(client.Fake)

				listPkgContentsErr := errors.New("dummyError")
				fakeNodeClient.ListPkgContentsReturns(nil, listPkgContentsErr)

				objectUnderTest := nodeProvider{
					nodeClient: fakeNodeClient,
				}

				/* act */
				_, actualErr := objectUnderTest.TryResolve(
					context.Background(),
					"dummyDataRef",
				)

				/* assert */
				Expect(actualErr).To(Equal(listPkgContentsErr))
			})
		})
		Context("nodeClient.ListPkgContents doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedDataRef := "dummyDataRef"

				fakeNodeClient := new(client.Fake)

				providedPullCreds := &model.PullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				}

				fakePuller := new(fakePuller)

				objectUnderTest := nodeProvider{
					nodeClient: fakeNodeClient,
					pullCreds:  providedPullCreds,
					puller:     fakePuller,
				}

				/* act */
				actualHandle, actualErr := objectUnderTest.TryResolve(
					context.Background(),
					providedDataRef,
				)

				/* assert */
				Expect(actualHandle).To(Equal(newNodeHandle(fakeNodeClient, providedDataRef, providedPullCreds)))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
