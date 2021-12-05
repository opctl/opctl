package node

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("GetDataReq", func() {
	Context("req.DataRef empty", func() {
		It("should return expected result", func() {
			/* arrange */
			objectUnderTest := core{}

			/* act */
			actualData, actualErr := objectUnderTest.GetData(
				context.Background(),
				model.GetDataReq{},
			)

			/* assert */
			Expect(actualData).To(BeNil())
			Expect(actualErr).To(MatchError(`"" not a valid data ref`))
		})
	})
})
