package dockercompose

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("yamlCodec", func() {
  Context("executing .toYaml() then .fromYaml", func() {
    It("should roundtrip a dockerComposeFile", func() {

      /* arrange */
      expectedDockerComposeFile := dockerComposeFile{
        Version: "2",
        Services:map[string]dockerComposeFileService{
          "dev-op-name":dockerComposeFileService{
            Image:"alpine:3.3",
          },
        },
      }
      objectUnderTest := _yamlCodec{}

      /* act */
      dockercomposeFileBytes, _ := objectUnderTest.toYaml(&expectedDockerComposeFile)
      actualDockerComposeFile := dockerComposeFile{}
      objectUnderTest.fromYaml(dockercomposeFileBytes, &actualDockerComposeFile)

      /* assert */
      Expect(actualDockerComposeFile).To(Equal(expectedDockerComposeFile))

    })
  })
})
