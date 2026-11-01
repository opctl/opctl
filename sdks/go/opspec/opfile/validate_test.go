package opfile

import (
	"fmt"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Validate caches", func() {
	It("should accept a shorthand cache name", func() {
		Expect(Validate([]byte(`
name: cache-example
run:
  container:
    image: { ref: node:22 }
    caches:
      /root/.npm:
`))).To(BeEmpty())
	})

	for _, caches := range []map[string]interface{}{
		{"relative/path": "gomod"},
		{"/cache": 42},
		{"/cache": true},
		{"/cache": []interface{}{}},
		{"/cache": map[string]interface{}{"unknown": "locked"}},
		{"/cache": map[string]interface{}{"sharing": "private"}},
		{"/cache": map[string]interface{}{"sharing": nil}},
		{"/cache": map[string]interface{}{"sharing": true}},
		{"/cache": map[string]interface{}{"name": 42}},
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
