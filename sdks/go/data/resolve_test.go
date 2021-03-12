package data

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

var _ = Context("Resolve", func() {
	Context("providers[0].TryResolve errs", func() {
		It("should return error", func() {
			/* arrange */
			provider0 := fs.New()
			providedProviders := []model.DataProvider{provider0}
			dataRef := "\\not/exist"

			/* act */
			_, actualErr := Resolve(
				context.Background(),
				dataRef,
				providedProviders...,
			)

			/* assert */
			Expect(actualErr).To(MatchError(ErrDataResolution{
				dataRef: dataRef,
				errs:    []error{errors.Wrap(errors.New("not found"), provider0.Label())},
			}.Error()))
		})
	})
	Context("providers[0].TryResolve doesn't err", func() {
		It("should return expected results", func() {
			wd, err := os.Getwd()
			if nil != err {
				panic(err)
			}
			opRef := filepath.Join(wd, "testdata/testop")
			provider0 := fs.New(filepath.Dir(opRef))

			providedProviders := []model.DataProvider{provider0}

			/* act */
			actualHandle, actualErr := Resolve(
				context.Background(),
				opRef,
				providedProviders...,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(actualHandle.Ref()).To(Equal(opRef))
		})
	})
})
