package core

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("yamlCodec", func() {
  Context("executing .toYaml() then .fromYaml", func() {
    It("should roundtrip a devOpFile", func() {

      /* arrange */
      expectedDevOpFile := devOpFile{Description:"dev op description"}
      objectUnderTest := yamlCodecImpl{}

      /* act */
      devOpFileBytes, _ := objectUnderTest.toYaml(&expectedDevOpFile)
      actualDevOpFile := devOpFile{}
      objectUnderTest.fromYaml(devOpFileBytes, &actualDevOpFile)

      /* assert */
      Expect(actualDevOpFile).To(Equal(expectedDevOpFile))

    })
    It("should roundtrip a pipelineFile", func() {

      /* arrange */
      expectedPipelineFile := pipelineFile{
        Description:"pipeline description",
        Stages:[]pipelineFileStage{
          pipelineFileStage{Name:"pipeline stage name", Type:pipelineStageType},
          pipelineFileStage{Name:"dev op stage name", Type:devOpStageType},
        },
      }
      objectUnderTest := yamlCodecImpl{}

      /* act */
      pipelineFileBytes, _ := objectUnderTest.toYaml(&expectedPipelineFile)
      actualPipelineFile := pipelineFile{}
      objectUnderTest.fromYaml(pipelineFileBytes, &actualPipelineFile)

      /* assert */
      Expect(actualPipelineFile).To(Equal(expectedPipelineFile))

    })
  })
})
