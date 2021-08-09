package core

import (
	"context"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("ListDescendants", func() {
	Context("req.DataRef empty", func() {
		It("should return expected result", func() {
			/* arrange */
			objectUnderTest := core{}

			/* act */
			actualDescendants, actualErr := objectUnderTest.ListDescendants(
				context.Background(),
				model.ListDescendantsReq{},
			)

			/* assert */
			Expect(actualDescendants).To(BeEmpty())
			Expect(actualErr).To(MatchError(`"" not a valid data ref`))
		})
	})
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	Context("req.DataRef absolute path", func() {
		It("should return expected result", func() {
			/* arrange */
			providedDataRef := path.Join(wd, "testdata/listDescendants")
			objectUnderTest := core{}

			expectedDescendants := []*model.DirEntry{
				{Path: "/empty.txt", Size: 0, Mode: 420},
			}

			/* act */
			actualDescendants, actualErr := objectUnderTest.ListDescendants(
				context.Background(),
				model.ListDescendantsReq{
					DataRef: providedDataRef,
				},
			)

			/* assert */
			Expect(actualDescendants).To(ConsistOf(expectedDescendants))
			Expect(actualErr).To(BeNil())
		})
	})
})
