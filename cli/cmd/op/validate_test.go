package op

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/model/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"io"
	"strings"
)

var _ = Describe("validate an op declaration", func() {
	It("rejects Go ops explicitly, both declarations as ambiguous, and accepts YAML ops", func() {
		for _, tc := range []struct {
			files map[string]string
			want  error
		}{
			{map[string]string{"op.go": "package main"}, errors.New("opctl op validate does not support Go ops yet")},
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
			err := validateOpFile(context.Background(), dir)
			if tc.want == nil {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(tc.want))
			}
		}
	})
})
