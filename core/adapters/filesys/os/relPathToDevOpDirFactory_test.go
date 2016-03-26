package os

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "path"
)

var _ = Describe("relPathToDevOpDirFactory", func() {
  Context("executing .Construct", func() {
    It("should return expected path", func() {

      /* arrange */
      providedDevOpName := "providedDevOpName"
      objectUnderTest := newRelPathToDevOpDirFactory()
      expectedRelPathToDevOpDir := path.Join(relPathToDevOpsDir, providedDevOpName)

      /* act */
      actualRelPathToDevOpDir, _ := objectUnderTest.Construct(providedDevOpName)

      /* assert */
      Expect(actualRelPathToDevOpDir).To(Equal(expectedRelPathToDevOpDir))

    })
  })
})
