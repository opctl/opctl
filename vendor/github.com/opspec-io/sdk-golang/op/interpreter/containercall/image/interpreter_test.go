package image

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/interpreter/expression"
)

var _ = Context("Interpreter", func() {
	Context("NewInterpreter", func() {
		It("shouldn't return nil", func() {
			/* arrange/act/assert */
			Expect(NewInterpreter()).To(Not(BeNil()))
		})
	})
	Context("Interpret", func() {
		Context("scgContainerCallImage is nil", func() {
			It("should return expected error", func() {
				/* arrange */
				objectUnderTest := _interpreter{
					expression: new(expression.Fake),
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					nil,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualError).To(Equal(errors.New("image required")))
			})
		})
		Context("scgContainerCallImage isn't nill", func() {
			It("should call expression.EvalToString w/ expected args", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Value{
					"name1": {String: &providedString1},
				}

				providedOpHandle := new(data.FakeHandle)

				providedSCGContainerCallImage := &model.SCGContainerCallImage{
					Ref: "dummyImageRef",
					PullCreds: &model.SCGPullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				}

				fakeExpression := new(expression.Fake)
				fakeExpression.EvalToStringReturns(&model.Value{String: new(string)}, nil)

				objectUnderTest := _interpreter{
					expression: fakeExpression,
				}

				/* act */
				objectUnderTest.Interpret(
					providedCurrentScope,
					providedSCGContainerCallImage,
					providedOpHandle,
				)

				/* assert */
				actualImageRefScope,
					actualImageRef,
					actualImageRefOpHandle := fakeExpression.EvalToStringArgsForCall(0)
				Expect(actualImageRef).To(Equal(providedSCGContainerCallImage.Ref))
				Expect(actualImageRefScope).To(Equal(providedCurrentScope))
				Expect(actualImageRefOpHandle).To(Equal(providedOpHandle))

				actualUsernameScope,
					actualUsername,
					actualUsernameOpHandle := fakeExpression.EvalToStringArgsForCall(1)
				Expect(actualUsername).To(Equal(providedSCGContainerCallImage.PullCreds.Username))
				Expect(actualUsernameScope).To(Equal(providedCurrentScope))
				Expect(actualUsernameOpHandle).To(Equal(providedOpHandle))

				actualPasswordScope,
					actualPassword,
					actualPasswordOpHandle := fakeExpression.EvalToStringArgsForCall(2)
				Expect(actualPassword).To(Equal(providedSCGContainerCallImage.PullCreds.Password))
				Expect(actualPasswordScope).To(Equal(providedCurrentScope))
				Expect(actualPasswordOpHandle).To(Equal(providedOpHandle))
			})
			It("should return expected dcg.Image", func() {

				/* arrange */
				providedSCGContainerCallImage := &model.SCGContainerCallImage{
					Ref:       "dummyImageRef",
					PullCreds: &model.SCGPullCreds{},
				}

				fakeExpression := new(expression.Fake)

				expectedImageRef := "expectedImageRef"
				fakeExpression.EvalToStringReturnsOnCall(0, &model.Value{String: &expectedImageRef}, nil)

				expectedUsername := "expectedUsername"
				fakeExpression.EvalToStringReturnsOnCall(1, &model.Value{String: &expectedUsername}, nil)

				expectedPassword := "expectedPassword"
				fakeExpression.EvalToStringReturnsOnCall(2, &model.Value{String: &expectedPassword}, nil)

				expectedImage := &model.DCGContainerCallImage{
					Ref: expectedImageRef,
					PullCreds: &model.DCGPullCreds{
						Username: expectedUsername,
						Password: expectedPassword,
					},
				}

				objectUnderTest := _interpreter{
					expression: fakeExpression,
				}

				/* act */
				actualDCGContainerCallImage, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGContainerCallImage,
					new(data.FakeHandle),
				)

				/* assert */
				Expect(actualDCGContainerCallImage).To(Equal(expectedImage))
			})
		})
	})
})
