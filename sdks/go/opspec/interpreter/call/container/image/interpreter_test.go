package image

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/model"
	stringPkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/string"
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
				objectUnderTest := _interpreter{}

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
			It("should call stringInterpreter.Interpret w/ expected args", func() {
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

				fakeStringInterpreter := new(stringPkg.FakeInterpreter)
				fakeStringInterpreter.InterpretReturns(&model.Value{String: new(string)}, nil)

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
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
					actualImageRefOpHandle := fakeStringInterpreter.InterpretArgsForCall(0)
				Expect(actualImageRef).To(Equal(providedSCGContainerCallImage.Ref))
				Expect(actualImageRefScope).To(Equal(providedCurrentScope))
				Expect(actualImageRefOpHandle).To(Equal(providedOpHandle))

				actualUsernameScope,
					actualUsername,
					actualUsernameOpHandle := fakeStringInterpreter.InterpretArgsForCall(1)
				Expect(actualUsername).To(Equal(providedSCGContainerCallImage.PullCreds.Username))
				Expect(actualUsernameScope).To(Equal(providedCurrentScope))
				Expect(actualUsernameOpHandle).To(Equal(providedOpHandle))

				actualPasswordScope,
					actualPassword,
					actualPasswordOpHandle := fakeStringInterpreter.InterpretArgsForCall(2)
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

				fakeStringInterpreter := new(stringPkg.FakeInterpreter)

				expectedImageRef := "expectedImageRef"
				fakeStringInterpreter.InterpretReturnsOnCall(0, &model.Value{String: &expectedImageRef}, nil)

				expectedUsername := "expectedUsername"
				fakeStringInterpreter.InterpretReturnsOnCall(1, &model.Value{String: &expectedUsername}, nil)

				expectedPassword := "expectedPassword"
				fakeStringInterpreter.InterpretReturnsOnCall(2, &model.Value{String: &expectedPassword}, nil)

				expectedImage := &model.DCGContainerCallImage{
					Ref: expectedImageRef,
					PullCreds: &model.PullCreds{
						Username: expectedUsername,
						Password: expectedPassword,
					},
				}

				objectUnderTest := _interpreter{
					stringInterpreter: fakeStringInterpreter,
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
