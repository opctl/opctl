package cliparamsatisfier

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Describe("inputSrcFactory", func() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	argsYmlTestDataPath := filepath.Join(wd, "inputsrc/ymlfile/testdata/args.yml")
	Context("NewCLIPromptInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_inputSrcFactory{}.NewCliPromptInputSrc(nil)).To(Not(BeNil()))
		})
	})
	Context("NewEnvVarInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_inputSrcFactory{}.NewEnvVarInputSrc()).To(Not(BeNil()))
		})
	})
	Context("NewParamDefaultInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_inputSrcFactory{}.NewParamDefaultInputSrc(
				map[string]*model.Param{},
			)).To(Not(BeNil()))
		})
	})
	Context("NewSliceInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_inputSrcFactory{}.NewSliceInputSrc([]string{}, "")).To(Not(BeNil()))
		})
	})
	Context("NewYMLFileInputSrc()", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(_inputSrcFactory{}.NewYMLFileInputSrc(argsYmlTestDataPath)).To(Not(BeNil()))
		})
	})
})
