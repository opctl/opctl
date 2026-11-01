package opfile

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("container stream schema", func() {
	It("accepts only whole identifier write targets", func() {
		for _, field := range []string{"stdout", "stderr"} {
			for _, target := range []string{`"$(version)"`, `"$(_A-12)"`} {
				Expect(Validate([]byte(fmt.Sprintf(`{"name":"test","run":{"container":{"image":{"ref":"alpine"},%q:%s}}}`, field, target)))).To(BeEmpty())
			}
			for _, target := range []string{`"version"`, `null`, `""`, `"$(v/path)"`, `"$(v.key)"`, `"$(v) suffix"`, `"$($(v))"`} {
				Expect(Validate([]byte(fmt.Sprintf(`{"name":"test","run":{"container":{"image":{"ref":"alpine"},%q:%s}}}`, field, target)))).NotTo(BeEmpty())
			}
		}
	})
})
