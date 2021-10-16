package core

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	containerRuntimeFakes "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
	. "github.com/opctl/opctl/sdks/go/node/core/internal/fakes"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

var _ = Context("parallelCaller", func() {
	Context("newParallelCaller", func() {
		It("should return parallelCaller", func() {
			/* arrange/act/assert */
			Expect(newParallelCaller(new(FakeCaller))).To(Not(BeNil()))
		})
	})

	Context("Call", func() {
		Context("caller errors", func() {
			It("should return expected results", func() {
				/* arrange */
				dbDir, err := ioutil.TempDir("", "")
				if err != nil {
					panic(err)
				}

				db, err := badger.Open(
					badger.DefaultOptions(dbDir).WithLogger(nil),
				)
				if err != nil {
					panic(err)
				}
				pubSub := pubsub.New(db)

				objectUnderTest := _parallelCaller{
					caller: newCaller(
						newContainerCaller(
							new(containerRuntimeFakes.FakeContainerRuntime),
							pubSub,
							newStateStore(
								context.Background(),
								db,
								pubSub,
							),
						),
						dbDir,
						pubSub,
					),
				}

				/* act */
				_, actualErr := objectUnderTest.Call(
					context.Background(),
					"callID",
					map[string]*model.Value{},
					"rootCallID",
					"opPath",
					[]*model.CallSpec{
						{
							// intentionally invalid, will produce validation error
							Container: &model.ContainerCallSpec{},
						},
					},
				)

				/* assert */
				Expect(actualErr).To(MatchError("image required"))
			})
		})

		It("should start each child as expected", func() {
			/* arrange */
			dbDir, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			db, err := badger.Open(
				badger.DefaultOptions(dbDir).WithLogger(nil),
			)
			if err != nil {
				panic(err)
			}
			pubSub := pubsub.New(db)

			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			providedOpRef := "providedOpRef"
			providedParentID := "providedParentID"
			providedRootID := "providedRootID"
			childOpRef := filepath.Join(wd, "testdata/parallelCaller")
			input1Key := "input1"
			childOp1Path := filepath.Join(childOpRef, "op1")
			childOp2Path := filepath.Join(childOpRef, "op2")

			ctx := context.Background()

			fakeContainerRuntime := new(containerRuntimeFakes.FakeContainerRuntime)
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

			eventChannel, err := pubSub.Subscribe(
				ctx,
				model.EventFilter{},
			)
			if err != nil {
				panic(err)
			}

			input1Value := "input1Value"
			providedInboundScope := map[string]*model.Value{
				input1Key: {String: &input1Value},
			}

			objectUnderTest := _parallelCaller{
				caller: newCaller(
					newContainerCaller(
						fakeContainerRuntime,
						pubSub,
						newStateStore(
							ctx,
							db,
							pubSub,
						),
					),
					dbDir,
					pubSub,
				),
			}

			/* act */
			_, actualErr := objectUnderTest.Call(
				ctx,
				providedParentID,
				providedInboundScope,
				providedRootID,
				providedOpRef,
				[]*model.CallSpec{
					{
						Op: &model.OpCallSpec{
							Ref: childOp1Path,
							Inputs: map[string]interface{}{
								input1Key: nil,
							},
						},
					},
					{
						Op: &model.OpCallSpec{
							Ref: childOp2Path,
							Inputs: map[string]interface{}{
								input1Key: nil,
							},
						},
					},
				},
			)

			/* assert */
			Expect(actualErr).To(BeNil())

			actualChildCalls := []model.CallStarted{}
			go func() {
				for event := range eventChannel {
					if event.CallStarted != nil && event.CallStarted.Call.Op != nil {
						// ignore props we can't readily assert
						event.CallStarted.Call.Op.ChildCallCallSpec = nil
						event.CallStarted.Call.Op.ChildCallID = ""
						event.CallStarted.Call.Op.OpID = ""
						event.CallStarted.Call.ID = ""

						actualChildCalls = append(actualChildCalls, *event.CallStarted)
					}
				}
			}()

			Eventually(
				func() []model.CallStarted { return actualChildCalls },
			).Should(
				ContainElements(
					[]model.CallStarted{
						{
							Call: model.Call{
								Op: &model.OpCall{
									BaseCall: model.BaseCall{
										OpPath: childOp1Path,
									},
									Inputs:            providedInboundScope,
									ChildCallCallSpec: nil,
								},
								ParentID: &providedParentID,
								RootID:   providedRootID,
							},
							Ref: providedOpRef,
						},
						{
							Call: model.Call{
								Op: &model.OpCall{
									BaseCall: model.BaseCall{
										OpPath: childOp2Path,
									},
									Inputs:            providedInboundScope,
									ChildCallCallSpec: nil,
								},
								ParentID: &providedParentID,
								RootID:   providedRootID,
							},
							Ref: providedOpRef,
						},
					},
				),
			)
		})

		It("adds outputs to scope", func() {
			/* arrange */
			providedOpRef := "providedOpRef"
			providedParentID := "providedParentID"
			providedRootID := "providedRootID"
			input1Value := "input1Value"
			providedInboundScope := map[string]*model.Value{
				"input": {String: &input1Value},
			}

			fakeCaller := new(FakeCaller)

			expectedOutput0 := "outputVal0"
			expectedOutputs0 := map[string]*model.Value{
				"output0": {String: &expectedOutput0},
			}
			fakeCaller.CallReturnsOnCall(0, expectedOutputs0, nil)
			expectedOutput1 := "outputVal1"
			expectedOutputs1 := map[string]*model.Value{
				"output1": {String: &expectedOutput1},
			}
			fakeCaller.CallReturnsOnCall(1, expectedOutputs1, nil)

			objectUnderTest := _parallelCaller{
				caller: fakeCaller,
			}

			/* act */
			actualOutputs, actualErr := objectUnderTest.Call(
				context.Background(),
				providedParentID,
				providedInboundScope,
				providedRootID,
				providedOpRef,
				[]*model.CallSpec{{}, {}},
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(actualOutputs).To(Equal(map[string]*model.Value{
				"input":   {String: &input1Value},
				"output0": {String: &expectedOutput0},
				"output1": {String: &expectedOutput1},
			}))
		})

		It("cancels other children when one fails", func() {
			/* arrange */
			providedOpRef := "providedOpRef"
			providedParentID := "providedParentID"
			providedRootID := "providedRootID"

			fakeCaller := new(FakeCaller)

			expectedError := errors.New("fail")
			fakeCaller.CallReturnsOnCall(0, nil, expectedError)

			objectUnderTest := _parallelCaller{
				caller: fakeCaller,
			}

			/* act */
			actualOutputs, actualErr := objectUnderTest.Call(
				context.Background(),
				providedParentID,
				map[string]*model.Value{},
				providedRootID,
				providedOpRef,
				[]*model.CallSpec{{}, {}},
			)

			/* assert */
			Expect(actualErr).To(MatchError(expectedError))
			Expect(actualOutputs).To(BeNil())
		})

		It("cancels early when context is cancelled", func() {
			/* arrange */
			providedOpRef := "providedOpRef"
			providedParentID := "providedParentID"
			providedRootID := "providedRootID"

			objectUnderTest := _parallelCaller{
				caller: new(FakeCaller),
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			/* act */
			actualOutputs, actualErr := objectUnderTest.Call(
				ctx,
				providedParentID,
				map[string]*model.Value{},
				providedRootID,
				providedOpRef,
				[]*model.CallSpec{{}, {}},
			)

			/* assert */
			Expect(actualErr).To(MatchError(context.Canceled))
			Expect(actualOutputs).To(BeNil())
		})
	})

	Context("needs", func() {
		It("cancels needed calls once all dependents have finished", func() {
			/* arrange */
			providedOpRef := "providedOpRef"
			providedParentID := "providedParentID"
			providedRootID := "providedRootID"

			neededName := "needed"

			fakeCaller := new(FakeCaller)
			fakeCaller.CallStub = func(
				ctx context.Context,
				id string,
				scope map[string]*model.Value,
				callSpec *model.CallSpec,
				opPath string,
				parentCallID *string,
				rootCallID string,
			) (map[string]*model.Value, error) {
				if callSpec.Name != nil && *callSpec.Name == neededName {
					// this error will be ignored, since this is a needed call
					return nil, errors.New("done")
				} else {
					return map[string]*model.Value{}, nil
				}
			}

			objectUnderTest := _parallelCaller{
				caller: fakeCaller,
			}

			/* act */
			_, actualErr := objectUnderTest.Call(
				context.Background(),
				providedParentID,
				map[string]*model.Value{},
				providedRootID,
				providedOpRef,
				[]*model.CallSpec{
					{Name: &neededName},
					{Needs: []string{neededName}},
				},
			)

			/* assert */
			Expect(actualErr).To(BeNil())
		})
	})
})
