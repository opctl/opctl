package image

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	Context("containerCallImageSpec is nil", func() {
		It("should return expected error", func() {
			/* arrange */
			/* act */
			_, actualError := Interpret(
				map[string]*model.Value{},
				nil,
				"dummyScratchDir",
			)

			/* assert */
			Expect(actualError).To(MatchError("image required"))
		})
	})
	Context("containerCallImageSpec isn't nil", func() {
		It("should return expected result", func() {

			/* arrange */
			archVariable := "archVariable"
			archValue := "archValue"
			refVariable := "refVariable"
			refValue := "refValue"
			usernameVariable := "usernameVariable"
			usernameValue := "usernameValue"
			passwordVariable := "passwordVariable"
			passwordValue := "passwordValue"

			providedScope := map[string]*model.Value{
				archVariable: {
					String: &archValue,
				},
				usernameVariable: {
					String: &usernameValue,
				},
				passwordVariable: {
					String: &passwordValue,
				},
				refVariable: {
					String: &refValue,
				},
			}

			providedContainerCallImageSpec := &model.ContainerCallImageSpec{
				Ref: fmt.Sprintf("$(%s)", refVariable),
				Platform: &model.OCIImagePlatformSpec{
					Arch: fmt.Sprintf("$(%s)", archVariable),
				},
				PullCreds: &model.CredsSpec{
					Username: fmt.Sprintf("$(%s)", usernameVariable),
					Password: fmt.Sprintf("$(%s)", passwordVariable),
				},
			}

			parsedImageRef, err := reference.ParseAnyReference(strings.ToLower(refValue))
			if err != nil {
				panic(err)
			}

			expectedImageRef := parsedImageRef.String()

			expectedImage := &model.ContainerCallImage{
				Ref: &expectedImageRef,
				Platform: &model.OCIImagePlatform{
					Arch: &archValue,
				},
				PullCreds: &model.Creds{
					Username: usernameValue,
					Password: passwordValue,
				},
			}

			/* act */
			actualContainerCallImage, actualErr := Interpret(
				providedScope,
				providedContainerCallImageSpec,
				"dummyScratchDir",
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(*actualContainerCallImage).To(Equal(*expectedImage))
		})
	})
})
