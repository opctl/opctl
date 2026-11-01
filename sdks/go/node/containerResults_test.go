package node

import (
	"context"
	"github.com/opctl/opctl/sdks/go/model"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("container result collectors", func() {
	It("assembles UTF-8 and strips exactly one final newline", func() {
		for input, want := range map[string]string{"": "", "v\n": "v", "v\r\n": "v", "v\n\n": "v\n", " v \n": " v ", "v\r": "v\r", "é": "é"} {
			b := newResultStream("stdout", func() {})
			for _, c := range []byte(input) {
				_, err := b.Write([]byte{c})
				Expect(err).NotTo(HaveOccurred())
			}
			got, err := b.value()
			Expect(err).NotTo(HaveOccurred())
			Expect(*got.String).To(Equal(want))
		}
		b := newResultStream("stderr", func() {})
		b.Write([]byte{255})
		_, err := b.value()
		Expect(err).To(MatchError(ContainSubstring("UTF-8")))
	})
	It("bounds concurrent writes and cancels even when writer errors are ignored", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		b := newResultStream("stdout", cancel)
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); b.Write([]byte(strings.Repeat("x", containerResultLimit/8))) }()
		}
		wg.Wait()
		_, err := b.value()
		Expect(err).NotTo(HaveOccurred())
		b.Write([]byte("x"))
		Expect(ctx.Err()).To(HaveOccurred())
		_, err = b.value()
		Expect(err).To(MatchError(ContainSubstring("stdout exceeds 1 MiB")))
	})
	It("bounds stderr independently and retains exact bytes below the cap", func() {
		out := newResultStream("stdout", func() {})
		stderr := newResultStream("stderr", func() {})
		_, err := stderr.Write([]byte(strings.Repeat("e", containerResultLimit)))
		Expect(err).NotTo(HaveOccurred())
		_, err = out.Write([]byte("public"))
		Expect(err).NotTo(HaveOccurred())
		_, err = stderr.Write([]byte("x"))
		Expect(err).To(MatchError(ContainSubstring("stderr exceeds 1 MiB")))
		value, err := out.value()
		Expect(err).NotTo(HaveOccurred())
		Expect(*value.String).To(Equal("public"))
	})

	It("parses native JSON without evaluating references", func() {
		result, err := parseContainerOutputs([]byte(`{"s":"$(token)\n","n":3,"b":true,"a":["$(literal)",{}],"o":{"$(key)":"value"}}`))
		Expect(err).NotTo(HaveOccurred())
		Expect(*result["s"].String).To(Equal("$(token)\n"))
		Expect(*result["n"].Number).To(Equal(float64(3)))
		Expect(*result["b"].Boolean).To(BeTrue())
		Expect(*result["o"].Object).To(HaveKey("$(key)"))
		Expect(*result["a"].Array).To(HaveLen(2))
	})
	It("rejects invalid data without quoting any content", func() {
		for _, data := range []string{"", "\xef\xbb\xbf{}", "[]", "{} {}", "{}secret-marker", `{"secret-marker":}`, `{"a":null}`, `{"a":[null]}`, `{"a":{"x":null}}`, `{"a":1,"a":2}`, `{"a":{"x":1,"x":2}}`, `{"a":1e9999}`, "{\"a\":\"\xff\"}", strings.Repeat("[", 101) + strings.Repeat("]", 101)} {
			_, err := parseContainerOutputs([]byte(data))
			Expect(err).To(HaveOccurred(), data)
			Expect(err.Error()).NotTo(ContainSubstring("secret-marker"))
		}
		_, err := parseContainerOutputs([]byte(`{"a":` + strings.Repeat("[", 99) + strings.Repeat("]", 99) + `}`))
		Expect(err).NotTo(HaveOccurred())
		_, err = parseContainerOutputs([]byte(`{"a":` + strings.Repeat("[", 100) + strings.Repeat("]", 100) + `}`))
		Expect(err).To(HaveOccurred())
	})
	It("accepts the exact file cap and rejects cap plus one", func() {
		dir, err := os.MkdirTemp("", "outputs-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "values.json")
		data := `{}` + strings.Repeat(" ", containerResultLimit-2)
		Expect(os.WriteFile(path, []byte(data), 0600)).To(Succeed())
		_, err = readContainerOutputs(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(path, []byte(data+" "), 0600)).To(Succeed())
		_, err = readContainerOutputs(path)
		Expect(err).To(MatchError(ContainSubstring("1 MiB")))
		Expect(os.Remove(path)).To(Succeed())
		_, err = readContainerOutputs(path)
		Expect(err).To(HaveOccurred())
		Expect(os.Symlink("missing", path)).To(Succeed())
		_, err = readContainerOutputs(path)
		Expect(err).To(HaveOccurred())
		Expect(os.Remove(path)).To(Succeed())
		Expect(os.Mkdir(path, 0700)).To(Succeed())
		_, err = readContainerOutputs(path)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("output destination intent", func() {
	It("protects renamed output sources without changing declarations", func() {
		secret := &model.ParamSpec{String: &model.StringParamSpec{IsSecret: true}}
		public := &model.ParamSpec{String: &model.StringParamSpec{}}
		ctx := withProtectedSources(context.Background(), map[string]*model.ParamSpec{"destination": secret})
		params := map[string]*model.ParamSpec{"source": public, "declared": secret}
		result := outputDestinations(ctx, params, map[string]string{"source": "$(destination)"}, nil)
		Expect(secretParam(result["source"])).To(BeTrue())
		Expect(secretParam(result["declared"])).To(BeTrue())
		Expect(secretParam(params["source"])).To(BeFalse())
		value := "private"
		result = outputDestinations(context.Background(), params, map[string]string{"source": "$(incoming)"}, map[string]*model.Value{"incoming": {String: &value, Sensitive: true}})
		Expect(secretParam(result["source"])).To(BeTrue())
	})
})
