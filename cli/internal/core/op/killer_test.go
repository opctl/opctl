package op

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
)

var _ = Context("Killer", func() {
	Context("Invoke", func() {
		It("Returns errors from node killing", func() {
			/* arrange */
			expectedError := errors.New("expected")

			fakeCore := new(nodeFakes.FakeOpNode)
			fakeCore.KillOpReturns(expectedError)

			objectUnderTest := newKiller(fakeCore)

			/* act */
			err := objectUnderTest.Kill(context.Background(), "opID")

			/* assert */
			Expect(err).To(MatchError(expectedError))
		})
		It("should call core.KillOp w/ expected args", func() {
			/* arrange */
			providedCtx := context.TODO()

			expectedCtx := providedCtx
			expectedReq := model.KillOpReq{
				OpID:       "dummyOpID",
				RootCallID: "dummyOpID",
			}

			fakeCore := new(nodeFakes.FakeOpNode)

			objectUnderTest := _killer{
				core: fakeCore,
			}

			/* act */
			err := objectUnderTest.Kill(expectedCtx, expectedReq.OpID)

			/* assert */
			actualCtx, actualReq := fakeCore.KillOpArgsForCall(0)
			Expect(err).To(BeNil())
			Expect(actualCtx).To(Equal(expectedCtx))
			Expect(actualReq).To(BeEquivalentTo(expectedReq))
		})
		Context("errors", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")
				fakeCore := new(nodeFakes.FakeOpNode)
				fakeCore.KillOpReturns(expectedError)

				objectUnderTest := _killer{
					core: fakeCore,
				}

				/* act */
				err := objectUnderTest.Kill(context.TODO(), "")

				/* assert */
				Expect(err).To(MatchError(expectedError))
			})
		})
	})
})
