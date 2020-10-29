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
		Context("containerCallImageSpec is nil", func() {
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
		Context("containerCallImageSpec isn't nil", func() {
			It("should call stringInterpreter.Interpret w/ expected args", func() {
				/* arrange */
				providedString1 := "dummyString1"
				providedCurrentScope := map[string]*model.Value{
					"name1": {String: &providedString1},
				}

				providedContainerCallImageSpec := &model.ContainerCallImageSpec{
					PullCreds: &model.CredsSpec{
						Username: "dummyUsername",
						Password: "dummyPassword",
					},
				}

				fakeDirInterpreter := new(dirFakes.FakeInterpreter)
				fakeDirInterpreter.InterpretReturns(nil, errors.New("dummyError"))

				fakeStrInterpreter := new(strFakes.FakeInterpreter)
				imageRef := "imageRef"
				fakeStrInterpreter.InterpretReturns(&model.Value{String: &imageRef}, nil)

				objectUnderTest := _interpreter{
					dirInterpreter:    fakeDirInterpreter,
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				_, err := objectUnderTest.Interpret(
					providedCurrentScope,
					providedContainerCallImageSpec,
					"dummyScratchDir",
				)

				/* assert */
				Expect(err).To(BeNil())
				actualImageRefScope,
					actualImageRef := fakeStrInterpreter.InterpretArgsForCall(0)
				Expect(actualImageRef).To(Equal(providedContainerCallImageSpec.Ref))
				Expect(actualImageRefScope).To(Equal(providedCurrentScope))

				actualUsernameScope,
					actualUsername := fakeStrInterpreter.InterpretArgsForCall(1)
				Expect(actualUsername).To(Equal(providedContainerCallImageSpec.PullCreds.Username))
				Expect(actualUsernameScope).To(Equal(providedCurrentScope))

				actualPasswordScope,
					actualPassword := fakeStrInterpreter.InterpretArgsForCall(2)
				Expect(actualPassword).To(Equal(providedContainerCallImageSpec.PullCreds.Password))
				Expect(actualPasswordScope).To(Equal(providedCurrentScope))
			})
			It("should return expected dcg.Image", func() {

				/* arrange */
				providedContainerCallImageSpec := &model.ContainerCallImageSpec{
					Ref:       "dummyRef",
					PullCreds: &model.CredsSpec{},
				}

				fakeDirInterpreter := new(dirFakes.FakeInterpreter)
				fakeDirInterpreter.InterpretReturns(nil, errors.New("dummyError"))

				fakeStrInterpreter := new(strFakes.FakeInterpreter)

				expectedImageRef := "docker.io/library/expectedimageref"
				fakeStrInterpreter.InterpretReturnsOnCall(0, &model.Value{String: &expectedImageRef}, nil)

				expectedUsername := "expectedUsername"
				fakeStrInterpreter.InterpretReturnsOnCall(1, &model.Value{String: &expectedUsername}, nil)

				expectedPassword := "expectedPassword"
				fakeStrInterpreter.InterpretReturnsOnCall(2, &model.Value{String: &expectedPassword}, nil)

				expectedImage := &model.ContainerCallImage{
					Ref: &expectedImageRef,
					PullCreds: &model.Creds{
						Username: expectedUsername,
						Password: expectedPassword,
					},
				}

				objectUnderTest := _interpreter{
					dirInterpreter:    fakeDirInterpreter,
					stringInterpreter: fakeStrInterpreter,
				}

				/* act */
				actualContainerCallImage, actualErr := objectUnderTest.Interpret(
					map[string]*model.Value{},
					providedContainerCallImageSpec,
					"dummyScratchDir",
				)

				/* assert */
				Expect(actualErr).To(BeNil())
				Expect(actualContainerCallImage).To(Equal(expectedImage))
			})
		})
	})
})
