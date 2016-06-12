package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opctl/engine/core/models"
)

var _ = Describe("_yamlCodec", func() {
  Context("newYamlCodec()", func() {
    It("should return an instance of _yamlCodec", func() {

      /* arrange/act */
      objectUnderTest := newYamlCodec()

      /* assert */
      Expect(objectUnderTest).To(BeAssignableToTypeOf(&_yamlCodec{}))

    })
  })
  Context("executing .toYaml() then .fromYaml", func() {
    It("should roundtrip a opFile", func() {

      /* arrange */
      expectedOpFile := models.OpFile{Description:"op description"}
      objectUnderTest := _yamlCodec{}

      /* act */
      opFileBytes, _ := objectUnderTest.ToYaml(&expectedOpFile)
      actualOpFile := models.OpFile{}
      objectUnderTest.FromYaml(opFileBytes, &actualOpFile)

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
      opFileBytes, _ := objectUnderTest.ToYaml(&expectedOpFile)
      actualOpFile := models.OpFile{}
      objectUnderTest.FromYaml(opFileBytes, &actualOpFile)

      /* assert */
      Expect(actualOpFile).To(Equal(expectedOpFile))

    })
  })
})
