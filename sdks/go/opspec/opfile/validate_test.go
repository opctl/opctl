package opfile

import (
	"fmt"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate caches", func() {
	for _, caches := range []map[string]interface{}{
		{"relative/path": "gomod"},
		{"/cache": nil},
		{"/cache": 42},
		{"/cache": map[string]interface{}{}},
	} {
		It(fmt.Sprintf("should reject %v", caches), func() {
			providedBytes, err := yaml.Marshal(model.OpSpec{
				Name: "cache-example",
				Run: &model.CallSpec{
					Container: &model.ContainerCallSpec{
						Image:  &model.ContainerCallImageSpec{Ref: "golang:1.23"},
						Caches: caches,
					},
				},
			})
			Expect(err).To(BeNil())

			Expect(Validate(providedBytes)).NotTo(BeEmpty())
		})
	}
})
