package sdk_golang

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opctl/engine/core/models"
)

var _ = Describe("yamlCodec", func() {
  Context("executing .toYaml() then .fromYaml", func() {
    It("should roundtrip a opFile", func() {

      /* arrange */
      expectedOpFile := models.OpFile{Description:"op description"}
      objectUnderTest := _yamlCodec{}

      /* act */
      opFileBytes, _ := objectUnderTest.toYaml(&expectedOpFile)
      actualOpFile := models.OpFile{}
      objectUnderTest.fromYaml(opFileBytes, &actualOpFile)

      /* assert */
      Expect(actualOpFile).To(Equal(expectedOpFile))

    })
    It("should roundtrip a opFile", func() {

      /* arrange */
      expectedOpFile := models.OpFile{
        Description:"op description",
        SubOps:[]models.OpFileSubOp{
          models.OpFileSubOp{Url:"op1-name"},
          models.OpFileSubOp{Url:"op2-name"},
        },
      }
      objectUnderTest := _yamlCodec{}

      /* act */
      opFileBytes, _ := objectUnderTest.toYaml(&expectedOpFile)
      actualOpFile := models.OpFile{}
      objectUnderTest.fromYaml(opFileBytes, &actualOpFile)

      /* assert */
      Expect(actualOpFile).To(Equal(expectedOpFile))

    })
  })
})
