package models

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "time"
)

var _ = Describe("devOpRunViewBuilder", func() {
  Context("executing .Build", func() {
    It("should return expected DevOpRunView", func() {

      /* arrange */
      expectedDevOpRunView := newDevOpRunView("devOpName", time.Now().Unix(), time.Now().Unix(), 1)

      objectUnderTest := NewDevOpRunViewBuilder()
      objectUnderTest.SetDevOpName(expectedDevOpRunView.DevOpName())
      objectUnderTest.SetStartedAtEpochTime(expectedDevOpRunView.StartedAtEpochTime())
      objectUnderTest.SetEndedAtEpochTime(expectedDevOpRunView.EndedAtEpochTime())
      objectUnderTest.SetExitCode(expectedDevOpRunView.ExitCode())

      /* act */
      actualDevOpRunView := objectUnderTest.Build()

      /* assert */
      Expect(actualDevOpRunView).To(Equal(expectedDevOpRunView))

    })
  })
})
