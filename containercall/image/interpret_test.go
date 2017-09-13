package image

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
)

var _ = Context("Image", func() {
	Context("Interpret", func() {

		Context("scgContainerCallImage is nil", func() {
			It("should return expected error", func() {
				/* arrange */
				objectUnderTest := _Image{
					expression: new(expression.Fake),
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					nil,
					new(pkg.FakeHandle),
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

				providedPkgHandle := new(pkg.FakeHandle)

				providedSCGContainerCallImage := &model.SCGContainerCallImage{
					Ref: "dummyImageRef",
					PullCreds: &model.SCGPullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				}

				fakeExpression := new(expression.Fake)

				objectUnderTest := _Image{
					expression: fakeExpression,
				}

				/* act */
				objectUnderTest.Interpret(
					providedCurrentScope,
					providedSCGContainerCallImage,
					providedPkgHandle,
				)

				/* assert */
				actualImageRefScope,
					actualImageRef,
					actualImageRefPkgHandle := fakeExpression.EvalToStringArgsForCall(0)
				Expect(actualImageRef).To(Equal(providedSCGContainerCallImage.Ref))
				Expect(actualImageRefScope).To(Equal(providedCurrentScope))
				Expect(actualImageRefPkgHandle).To(Equal(providedPkgHandle))

				actualUsernameScope,
					actualUsername,
					actualUsernamePkgHandle := fakeExpression.EvalToStringArgsForCall(1)
				Expect(actualUsername).To(Equal(providedSCGContainerCallImage.PullCreds.Username))
				Expect(actualUsernameScope).To(Equal(providedCurrentScope))
				Expect(actualUsernamePkgHandle).To(Equal(providedPkgHandle))

				actualPasswordScope,
					actualPassword,
					actualPasswordPkgHandle := fakeExpression.EvalToStringArgsForCall(2)
				Expect(actualPassword).To(Equal(providedSCGContainerCallImage.PullCreds.Password))
				Expect(actualPasswordScope).To(Equal(providedCurrentScope))
				Expect(actualPasswordPkgHandle).To(Equal(providedPkgHandle))
			})
			It("should return expected dcg.Image", func() {

				/* arrange */
				providedSCGContainerCallImage := &model.SCGContainerCallImage{
					Ref:       "dummyImageRef",
					PullCreds: &model.SCGPullCreds{},
				}

				fakeExpression := new(expression.Fake)

				expectedImageRef := "expectedImageRef"
				fakeExpression.EvalToStringReturnsOnCall(0, expectedImageRef, nil)

				expectedUsername := "expectedUsername"
				fakeExpression.EvalToStringReturnsOnCall(1, expectedUsername, nil)

				expectedPassword := "expectedPassword"
				fakeExpression.EvalToStringReturnsOnCall(2, expectedPassword, nil)

				expectedImage := &model.DCGContainerCallImage{
					Ref: expectedImageRef,
					PullCreds: &model.DCGPullCreds{
						Username: expectedUsername,
						Password: expectedPassword,
					},
				}

				objectUnderTest := _Image{
					expression: fakeExpression,
				}

				/* act */
				actualDCGContainerCallImage, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGContainerCallImage,
					new(pkg.FakeHandle),
				)

				/* assert */
				Expect(actualDCGContainerCallImage).To(Equal(expectedImage))
			})
		})
	})
})
