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
      objectUnderTest.SetDevOpName(expectedDevOpRunView.devOpName)
      objectUnderTest.SetStartedAtPosixTime(expectedDevOpRunView.StartedAtPosixTime())
      objectUnderTest.SetEndedAtPosixTime(expectedDevOpRunView.EndedAtPosixTime())
      objectUnderTest.SetExitCode(expectedDevOpRunView.ExitCode())

      /* act */
      actualDevOpRunView := objectUnderTest.Build()

      /* assert */
      Expect(actualDevOpRunView).To(Equal(expectedDevOpRunView))

    })
  })
})
