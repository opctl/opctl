package loop

import (
	"errors"

	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/predicates"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/data"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/call/loop/forpkg"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("should call forInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedSCGLoop := &model.SCGLoop{
				For: &model.SCGLoopFor{},
			}

			providedOpHandle := new(data.FakeHandle)
			providedScope := map[string]*model.Value{}

			fakeForInterpreter := new(forpkg.FakeInterpreter)

			objectUnderTest := _interpreter{
				forInterpreter:        fakeForInterpreter,
				predicatesInterpreter: new(predicates.FakeInterpreter),
			}

			/* act */
			objectUnderTest.Interpret(
				providedOpHandle,
				providedSCGLoop,
				providedScope,
			)

			/* assert */
			actualOpHandle,
				actualSCGLoopFor,
				actualScope := fakeForInterpreter.InterpretArgsForCall(0)

			Expect(actualOpHandle).To(Equal(providedOpHandle))
			Expect(actualSCGLoopFor).To(Equal(providedSCGLoop.For))
			Expect(actualScope).To(Equal(providedScope))
		})
		Context("forInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeForInterpreter := new(forpkg.FakeInterpreter)

				expectedError := errors.New("expectedError")
				fakeForInterpreter.InterpretReturns(
					nil,
					expectedError,
				)

				objectUnderTest := _interpreter{
					forInterpreter:        fakeForInterpreter,
					predicatesInterpreter: new(predicates.FakeInterpreter),
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					new(data.FakeHandle),
					&model.SCGLoop{
						For: &model.SCGLoopFor{},
					},
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		It("should call predicatesInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedSCGLoop := &model.SCGLoop{
				Until: []*model.SCGPredicate{},
			}

			providedOpHandle := new(data.FakeHandle)
			providedScope := map[string]*model.Value{}

			fakePredicatesInterpreter := new(predicates.FakeInterpreter)

			objectUnderTest := _interpreter{
				forInterpreter:        new(forpkg.FakeInterpreter),
				predicatesInterpreter: fakePredicatesInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedOpHandle,
				providedSCGLoop,
				providedScope,
			)

			/* assert */
			actualOpHandle,
				actualSCGLoopUntil,
				actualScope := fakePredicatesInterpreter.InterpretArgsForCall(0)

			Expect(actualOpHandle).To(Equal(providedOpHandle))
			Expect(actualSCGLoopUntil).To(Equal(providedSCGLoop.Until))
			Expect(actualScope).To(Equal(providedScope))
		})
		Context("predicatesInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakePredicatesInterpreter := new(predicates.FakeInterpreter)

				expectedError := errors.New("expectedError")
				fakePredicatesInterpreter.InterpretReturns(
					false,
					expectedError,
				)

				objectUnderTest := _interpreter{
					forInterpreter:        new(forpkg.FakeInterpreter),
					predicatesInterpreter: fakePredicatesInterpreter,
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					new(data.FakeHandle),
					&model.SCGLoop{
						Until: []*model.SCGPredicate{},
					},
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		It("should return expected result", func() {
			/* arrange */
			expectedIndex := "Index"
			providedScgLoop := &model.SCGLoop{
				For:   &model.SCGLoopFor{},
				Index: &expectedIndex,
				Until: []*model.SCGPredicate{},
			}

			fakeForInterpreter := new(forpkg.FakeInterpreter)

			expectedDCGLoopFor := &model.DCGLoopFor{}
			fakeForInterpreter.InterpretReturns(
				expectedDCGLoopFor,
				nil,
			)

			fakePredicatesInterpreter := new(predicates.FakeInterpreter)

			expectedDCGLoopUntil := false
			fakePredicatesInterpreter.InterpretReturns(
				expectedDCGLoopUntil,
				nil,
			)

			expectedResult := &model.DCGLoop{
				For:   expectedDCGLoopFor,
				Index: providedScgLoop.Index,
				Until: &expectedDCGLoopUntil,
			}

			objectUnderTest := _interpreter{
				forInterpreter:        fakeForInterpreter,
				predicatesInterpreter: fakePredicatesInterpreter,
			}

			/* act */
			actualResult, _ := objectUnderTest.Interpret(
				new(data.FakeHandle),
				providedScgLoop,
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
	})
})
