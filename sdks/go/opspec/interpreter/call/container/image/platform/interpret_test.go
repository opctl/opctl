package platform

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	It("should return expected result", func() {
		/* arrange */
		archVariable := "archVariable"
		archValue := "archValue"

		providedScope := map[string]*model.Value{
			archVariable: {
				String: &archValue,
			},
		}

		providedImagePlatformSpec := &model.OCIImagePlatformSpec{
			Arch: fmt.Sprintf("$(%s)", archVariable),
		}

		expectedImagePlatform := &model.OCIImagePlatform{
			Arch: &archValue,
		}

		/* act */
		actualContainerCallImage, actualErr := Interpret(
			providedScope,
			providedImagePlatformSpec,
			"dummyScratchDir",
		)

		/* assert */
		Expect(actualErr).To(BeNil())
		Expect(actualContainerCallImage).To(Equal(expectedImagePlatform))
	})
})
