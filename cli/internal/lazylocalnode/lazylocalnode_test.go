package lazylocalnode

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	nodeFakes "github.com/opctl/opctl/cli/internal/nodeprovider/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	apiClientFakes "github.com/opctl/opctl/sdks/go/node/api/client/fakes"
)

var _ = Context("lazylocalnode", func() {
	It("conforms to the OpNode interfaces", func() {
		objectUnderTest := New(new(nodeFakes.FakeNodeProvider))
		test := func(n node.OpNode) {
			Expect(n).NotTo(BeNil())
		}
		test(objectUnderTest)
	})
	Context("wraps an APIClient, first ensuring the remote node is initialized", func() {
		It("for AddAuth", func() {
			// arrange
			fakeNodeProvider := new(nodeFakes.FakeNodeProvider)
			fakeNodeHandle := new(nodeFakes.FakeNodeHandle)
			fakeAPIClient := new(apiClientFakes.FakeAPIClient)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)
			objectUnderTest := New(fakeNodeProvider)
			arg1, arg2 := context.Background(), model.AddAuthReq{}

			// act
			objectUnderTest.AddAuth(arg1, arg2)

			// assert
			Expect(fakeNodeProvider.CreateNodeIfNotExistsCallCount()).To(Equal(1))
			Expect(fakeNodeHandle.APIClientCallCount()).To(Equal(1))
			aArg1, aArg2 := fakeAPIClient.AddAuthArgsForCall(0)
			Expect(aArg1).To(Equal(arg1))
			Expect(aArg2).To(Equal(arg2))
		})
		It("for GetEventStream", func() {
			// arrange
			fakeNodeProvider := new(nodeFakes.FakeNodeProvider)
			fakeNodeHandle := new(nodeFakes.FakeNodeHandle)
			fakeAPIClient := new(apiClientFakes.FakeAPIClient)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)
			objectUnderTest := New(fakeNodeProvider)
			arg1, arg2 := context.Background(), &model.GetEventStreamReq{}

			// act
			objectUnderTest.GetEventStream(arg1, arg2)

			// assert
			Expect(fakeNodeProvider.CreateNodeIfNotExistsCallCount()).To(Equal(1))
			Expect(fakeNodeHandle.APIClientCallCount()).To(Equal(1))
			aArg1, aArg2 := fakeAPIClient.GetEventStreamArgsForCall(0)
			Expect(aArg1).To(Equal(arg1))
			Expect(aArg2).To(Equal(arg2))
		})
		It("for KillOp", func() {
			// arrange
			fakeNodeProvider := new(nodeFakes.FakeNodeProvider)
			fakeNodeHandle := new(nodeFakes.FakeNodeHandle)
			fakeAPIClient := new(apiClientFakes.FakeAPIClient)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)
			objectUnderTest := New(fakeNodeProvider)
			arg1, arg2 := context.Background(), model.KillOpReq{}

			// act
			objectUnderTest.KillOp(arg1, arg2)

			// assert
			Expect(fakeNodeProvider.CreateNodeIfNotExistsCallCount()).To(Equal(1))
			Expect(fakeNodeHandle.APIClientCallCount()).To(Equal(1))
			aArg1, aArg2 := fakeAPIClient.KillOpArgsForCall(0)
			Expect(aArg1).To(Equal(arg1))
			Expect(aArg2).To(Equal(arg2))
		})
		It("for StartOp", func() {
			// arrange
			fakeNodeProvider := new(nodeFakes.FakeNodeProvider)
			fakeNodeHandle := new(nodeFakes.FakeNodeHandle)
			fakeAPIClient := new(apiClientFakes.FakeAPIClient)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)
			objectUnderTest := New(fakeNodeProvider)
			arg1, arg2 := context.Background(), model.StartOpReq{}

			// act
			objectUnderTest.StartOp(arg1, arg2)

			// assert
			Expect(fakeNodeProvider.CreateNodeIfNotExistsCallCount()).To(Equal(1))
			Expect(fakeNodeHandle.APIClientCallCount()).To(Equal(1))
			aArg1, aArg2 := fakeAPIClient.StartOpArgsForCall(0)
			Expect(aArg1).To(Equal(arg1))
			Expect(aArg2).To(Equal(arg2))
		})
		It("for GetData", func() {
			// arrange
			fakeNodeProvider := new(nodeFakes.FakeNodeProvider)
			fakeNodeHandle := new(nodeFakes.FakeNodeHandle)
			fakeAPIClient := new(apiClientFakes.FakeAPIClient)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)
			objectUnderTest := New(fakeNodeProvider)
			arg1, arg2 := context.Background(), model.GetDataReq{}

			// act
			objectUnderTest.GetData(arg1, arg2)

			// assert
			Expect(fakeNodeProvider.CreateNodeIfNotExistsCallCount()).To(Equal(1))
			Expect(fakeNodeHandle.APIClientCallCount()).To(Equal(1))
			aArg1, aArg2 := fakeAPIClient.GetDataArgsForCall(0)
			Expect(aArg1).To(Equal(arg1))
			Expect(aArg2).To(Equal(arg2))

		})
		It("for ListDescendants", func() {
			// arrange
			fakeNodeProvider := new(nodeFakes.FakeNodeProvider)
			fakeNodeHandle := new(nodeFakes.FakeNodeHandle)
			fakeAPIClient := new(apiClientFakes.FakeAPIClient)
			fakeNodeHandle.APIClientReturns(fakeAPIClient)
			fakeNodeProvider.CreateNodeIfNotExistsReturns(fakeNodeHandle, nil)
			objectUnderTest := New(fakeNodeProvider)
			arg1, arg2 := context.Background(), model.ListDescendantsReq{}

			// act
			objectUnderTest.ListDescendants(arg1, arg2)

			// assert
			Expect(fakeNodeProvider.CreateNodeIfNotExistsCallCount()).To(Equal(1))
			Expect(fakeNodeHandle.APIClientCallCount()).To(Equal(1))
			aArg1, aArg2 := fakeAPIClient.ListDescendantsArgsForCall(0)
			Expect(aArg1).To(Equal(arg1))
			Expect(aArg2).To(Equal(arg2))
		})
	})
	Context("passes through errors", func() {
		// arrange
		expectedErr := errors.New("expected")
		fakeNodeProvider := new(nodeFakes.FakeNodeProvider)
		fakeNodeProvider.CreateNodeIfNotExistsReturns(nil, expectedErr)
		objectUnderTest := New(fakeNodeProvider)

		It("for AddAuth", func() {
			// act
			err := objectUnderTest.AddAuth(context.Background(), model.AddAuthReq{})
			// assert
			Expect(err).To(MatchError(expectedErr))
		})
		It("for GetEventStream", func() {
			// act
			_, err := objectUnderTest.GetEventStream(context.Background(), &model.GetEventStreamReq{})
			// assert
			Expect(err).To(MatchError(expectedErr))
		})
		It("for KillOp", func() {
			// act
			err := objectUnderTest.KillOp(context.Background(), model.KillOpReq{})
			// assert
			Expect(err).To(MatchError(expectedErr))
		})
		It("for StartOp", func() {
			// act
			_, err := objectUnderTest.StartOp(context.Background(), model.StartOpReq{})
			// assert
			Expect(err).To(MatchError(expectedErr))
		})
		It("for GetData", func() {
			// act
			_, err := objectUnderTest.GetData(context.Background(), model.GetDataReq{})
			// assert
			Expect(err).To(MatchError(expectedErr))
		})
		It("for ListDescendants", func() {
			// act
			_, err := objectUnderTest.ListDescendants(context.Background(), model.ListDescendantsReq{})
			// assert
			Expect(err).To(MatchError(expectedErr))
		})
	})
})
