package image

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	stringPkg "github.com/opspec-io/sdk-golang/string"
)

var _ = Context("Image", func() {
	Context("Interpret", func() {

		Context("scgContainerCallImage is nil", func() {
			It("should return expected error", func() {
				/* arrange */
				objectUnderTest := _Image{
					string: new(stringPkg.Fake),
				}

				/* act */
				_, actualError := objectUnderTest.Interpret(
					map[string]*model.Value{},
					nil,
				)

				/* assert */
				Expect(actualError).To(Equal(errors.New("image required")))
			})
		})
		Context("scgContainerCallImage isn't nill", func() {
			It("should call string.Interpret w/ expected args", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Value{
					"name1": {String: &providedString1},
				}

				providedSCGContainerCallImage := &model.SCGContainerCallImage{
					Ref: "dummyImageRef",
					PullCreds: &model.SCGPullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				}

				fakeString := new(stringPkg.Fake)

				objectUnderTest := _Image{
					string: fakeString,
				}

				/* act */
				objectUnderTest.Interpret(
					providedCurrentScope,
					providedSCGContainerCallImage,
				)

				/* assert */
				actualImageRefScope, actualImageRef := fakeString.InterpretArgsForCall(0)
				Expect(actualImageRef).To(Equal(providedSCGContainerCallImage.Ref))
				Expect(actualImageRefScope).To(Equal(providedCurrentScope))

				actualUsernameScope, actualUsername := fakeString.InterpretArgsForCall(1)
				Expect(actualUsername).To(Equal(providedSCGContainerCallImage.PullCreds.Username))
				Expect(actualUsernameScope).To(Equal(providedCurrentScope))

				actualPasswordScope, actualPassword := fakeString.InterpretArgsForCall(2)
				Expect(actualPassword).To(Equal(providedSCGContainerCallImage.PullCreds.Password))
				Expect(actualPasswordScope).To(Equal(providedCurrentScope))
			})
			It("should return expected dcg.Image", func() {

				/* arrange */
				providedSCGContainerCallImage := &model.SCGContainerCallImage{
					Ref:       "dummyImageRef",
					PullCreds: &model.SCGPullCreds{},
				}

				fakeString := new(stringPkg.Fake)

				expectedImageRef := "expectedImageRef"
				fakeString.InterpretReturnsOnCall(0, expectedImageRef, nil)

				expectedUsername := "expectedUsername"
				fakeString.InterpretReturnsOnCall(1, expectedUsername, nil)

				expectedPassword := "expectedPassword"
				fakeString.InterpretReturnsOnCall(2, expectedPassword, nil)

				expectedImage := &model.DCGContainerCallImage{
					Ref: expectedImageRef,
					PullCreds: &model.DCGPullCreds{
						Username: expectedUsername,
						Password: expectedPassword,
					},
				}

				objectUnderTest := _Image{
					string: fakeString,
				}

				/* act */
				actualDCGContainerCallImage, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGContainerCallImage,
				)

				/* assert */
				Expect(actualDCGContainerCallImage).To(Equal(expectedImage))
			})
		})
	})
})
