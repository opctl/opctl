package core

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("core", func() {
	Context("ResolveData", func() {
		It("should call data.Resolve w/ expected args", func() {
			/* arrange */
			dataCachePath := os.TempDir()

			providedCtx := context.Background()
			// some public repo that's relatively small
			providedOpRef := "github.com/opspec-pkgs/_.op.create#3.3.1"

			objectUnderTest := core{
				dataCachePath: dataCachePath,
			}

			/* act */
			actualOp, actualErr := objectUnderTest.ResolveData(
				providedCtx,
				providedOpRef,
				nil,
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(*actualOp.Path()).To(Equal(filepath.Join(dataCachePath, providedOpRef)))
		})
	})
})
