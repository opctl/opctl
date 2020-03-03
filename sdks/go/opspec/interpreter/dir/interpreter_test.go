package dir

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	referenceFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/fakes"
)

var _ = Context("Interpret", func() {
	Context("expression is scope ref", func() {
		It("should call referenceInterpreter.Interpret w/ expected args", func() {
			/* arrange */
			providedScope := map[string]*model.Value{"": new(model.Value)}

			providedExpression := "$(providedReference)"
			providedOpRef := new(modelFakes.FakeDataHandle)

			expectedReference := strings.TrimSuffix(strings.TrimPrefix(providedExpression, interpolater.RefStart), interpolater.RefEnd)

			fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
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
				providedScope := map[string]*model.Value{"": new(model.Value)}

				providedExpression := "$(providedReference)"

				fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
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
					new(modelFakes.FakeDataHandle),
				)

				/* assert */
				Expect(actualValue).To(BeNil())
				Expect(actualErr).To(Equal(expectedErr))

			})
		})
		Context("referenceInterpreter.Interpret doesn't error", func() {
			Context("value.Dir nil", func() {
				It("should return expected result", func() {
					/* arrange */
					providedExpression := "$(providedReference)"
					expectedValue := &model.Value{}
					fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
					fakeReferenceInterpreter.InterpretReturns(expectedValue, nil)

					objectUnderTest := _interpreter{
						referenceInterpreter: fakeReferenceInterpreter,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						providedExpression,
						new(modelFakes.FakeDataHandle),
					)

					/* assert */
					Expect(actualValue).To(BeNil())
					Expect(actualErr).To(Equal(fmt.Errorf("unable to interpret %+v to dir", providedExpression)))
				})
			})
			Context("value.Dir not nil", func() {
				It("should return expected result", func() {
					/* arrange */
					expectedValue := &model.Value{
						Dir: new(string),
					}
					fakeReferenceInterpreter := new(referenceFakes.FakeInterpreter)
					fakeReferenceInterpreter.InterpretReturns(expectedValue, nil)

					objectUnderTest := _interpreter{
						referenceInterpreter: fakeReferenceInterpreter,
					}

					/* act */
					actualValue, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						"$(providedReference)",
						new(modelFakes.FakeDataHandle),
					)

					/* assert */
					Expect(actualValue).To(Equal(expectedValue))
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})
})
