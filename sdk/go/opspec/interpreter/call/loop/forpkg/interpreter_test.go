package forpkg

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/loopable"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("should call loopableInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedSCGLoopFor := &model.SCGLoopFor{
				Each: "Each",
			}

			providedOpHandle := new(data.FakeHandle)
			providedScope := map[string]*model.Value{}

			fakeLoopableInterpreter := new(loopable.FakeInterpreter)

			objectUnderTest := _interpreter{
				loopableInterpreter: fakeLoopableInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedOpHandle,
				providedSCGLoopFor,
				providedScope,
			)

			/* assert */
			actualSCGLoopForEach,
				actualOpHandle,
				actualScope := fakeLoopableInterpreter.InterpretArgsForCall(0)

			Expect(actualSCGLoopForEach).To(Equal(providedSCGLoopFor.Each))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
			Expect(actualScope).To(Equal(providedScope))
		})
		Context("loopableInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeLoopableInterpreter := new(loopable.FakeInterpreter)

				expectedError := errors.New("expectedError")
				fakeLoopableInterpreter.InterpretReturns(
					nil,
					expectedError,
				)

				objectUnderTest := _interpreter{
					loopableInterpreter: fakeLoopableInterpreter,
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					new(data.FakeHandle),
					&model.SCGLoopFor{},
					map[string]*model.Value{},
				)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		It("should return expected result", func() {
			/* arrange */
			providedSCGLoopFor := &model.SCGLoopFor{
				Value: new(string),
			}

			fakeLoopableInterpreter := new(loopable.FakeInterpreter)

			expectedDCGLoopForEach := &model.Value{}
			fakeLoopableInterpreter.InterpretReturns(
				expectedDCGLoopForEach,
				nil,
			)

			expectedResult := &model.DCGLoopFor{
				Each:  expectedDCGLoopForEach,
				Value: providedSCGLoopFor.Value,
			}

			objectUnderTest := _interpreter{
				loopableInterpreter: fakeLoopableInterpreter,
			}

			/* act */
			actualResult, _ := objectUnderTest.Interpret(
				new(data.FakeHandle),
				providedSCGLoopFor,
				map[string]*model.Value{},
			)

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))
		})
	})
})
