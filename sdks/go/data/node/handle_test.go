package node

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

var _ = Context("handle", func() {

	Context("GetContent", func() {

		It("should call core.GetData w/ expected args", func() {
			/* arrange */
			providedCtx := context.TODO()
			providedContentPath := "dummyContentPath"

			dataRef := "dummyDataRef"
			pullCreds := &model.Creds{Username: "dummyUsername", Password: "dummyPassword"}

			fakeCore := new(nodeFakes.FakeOpNode)

			objectUnderTest := handle{
				opNode:    fakeCore,
				dataRef:   dataRef,
				pullCreds: pullCreds,
			}

			/* act */
			objectUnderTest.GetContent(providedCtx, providedContentPath)

			/* assert */
			actualCtx,
				actualReq := fakeCore.GetDataArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualReq).To(Equal(model.GetDataReq{
				ContentPath: providedContentPath,
				PkgRef:      dataRef,
				PullCreds:   pullCreds,
			}))
		})
	})

	Context("ListDescendants", func() {
		It("should call core.ListDescendants w/ expected args", func() {
			/* arrange */
			providedCtx := context.TODO()

			dataRef := "dummyDataRef"
			pullCreds := &model.Creds{Username: "dummyUsername", Password: "dummyPassword"}

			fakeCore := new(nodeFakes.FakeOpNode)

			objectUnderTest := handle{
				opNode:    fakeCore,
				dataRef:   dataRef,
				pullCreds: pullCreds,
			}

			/* act */
			objectUnderTest.ListDescendants(providedCtx)

			/* assert */
			actualCtx,
				actualReq := fakeCore.ListDescendantsArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualReq).To(Equal(model.ListDescendantsReq{
				PkgRef:    dataRef,
				PullCreds: pullCreds,
			}))
		})
	})

	Context("Ref", func() {
		It("should return expected ref", func() {
			/* arrange */
			dataRef := "dummyDataRef"

			fakeCore := new(nodeFakes.FakeOpNode)

			objectUnderTest := handle{
				opNode:  fakeCore,
				dataRef: dataRef,
			}

			/* act */
			actualRef := objectUnderTest.Ref()

			/* assert */
			Expect(actualRef).To(Equal(dataRef))
		})
	})
})
