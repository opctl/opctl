package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("yamlCodec", func() {
  Context("executing .toYaml() then .fromYaml", func() {
    It("should roundtrip a operationFile", func() {

      /* arrange */
      expectedOperationFile := operationFile{Description:"operation description"}
      objectUnderTest := _yamlCodec{}

      /* act */
      operationFileBytes, _ := objectUnderTest.toYaml(&expectedOperationFile)
      actualOperationFile := operationFile{}
      objectUnderTest.fromYaml(operationFileBytes, &actualOperationFile)

      /* assert */
      Expect(actualOperationFile).To(Equal(expectedOperationFile))

    })
    It("should roundtrip a operationFile", func() {

      /* arrange */
      expectedOperationFile := operationFile{
        Description:"operation description",
        SubOperations:[]operationFileSubOperation{
          operationFileSubOperation{Name:"subOperation name"},
          operationFileSubOperation{Name:"subOperation name"},
        },
      }
      objectUnderTest := _yamlCodec{}

      /* act */
      operationFileBytes, _ := objectUnderTest.toYaml(&expectedOperationFile)
      actualOperationFile := operationFile{}
      objectUnderTest.fromYaml(operationFileBytes, &actualOperationFile)

      /* assert */
      Expect(actualOperationFile).To(Equal(expectedOperationFile))

    })
  })
})
