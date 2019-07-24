package dir

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/types"
)

var _ = Context("Interpret", func() {
	Context("expression is scope ref", func() {
		It("should call referenceInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedScope := map[string]*types.Value{"": new(types.Value)}

			providedExpression := "$(providedReference)"
			providedOpRef := new(data.FakeHandle)

			expectedReference := strings.TrimSuffix(strings.TrimPrefix(providedExpression, interpolater.RefStart), interpolater.RefEnd)

			fakeReferenceInterpreter := new(reference.FakeInterpreter)
			// err to trigger immediate return
			fakeReferenceInterpreter.InterpretReturns(nil, errors.New("dummyError"))

			objectUnderTest := _interpreter{
				referenceInterpreter: fakeReferenceInterpreter,
			}

			/* act */
			objectUnderTest.Interpret(
				providedScope,
				providedExpression,
				providedOpRef,
			)

			/* assert */
			expectedReference,
				actualScope,
				actualOpRef := fakeReferenceInterpreter.InterpretArgsForCall(0)

			Expect(expectedReference).To(Equal(expectedReference))
			Expect(actualScope).To(Equal(providedScope))
			Expect(actualOpRef).To(Equal(providedOpRef))
		})
		Context("referenceInterpreter.Interpret errs", func() {
			It("should return expected result", func() {

				/* arrange */
				providedScope := map[string]*types.Value{"": new(types.Value)}

				providedExpression := "$(providedReference)"

				fakeReferenceInterpreter := new(reference.FakeInterpreter)
				interpretErr := errors.New("dummyError")
				fakeReferenceInterpreter.InterpretReturns(nil, interpretErr)

				expectedErr := fmt.Errorf("unable to interpret %+v to dir; error was %v", providedExpression, interpretErr)

				objectUnderTest := _interpreter{
					referenceInterpreter: fakeReferenceInterpreter,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					providedScope,
					providedExpression,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(expectedErr))

			})
		})
		Context("referenceInterpreter.Interpret doesn't error", func() {
			It("should return expected result", func() {
				/* arrange */
				expectedValue := new(types.Value)
				fakeReferenceInterpreter := new(reference.FakeInterpreter)
				fakeReferenceInterpreter.InterpretReturns(expectedValue, nil)

				objectUnderTest := _interpreter{
					referenceInterpreter: fakeReferenceInterpreter,
				}

				/* act */
				actualValue, actualErr := objectUnderTest.Interpret(
					map[string]*types.Value{},
					"$(providedReference)",
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualValue).To(Equal(expectedValue))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
