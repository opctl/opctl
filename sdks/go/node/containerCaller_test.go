package node

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	. "github.com/opctl/opctl/sdks/go/node/pubsub/fakes"
)

var _ = Context("containerCaller", func() {
	dbDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	db, err := badger.Open(
		badger.DefaultOptions(dbDir).WithLogger(nil),
	)
	if err != nil {
		panic(err)
	}

	Context("newContainerCaller", func() {
		It("should return containerCaller", func() {
			/* arrange/act/assert */
			Expect(newContainerCaller(
				new(FakeContainerRuntime),
				new(FakePubSub),
				newStateStore(context.Background(), db, new(FakePubSub)),
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		It("should call containerRuntime.RunContainer w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedContainerCall := &model.ContainerCall{
				BaseCall: model.BaseCall{},
				Image:    &model.ContainerCallImage{},
			}
			providedRootCallID := "providedRootCallID"
			fakeContainerRuntime := new(FakeContainerRuntime)

			fakeContainerRuntime.RunContainerStub = func(
				ctx context.Context,
				req *model.ContainerCall,
				rootCallID string,
				eventPublisher pubsub.EventPublisher,
				stdOut io.WriteCloser,
				stdErr io.WriteCloser,
			) (*int64, error) {

				stdErr.Close()
				stdOut.Close()

				return nil, nil
			}

			fakePubSub := new(FakePubSub)

			objectUnderTest := _containerCaller{
				containerRuntime: fakeContainerRuntime,
				pubSub:           fakePubSub,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				providedContainerCall,
				map[string]*ipld.Node{},
				&model.ContainerCallSpec{},
				providedRootCallID,
			)

			/* assert */
			_,
				actualContainerCall,
				actualRootCallID,
				actualEventPublisher,
				_,
				_ := fakeContainerRuntime.RunContainerArgsForCall(0)
			Expect(actualContainerCall).To(Equal(providedContainerCall))
			Expect(actualRootCallID).To(Equal(providedRootCallID))
			Expect(actualEventPublisher).To(Equal(fakePubSub))
		})
		Context("containerRuntime.RunContainer errors", func() {
			It("should publish expected ContainerExited", func() {
				/* arrange */
				expectedErrorMessage := "expectedErrorMessage"
				fakePubSub := new(FakePubSub)

				fakeContainerRuntime := new(FakeContainerRuntime)

				fakeContainerRuntime.RunContainerStub = func(
					ctx context.Context,
					req *model.ContainerCall,
					rootCallID string,
					eventPublisher pubsub.EventPublisher,
					stdOut io.WriteCloser,
					stdErr io.WriteCloser,
				) (*int64, error) {

					stdErr.Close()
					stdOut.Close()

					return nil, errors.New(expectedErrorMessage)
				}

				objectUnderTest := _containerCaller{
					containerRuntime: fakeContainerRuntime,
					pubSub:           fakePubSub,
				}

				/* act */
				actualOutputs, actualErr := objectUnderTest.Call(
					context.Background(),
					&model.ContainerCall{
						BaseCall: model.BaseCall{},
						Image:    &model.ContainerCallImage{},
					},
					map[string]*ipld.Node{},
					&model.ContainerCallSpec{},
					"rootCallID",
				)

				/* assert */
				Expect(actualOutputs).To(Equal(map[string]*ipld.Node{}))
				Expect(actualErr).To(MatchError(expectedErrorMessage))
			})
		})
	})

	It("should return expected results", func() {
		/* arrange */
		providedOpPath := "providedOpPath"
		providedContainerCall := &model.ContainerCall{
			BaseCall: model.BaseCall{
				OpPath: providedOpPath,
			},
			ContainerID: "providedContainerID",
			Image:       &model.ContainerCallImage{},
		}
		providedInboundScope := map[string]*ipld.Node{}
		providedContainerCallSpec := &model.ContainerCallSpec{}

		fakeContainerRuntime := new(FakeContainerRuntime)

		expectedErr := errors.New("io: read/write on closed pipe")
		fakeContainerRuntime.RunContainerStub = func(
			ctx context.Context,
			req *model.ContainerCall,
			rootCallID string,
			eventPublisher pubsub.EventPublisher,
			stdOut io.WriteCloser,
			stdErr io.WriteCloser,
		) (*int64, error) {

			stdErr.Close()
			stdOut.Close()

			return nil, expectedErr
		}

		objectUnderTest := _containerCaller{
			containerRuntime: fakeContainerRuntime,
			pubSub:           new(FakePubSub),
		}

		/* act */
		actualOutputs, actualErr := objectUnderTest.Call(
			context.Background(),
			providedContainerCall,
			providedInboundScope,
			providedContainerCallSpec,
			"rootCallID",
		)

		/* assert */
		Expect(actualOutputs).To(Equal(map[string]*ipld.Node{}))
		Expect(actualErr).To(Equal(expectedErr))
	})
})
