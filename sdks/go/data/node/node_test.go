package node

import (
	"context"

	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

var _ = Context("_node", func() {
	Context("TryResolve", func() {
		It("should call apiClient.ListDescendants w/ expected args", func() {
			/* arrange */
			providedDataRef := "dummyDataRef"

			fakeCore := new(nodeFakes.FakeNode)

			providedPullCreds := &model.Creds{
				Username: "dummyUsername",
				Password: "dummyPassword",
			}

			objectUnderTest := _node{
				node:      fakeCore,
				pullCreds: providedPullCreds,
			}

			/* act */
			objectUnderTest.TryResolve(
				context.Background(),
				providedDataRef,
			)

			/* assert */
			actualContext,
				actualReq := fakeCore.ListDescendantsArgsForCall(0)

			Expect(actualContext).To(Equal(context.TODO()))
			Expect(actualReq).To(Equal(model.ListDescendantsReq{
				PkgRef:    providedDataRef,
				PullCreds: providedPullCreds,
			}))
		})
		Context("core.ListDescendants errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeCore := new(nodeFakes.FakeNode)

				listDirEntrysErr := errors.New("dummyError")
				fakeCore.ListDescendantsReturns(nil, listDirEntrysErr)

				objectUnderTest := _node{
					node: fakeCore,
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
		Context("core.ListDescendants doesn't err", func() {
			It("should return expected result", func() {
				/* arrange */
				providedDataRef := "dummyDataRef"

				fakeCore := new(nodeFakes.FakeNode)

				providedPullCreds := &model.Creds{
					Username: "dummyUsername",
					Password: "dummyPassword",
				}

				objectUnderTest := _node{
					node:      fakeCore,
					pullCreds: providedPullCreds,
				}

				/* act */
				actualHandle, actualErr := objectUnderTest.TryResolve(
					context.Background(),
					providedDataRef,
				)

				/* assert */
				Expect(actualHandle).To(Equal(newHandle(fakeCore, providedDataRef, providedPullCreds)))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
