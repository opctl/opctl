package op

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Creater", func() {
	Context("Create", func() {
		It("should create expected op", func() {
			/* arrange */
			providedPath, err := ioutil.TempDir("", "")
			if nil != err {
				panic(err)
			}
			providedName := "dummyName"
			providedDescription := "dummyDescription"

			expectedOpFileBytes, err := yaml.Marshal(&model.OpSpec{
				Description: providedDescription,
				Name:        providedName,
			})
			if nil != err {
				panic(err)
			}

			objectUnderTest := newCreater()

			/* act */
			err = objectUnderTest.Create(providedPath, providedDescription, providedName)

			/* assert */
			Expect(err).To(BeNil())
			actualOpFileBytes, err := ioutil.ReadFile(filepath.Join(providedPath, providedName, "op.yml"))
			if nil != err {
				panic(err)
			}
			Expect(string(actualOpFileBytes)).To(Equal(string(expectedOpFileBytes)))
		})
	})
})
