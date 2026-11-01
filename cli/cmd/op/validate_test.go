package op

import (
	"context"
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/model/fakes"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

var _ = Describe("validate an op declaration", func() {
	It("describes Go ops, rejects ambiguous declarations, and accepts YAML ops", func() {
		for _, tc := range []struct {
			files map[string]string
			want  error
		}{
			{map[string]string{"op.go": "package main"}, nil},
			{map[string]string{"op.go": "package main", "op.yml": "name: both"}, opfile.ErrBothOpFiles},
			{map[string]string{"op.yml": "name: yaml"}, nil},
		} {
			dir := new(fakes.FakeDataHandle)
			dir.GetContentStub = func(_ context.Context, path string) (model.ReadSeekCloser, error) {
				text, ok := tc.files[path]
				if !ok {
					return nil, errors.New("HTTP 404")
				}
				return struct {
					*strings.Reader
					io.Closer
				}{strings.NewReader(text), io.NopCloser(strings.NewReader(""))}, nil
			}
			api := new(nodeFakes.FakeNode)
			api.DescribeOpReturns(&model.OpSpec{Name: "go"}, nil)
			err := validateOpFile(context.Background(), dir, api, "")
			if len(tc.files) == 1 && tc.files["op.go"] != "" {
				Expect(api.DescribeOpCallCount()).To(Equal(1))
				failure := errors.New("unknown opctl tag key minimun")
				api.DescribeOpReturns(nil, failure)
				Expect(validateOpFile(context.Background(), dir, api, "")).To(Equal(failure))
			} else {
				Expect(api.DescribeOpCallCount()).To(Equal(0))
			}
			if tc.want == nil {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(tc.want))
			}
		}
	})
})
