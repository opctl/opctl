package opspec

import (
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/opspec/opfile"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Create", func() {

	It("should create expected op", func() {

		/* arrange */
		providedPath, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}
		providedPkgName := "dummyPkgName"
		providedPkgDescription := "dummyPkgDescription"

		expectedOpFileBytes, err := yaml.Marshal(&model.OpSpec{
			Description: providedPkgDescription,
			Name:        providedPkgName,
		})
		if err != nil {
			panic(err)
		}

		/* act */
		Create(
			providedPath,
			providedPkgName,
			providedPkgDescription,
		)

		/* assert */
		actualOpFileBytes, err := os.ReadFile(filepath.Join(providedPath, opfile.FileName))
		if err != nil {
			panic(err)
		}

		Expect(string(actualOpFileBytes)).To(Equal(string(expectedOpFileBytes)))
	})

})
