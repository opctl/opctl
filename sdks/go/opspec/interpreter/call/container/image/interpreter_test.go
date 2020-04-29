package image

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	dirFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/dir/fakes"
	strFakes "github.com/opctl/opctl/sdks/go/opspec/interpreter/str/fakes"
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
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualError).To(Equal(fmt.Errorf("image required")))
			})
		})
		Context("scgContainerCallImage isn't nil", func() {
			It("should call stringInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Value{
					"name1": {String: &providedString1},
				}

				providedSCGContainerCallImage := &model.SCGContainerCallImage{
					Ref: "dummyRef",
					PullCreds: &model.SCGPullCreds{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				}

				fakeDirInterpreter := new(dirFakes.FakeInterpreter)
				fakeDirInterpreter.InterpretReturns(nil, errors.New("dummyError"))

				fakeStrInterpreter := new(strFakes.FakeInterpreter)
				fakeStrInterpreter.InterpretReturns(&model.Value{String: new(string)}, nil)

				objectUnderTest := _interpreter{
					dirInterpreter:    fakeDirInterpreter,
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				objectUnderTest.Interpret(
					providedCurrentScope,
					providedSCGContainerCallImage,
					"dummyScratchDir",
				)

				/* assert */
				actualImageRefScope,
					actualImageRef := fakeStrInterpreter.InterpretArgsForCall(0)
				Expect(actualImageRef).To(Equal(providedSCGContainerCallImage.Ref))
				Expect(actualImageRefScope).To(Equal(providedCurrentScope))

				actualUsernameScope,
					actualUsername := fakeStrInterpreter.InterpretArgsForCall(1)
				Expect(actualUsername).To(Equal(providedSCGContainerCallImage.PullCreds.Username))
				Expect(actualUsernameScope).To(Equal(providedCurrentScope))

				actualPasswordScope,
					actualPassword := fakeStrInterpreter.InterpretArgsForCall(2)
				Expect(actualPassword).To(Equal(providedSCGContainerCallImage.PullCreds.Password))
				Expect(actualPasswordScope).To(Equal(providedCurrentScope))
			})
			It("should return expected dcg.Image", func() {

				/* arrange */
				providedSCGContainerCallImage := &model.SCGContainerCallImage{
					Ref:       "dummyRef",
					PullCreds: &model.SCGPullCreds{},
				}

				fakeDirInterpreter := new(dirFakes.FakeInterpreter)
				fakeDirInterpreter.InterpretReturns(nil, errors.New("dummyError"))

				fakeStrInterpreter := new(strFakes.FakeInterpreter)

				expectedImageRef := "expectedImageRef"
				fakeStrInterpreter.InterpretReturnsOnCall(0, &model.Value{String: &expectedImageRef}, nil)

				expectedUsername := "expectedUsername"
				fakeStrInterpreter.InterpretReturnsOnCall(1, &model.Value{String: &expectedUsername}, nil)

				expectedPassword := "expectedPassword"
				fakeStrInterpreter.InterpretReturnsOnCall(2, &model.Value{String: &expectedPassword}, nil)

				expectedImage := &model.DCGContainerCallImage{
					Ref: &expectedImageRef,
					PullCreds: &model.PullCreds{
						Username: expectedUsername,
						Password: expectedPassword,
					},
				}

				objectUnderTest := _interpreter{
					dirInterpreter:    fakeDirInterpreter,
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				actualDCGContainerCallImage, _ := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedSCGContainerCallImage,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualDCGContainerCallImage).To(Equal(expectedImage))
			})
		})
	})
})
