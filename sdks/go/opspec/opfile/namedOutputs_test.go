package opfile

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("outputs file binding schema", func() {
	It("requires a nonempty explicit identifier map and reserves declarations", func() {
		for _, value := range []string{`{"version":"$(v)"}`, `{"count":"$(n)","ready":"$(b)"}`} {
			Expect(Validate([]byte(fmt.Sprintf(`{"name":"test","run":{"container":{"image":{"ref":"alpine"},"outputs":%s}}}`, value)))).To(BeEmpty())
		}
		for _, value := range []string{`{}`, `true`, `null`, `{"a":null}`, `{"a":"bare"}`, `{"a/b":"$(v)"}`, `{"a":"$(v.key)"}`} {
			Expect(Validate([]byte(fmt.Sprintf(`{"name":"test","run":{"container":{"image":{"ref":"alpine"},"outputs":%s}}}`, value)))).NotTo(BeEmpty())
		}
		Expect(Validate([]byte(`{"name":"test","run":{"container":{"image":{"ref":"alpine"},"outputParams":{"x":{"string":{}}}}}}`))).NotTo(BeEmpty())
	})
})
