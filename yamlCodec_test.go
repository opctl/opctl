package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
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
    It("should roundtrip an opManifest", func() {

      /* arrange */
      expectedOpManifest := models.OpManifest{
        Manifest:models.Manifest{
          Name:"dummyName",
          Description:"dummyDescription",
          Version:"dummyVersion",
        },
        Inputs:[]models.Param{
          {
            Name:"dummyName",
            Description:"dummyDescription",
            IsSecret:false,
            String: &models.StringParam{
              Default:"dummyDefault",
            },
          },
        },
        Run:&models.RunDeclaration{Op:"dummyOpRef"},
      }

      objectUnderTest := _yamlCodec{}

      /* act */
      opManifestBytes, _ := objectUnderTest.ToYaml(&expectedOpManifest)
      actualOpManifest := models.OpManifest{}
      objectUnderTest.FromYaml(opManifestBytes, &actualOpManifest)

      /* assert */
      Expect(actualOpManifest).To(Equal(expectedOpManifest))

    })
  })
})
