package data

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Resolve", func() {
	Context("providers[0].TryResolve errs", func() {
		It("should return error", func() {
			/* arrange */
			provider0 := fs.New()

			providedProviders := []model.DataProvider{provider0}

			/* act */
			_, actualErr := Resolve(
				context.Background(),
				"\\not/exist",
				providedProviders...,
			)

			/* assert */
			Expect(actualErr).To(Equal(model.ErrDataRefResolution{}))
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
