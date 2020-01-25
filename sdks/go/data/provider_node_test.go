package data

import (
	"context"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

var _ = Context("nodeProvider", func() {
	Context("TryResolve", func() {
		It("should call apiClient.ListDescendants w/ expected args", func() {
			/* arrange */
			providedDataRef := "dummyDataRef"

			fakeAPIClient := new(client.Fake)

			providedPullCreds := &model.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			objectUnderTest := nodeProvider{
				apiClient: fakeAPIClient,
				pullCreds: providedPullCreds,
			}

			/* act */
			objectUnderTest.TryResolve(
				context.Background(),
				providedDataRef,
			)

			/* assert */
			actualContext,
				actualReq := fakeAPIClient.ListDescendantsArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualReq).To(Equal(model.ListDescendantsReq{
				PkgRef:    providedDataRef,
				PullCreds: providedPullCreds,
			}))
		})
		Context("apiClient.ListDirEntryd errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeAPIClient := new(client.Fake)

				listDirEntrysErr := errors.New("dummyError")
				fakeAPIClient.ListDescendantsReturns(nil, listDirEntrysErr)

				objectUnderTest := nodeProvider{
					apiClient: fakeAPIClient,
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
		Context("apiClient.ListDescendants doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedDataRef := "dummyDataRef"

				fakeAPIClient := new(client.Fake)

				providedPullCreds := &model.PullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				}

				fakePuller := new(fakePuller)

				objectUnderTest := nodeProvider{
					apiClient: fakeAPIClient,
					pullCreds: providedPullCreds,
					puller:    fakePuller,
				}

				/* act */
				actualHandle, actualErr := objectUnderTest.TryResolve(
					context.Background(),
					providedDataRef,
				)

				/* assert */
				Expect(actualHandle).To(Equal(newNodeHandle(fakeAPIClient, providedDataRef, providedPullCreds)))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
