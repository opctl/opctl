package os

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "path"
)

var _ = Describe("relPathToPipelineDirFactory", func() {
  Context("executing .Construct", func() {
    It("should return expected path", func() {

      /* arrange */
      providedPipelineName := "providedPipelineName"
      objectUnderTest := newRelPathToPipelineDirFactory()
      expectedRelPathToPipelineDir := path.Join(relPathToPipelinesDir, providedPipelineName)

      /* act */
      actualRelPathToPipelineDir, _ := objectUnderTest.Construct(providedPipelineName)

      /* assert */
      Expect(actualRelPathToPipelineDir).To(Equal(expectedRelPathToPipelineDir))

    })
  })
})
