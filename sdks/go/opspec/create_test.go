package opspec

import (
	"io/ioutil"
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
		providedPath, err := ioutil.TempDir("", "")
		if nil != err {
			panic(err)
		}
		providedPkgName := "dummyPkgName"
		providedPkgDescription := "dummyPkgDescription"

		expectedOpFileBytes, err := yaml.Marshal(&model.OpSpec{
			Description: providedPkgDescription,
			Name:        providedPkgName,
		})
		if nil != err {
			panic(err)
		}

		/* act */
		Create(
			providedPath,
			providedPkgName,
			providedPkgDescription,
		)

		/* assert */
		actualOpFileBytes, err := ioutil.ReadFile(filepath.Join(providedPath, opfile.FileName))
		if nil != err {
			panic(err)
		}

		Expect(string(actualOpFileBytes)).To(Equal(string(expectedOpFileBytes)))
	})

})
