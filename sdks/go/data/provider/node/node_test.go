package node

import (
	"context"

	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	clientFakes "github.com/opctl/opctl/sdks/go/node/api/client/fakes"
)

var _ = Context("_node", func() {
	Context("TryResolve", func() {
		It("should call apiClient.ListDescendants w/ expected args", func() {
			/* arrange */
			providedDataRef := "dummyDataRef"

			fakeAPIClient := new(clientFakes.FakeClient)

			providedPullCreds := &model.PullCreds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			objectUnderTest := _node{
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
				fakeAPIClient := new(clientFakes.FakeClient)

				listDirEntrysErr := errors.New("dummyError")
				fakeAPIClient.ListDescendantsReturns(nil, listDirEntrysErr)

				objectUnderTest := _node{
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

				fakeAPIClient := new(clientFakes.FakeClient)

				providedPullCreds := &model.PullCreds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				}

				objectUnderTest := _node{
					apiClient: fakeAPIClient,
					pullCreds: providedPullCreds,
				}

				/* act */
				actualHandle, actualErr := objectUnderTest.TryResolve(
					context.Background(),
					providedDataRef,
				)

				/* assert */
				Expect(actualHandle).To(Equal(newHandle(fakeAPIClient, providedDataRef, providedPullCreds)))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
