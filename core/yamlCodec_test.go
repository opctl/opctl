package core

var _ = Describe("yamlCodec", func() {
  Context("executing .toYaml() then .fromYaml", func() {
    It("should roundtrip a opFile", func() {

      /* arrange */
      expectedOpFile := opFile{Description:"op description"}
      objectUnderTest := _yamlCodec{}

      /* act */
      opFileBytes, _ := objectUnderTest.toYaml(&expectedOpFile)
      actualOpFile := opFile{}
      objectUnderTest.fromYaml(opFileBytes, &actualOpFile)

      /* assert */
      Expect(actualOpFile).To(Equal(expectedOpFile))

    })
    It("should roundtrip a opFile", func() {

      /* arrange */
      expectedOpFile := opFile{
        Description:"op description",
        SubOps:[]opFileSubOp{
          opFileSubOp{Url:"op1-name"},
          opFileSubOp{Url:"op2-name"},
        },
      }
      objectUnderTest := _yamlCodec{}

      /* act */
      opFileBytes, _ := objectUnderTest.toYaml(&expectedOpFile)
      actualOpFile := opFile{}
      objectUnderTest.fromYaml(opFileBytes, &actualOpFile)

      /* assert */
      Expect(actualOpFile).To(Equal(expectedOpFile))

    })
  })
})
