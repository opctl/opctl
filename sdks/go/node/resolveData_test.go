package node

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("core", func() {
	Context("ResolveData", func() {
		It("should call data.Resolve w/ expected args", func() {
			/* arrange */
			dataCachePath, err := os.MkdirTemp("", "")
			if err != nil {
				panic(err)
			}

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
				&model.Creds{},
			)

			/* assert */
			Expect(actualErr).To(BeNil())
			Expect(actualOp.Ref()).To(Equal(providedOpRef))
		})
	})
})
