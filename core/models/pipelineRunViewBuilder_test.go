package models

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "time"
)

var _ = Describe("pipelineRunViewBuilder", func() {
  Context("executing .Build", func() {
    It("should return expected PipelineRunView", func() {

      /* arrange */
      expectedPipelineRunView := newPipelineRunView("pipelineName", nil, time.Now().Unix(), time.Now().Unix(), 1)

      objectUnderTest := NewPipelineRunViewBuilder()
      objectUnderTest.SetPipelineName(expectedPipelineRunView.PipelineName())
      objectUnderTest.SetStartedAtEpochTime(expectedPipelineRunView.StartedAtEpochTime())
      objectUnderTest.SetEndedAtEpochTime(expectedPipelineRunView.EndedAtEpochTime())
      objectUnderTest.SetExitCode(expectedPipelineRunView.ExitCode())

      /* act */
      actualPipelineRunView := objectUnderTest.Build()

      /* assert */
      Expect(actualPipelineRunView).To(Equal(expectedPipelineRunView))

    })
  })
})
