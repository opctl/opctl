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
    It("should roundtrip an opBundleManifest", func() {

      /* arrange */
      expectedOpBundleManifest := models.OpBundleManifest{
        BundleManifest:models.BundleManifest{
          Name:"dummyName",
          Description:"dummyDescription",
          Version:"dummyVersion",
        },
        Inputs:[]models.Param{
          {
            String: &models.StringParam{
              Name:"dummyName",
              Default:"dummyDefault",
              Description:"dummyDescription",
              IsSecret:false,
            },
          },
        },
        Run:&models.RunStatement{Op:"dummyOpRef"},
      }

      objectUnderTest := _yamlCodec{}

      /* act */
      opBundleManifestBytes, _ := objectUnderTest.ToYaml(&expectedOpBundleManifest)
      actualOpBundleManifest := models.OpBundleManifest{}
      objectUnderTest.FromYaml(opBundleManifestBytes, &actualOpBundleManifest)

      /* assert */
      Expect(actualOpBundleManifest).To(Equal(expectedOpBundleManifest))

    })
  })
})
