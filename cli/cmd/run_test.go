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
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
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

	It("describes Go inputs through the node and passes them to the satisfier", func() {
		dir := handle(map[string]string{"op.go": "package main"})
		dir.RefReturns("/ops/go")
		api := new(nodeFakes.FakeNode)
		api.DescribeOpReturns(&model.OpSpec{Inputs: map[string]*model.ParamSpec{"env": {String: &model.StringParamSpec{Default: "dev"}}}}, nil)
		spec, _, err := getRunOpFile(context.Background(), dir, api)
		Expect(err).NotTo(HaveOccurred())
		_, req := api.DescribeOpArgsForCall(0)
		Expect(req.Op.Ref).To(Equal("/ops/go"))
		satisfier := cliparamsatisfier.New(nil)
		args, err := satisfier.Satisfy(cliparamsatisfier.NewInputSourcer(satisfier.NewSliceInputSrc([]string{"env=prod"}, "=")), spec.Inputs)
		Expect(err).NotTo(HaveOccurred())
		Expect(*args["env"].String).To(Equal("prod"))
		failure := errors.New("invalid Go interface")
		api.DescribeOpReturns(nil, failure)
		_, _, err = getRunOpFile(context.Background(), dir, api)
		Expect(err).To(Equal(failure))
	})

	It("retains the YAML op's inputs and defaults", func() {
		spec, _, err := getRunOpFile(context.Background(), handle(map[string]string{"op.yml": "name: yaml\ninputs:\n  greeting:\n    string: {default: hello}\n"}), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(spec.Inputs["greeting"].String.Default).To(Equal("hello"))
	})
	It("rejects both declarations before considering YAML validity", func() {
		_, _, err := getRunOpFile(context.Background(), handle(map[string]string{"op.go": "package main", "op.yml": "invalid"}), nil)
		Expect(err).To(Equal(opfile.ErrBothOpFiles))
	})
	It("keeps the original YAML read error when no declaration is readable", func() {
		dir := new(fakes.FakeDataHandle)
		failure := errors.New("remote error")
		dir.GetContentReturns(nil, failure)
		_, _, err := getRunOpFile(context.Background(), dir, nil)
		Expect(err).To(Equal(failure))
	})
})
