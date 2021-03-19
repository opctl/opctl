package data

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/fs"
	aggregateError "github.com/opctl/opctl/sdks/go/internal/aggregate_error"
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
			var expected aggregateError.ErrAggregate
			expected.AddError(errors.Wrap(fmt.Errorf("skipped"), provider0.Label()))
			Expect(actualErr).To(MatchError(errors.Wrap(expected, "unable to resolve op '\\not/exist'").Error()))
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
