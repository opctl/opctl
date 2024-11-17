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
			expected.AddError(fmt.Errorf("%s: %w", provider0.Label(), model.ErrDataSkipped{}))
			Expect(actualErr).To(MatchError(fmt.Errorf("%w op \"%s\": %w", model.ErrDataUnableToResolve{}, dataRef, expected)))
		})
	})
	Context("providers[0].TryResolve doesn't err", func() {
		It("should return expected results", func() {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			opRef := filepath.Join(wd, "testdata/testop")
			provider0 := fs.New()

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
