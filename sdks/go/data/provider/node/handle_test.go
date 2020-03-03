package node

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	clientFakes "github.com/opctl/opctl/sdks/go/node/api/client/fakes"
)

var _ = Context("handle", func() {

	Context("GetContent", func() {

		It("should call client.GetData w/ expected args", func() {
			/* arrange */
			providedCtx := context.TODO()
			providedContentPath := "dummyContentPath"

			dataRef := "dummyDataRef"
			pullCreds := &model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"}

			fakeClient := new(clientFakes.FakeClient)

			objectUnderTest := handle{
				client:    fakeClient,
				dataRef:   dataRef,
				pullCreds: pullCreds,
			}

			/* act */
			objectUnderTest.GetContent(providedCtx, providedContentPath)

			/* assert */
			actualCtx,
				actualReq := fakeClient.GetDataArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualReq).To(Equal(model.GetDataReq{
				ContentPath: providedContentPath,
				PkgRef:      dataRef,
				PullCreds:   pullCreds,
			}))
		})
	})

	Context("ListDescendants", func() {
		It("should call client.ListDescendants w/ expected args", func() {
			/* arrange */
			providedCtx := context.TODO()

			dataRef := "dummyDataRef"
			pullCreds := &model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"}

			fakeClient := new(clientFakes.FakeClient)

			objectUnderTest := handle{
				client:    fakeClient,
				dataRef:   dataRef,
				pullCreds: pullCreds,
			}

			/* act */
			objectUnderTest.ListDescendants(providedCtx)

			/* assert */
			actualCtx,
				actualReq := fakeClient.ListDescendantsArgsForCall(0)

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

			fakeClient := new(clientFakes.FakeClient)

			objectUnderTest := handle{
				client:  fakeClient,
				dataRef: dataRef,
			}

			/* act */
			actualRef := objectUnderTest.Ref()

			/* assert */
			Expect(actualRef).To(Equal(dataRef))
		})
	})
})
