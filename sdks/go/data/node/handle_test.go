package node

import (
	"context"
	"path"

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

			fakeCore := new(nodeFakes.FakeNode)

			objectUnderTest := handle{
				node:      fakeCore,
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
				DataRef:   path.Join(objectUnderTest.dataRef, providedContentPath),
				PullCreds: pullCreds,
			}))
		})
	})

	Context("ListDescendants", func() {
		It("should call core.ListDescendants w/ expected args", func() {
			/* arrange */
			providedCtx := context.TODO()

			dataRef := "dummyDataRef"
			pullCreds := &model.Creds{Username: "dummyUsername", Password: "dummyPassword"}

			fakeCore := new(nodeFakes.FakeNode)

			objectUnderTest := handle{
				node:      fakeCore,
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
				DataRef:   dataRef,
				PullCreds: pullCreds,
			}))
		})
	})

	Context("Ref", func() {
		It("should return expected ref", func() {
			/* arrange */
			dataRef := "dummyDataRef"

			fakeCore := new(nodeFakes.FakeNode)

			objectUnderTest := handle{
				node:    fakeCore,
				dataRef: dataRef,
			}

			/* act */
			actualRef := objectUnderTest.Ref()

			/* assert */
			Expect(actualRef).To(Equal(dataRef))
		})
	})
})
