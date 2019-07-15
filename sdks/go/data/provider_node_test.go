package data

import (
	"context"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/api/client"
)

var _ = Context("fsProvider", func() {
	Context("TryResolve", func() {
		It("should call nodeClient.ListDescendants w/ expected args", func() {
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
				actualReq := fakeNodeClient.ListDescendantsArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualReq).To(Equal(model.ListDescendantsReq{
				PkgRef:    providedDataRef,
				PullCreds: providedPullCreds,
			}))
		})
		Context("nodeClient.ListDirEntryd errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeNodeClient := new(client.Fake)

				listDirEntrysErr := errors.New("dummyError")
				fakeNodeClient.ListDescendantsReturns(nil, listDirEntrysErr)

				objectUnderTest := nodeProvider{
					nodeClient: fakeNodeClient,
				}

				/* act */
				_, actualErr := objectUnderTest.TryResolve(
					context.Background(),
					"dummyDataRef",
				)

				/* assert */
				Expect(actualErr).To(Equal(listDirEntrysErr))
			})
		})
		Context("nodeClient.ListDescendants doesn't err", func() {
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
