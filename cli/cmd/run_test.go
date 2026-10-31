package cmd

import (
	"context"
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/model/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

var _ = Describe("load an op for run", func() {
	handle := func(files map[string]string) *fakes.FakeDataHandle {
		dir := new(fakes.FakeDataHandle)
		dir.GetContentStub = func(_ context.Context, path string) (model.ReadSeekCloser, error) {
			text, ok := files[path]
			if !ok {
				return nil, errors.New("HTTP 404")
			}
			return struct {
				*strings.Reader
				io.Closer
			}{strings.NewReader(text), io.NopCloser(strings.NewReader(""))}, nil
		}
		return dir
	}
	It("runs a Go op without prompting and ignores unknown arguments", func() {
		spec, err := getRunOpFile(context.Background(), handle(map[string]string{"op.go": "package main"}))
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Inputs).To(BeEmpty())
		satisfier := cliparamsatisfier.New(nil)
		args, err := satisfier.Satisfy(cliparamsatisfier.NewInputSourcer(satisfier.NewSliceInputSrc([]string{"ignored=value"}, "="), satisfier.NewCliPromptInputSrc(spec.Inputs)), spec.Inputs)
		Expect(err).NotTo(HaveOccurred())
		Expect(args).To(BeEmpty())
	})
	It("retains the YAML op's inputs and defaults", func() {
		spec, err := getRunOpFile(context.Background(), handle(map[string]string{"op.yml": "name: yaml\ninputs:\n  greeting:\n    string: {default: hello}\n"}))
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Inputs["greeting"].String.Default).To(Equal("hello"))
	})
	It("rejects both declarations before considering YAML validity", func() {
		_, err := getRunOpFile(context.Background(), handle(map[string]string{"op.go": "package main", "op.yml": "invalid"}))
		Expect(err).To(Equal(opfile.ErrBothOpFiles))
	})
	It("keeps the original YAML read error when no declaration is readable", func() {
		dir := new(fakes.FakeDataHandle)
		failure := errors.New("remote error")
		dir.GetContentReturns(nil, failure)
		_, err := getRunOpFile(context.Background(), dir)
		Expect(err).To(Equal(failure))
	})
})
