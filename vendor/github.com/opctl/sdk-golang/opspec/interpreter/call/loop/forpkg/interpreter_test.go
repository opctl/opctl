package forpkg

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/array"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		It("should call arrayInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedSCGLoopFor := &model.SCGLoopFor{
				Each: "Each",
			}

			providedOpHandle := new(data.FakeHandle)
			providedScope := map[string]*model.Value{}

			fakeArrayInterpreter := new(array.FakeInterpreter)

			objectUnderTest := _interpreter{
				arrayInterpreter: fakeArrayInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedOpHandle,
				providedSCGLoopFor,
				providedScope,
			)

			/* assert */
			actualScope,
				actualSCGLoopForEach,
				actualOpHandle := fakeArrayInterpreter.InterpretArgsForCall(0)

			Expect(actualScope).To(Equal(providedScope))
			Expect(actualSCGLoopForEach).To(Equal(providedSCGLoopFor.Each))
			Expect(actualOpHandle).To(Equal(providedOpHandle))
		})
		Context("arrayInterpreter.Interpret errs", func() {
			It("should return expected result", func() {
				/* arrange */
				fakeArrayInterpreter := new(array.FakeInterpreter)

				expectedError := errors.New("expectedError")
				fakeArrayInterpreter.InterpretReturns(
					nil,
					expectedError,
				)

				objectUnderTest := _interpreter{
					arrayInterpreter: fakeArrayInterpreter,
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

			fakeArrayInterpreter := new(array.FakeInterpreter)

			expectedDCGLoopForEach := &model.Value{}
			fakeArrayInterpreter.InterpretReturns(
				expectedDCGLoopForEach,
				nil,
			)

			expectedResult := &model.DCGLoopFor{
				Each:  expectedDCGLoopForEach,
				Value: providedSCGLoopFor.Value,
			}

			objectUnderTest := _interpreter{
				arrayInterpreter: fakeArrayInterpreter,
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
