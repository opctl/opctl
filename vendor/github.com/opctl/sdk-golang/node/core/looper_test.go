package core

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

var _ = Context("looper", func() {
	Context("newLooper", func() {
		It("should return looper", func() {
			/* arrange/act/assert */
			Expect(newLooper(
				new(fakeCaller),
				new(pubsub.Fake),
			)).To(Not(BeNil()))
		})
	})

	Context("Loop", func() {
		Context("initial dcgLoop.Until true", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeLoopInterpreter := new(loop.FakeInterpreter)
				until := true
				fakeLoopInterpreter.InterpretReturns(&model.DCGLoop{Until: &until}, nil)

				fakeCaller := new(fakeCaller)

				objectUnderTest := _looper{
					caller:              fakeCaller,
					loopInterpreter:     fakeLoopInterpreter,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.Loop(
					"id",
					map[string]*model.Value{},
					&model.SCG{Loop: &model.SCGLoop{}},
					new(data.FakeHandle),
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCallCount()).To(Equal(0))
			})
		})
		Context("initial dcgLoop.For.Each empty", func() {
			It("should not call caller.Call", func() {
				/* arrange */
				fakeLoopInterpreter := new(loop.FakeInterpreter)
				fakeLoopInterpreter.InterpretReturns(
					&model.DCGLoop{
						For: &model.DCGLoopFor{
							Each: &model.Value{
								Array: []interface{}{},
							},
						},
					},
					nil,
				)

				fakeCaller := new(fakeCaller)

				objectUnderTest := _looper{
					caller:              fakeCaller,
					loopInterpreter:     fakeLoopInterpreter,
					uniqueStringFactory: new(uniquestring.Fake),
				}

				/* act */
				objectUnderTest.Loop(
					"id",
					map[string]*model.Value{},
					&model.SCG{Loop: &model.SCGLoop{}},
					new(data.FakeHandle),
					"rootOpID",
				)

				/* assert */
				Expect(fakeCaller.CallCallCount()).To(Equal(0))
			})
		})
		Context("initial dcgLoop.Until false", func() {
			It("should call caller.Call w/ expected args", func() {
				/* arrange */
				providedScope := map[string]*model.Value{}
				index := "index"
				providedSCG := &model.SCG{
					Loop: &model.SCGLoop{
						Index: &index,
					},
				}
				providedOpHandle := new(data.FakeHandle)
				providedRootOpID := "providedRootOpID"

				fakeLoopInterpreter := new(loop.FakeInterpreter)
				until := false
				fakeLoopInterpreter.InterpretReturns(
					&model.DCGLoop{
						Until: &until,
						Index: providedSCG.Loop.Index,
					},
					nil,
				)

				zero := float64(0)
				expectedScope := map[string]*model.Value{
					index: &model.Value{Number: &zero},
				}

				fakeCaller := new(fakeCaller)
				// error to trigger immediate return
				fakeCaller.CallReturns(errors.New("dummyError"))

				expectedID := "expectedID"
				fakeUniqueStringFactory := new(uniquestring.Fake)
				fakeUniqueStringFactory.ConstructReturns(expectedID, nil)

				objectUnderTest := _looper{
					caller:              fakeCaller,
					loopInterpreter:     fakeLoopInterpreter,
					uniqueStringFactory: fakeUniqueStringFactory,
				}

				/* act */
				objectUnderTest.Loop(
					"id",
					providedScope,
					providedSCG,
					providedOpHandle,
					providedRootOpID,
				)

				/* assert */
				actualCallID,
					actualScope,
					actualSCG,
					actualOpHandle,
					actualRootOpID := fakeCaller.CallArgsForCall(0)

				Expect(actualCallID).To(Equal(expectedID))
				Expect(actualScope).To(Equal(expectedScope))
				Expect(actualSCG).To(Equal(providedSCG))
				Expect(actualOpHandle).To(Equal(providedOpHandle))
				Expect(actualRootOpID).To(Equal(providedRootOpID))
			})
		})
	})
})
