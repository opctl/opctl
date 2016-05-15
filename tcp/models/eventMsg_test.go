package models

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/dev-op-spec/engine/core/models"
  "time"
)

var _ = Describe("EventMsg", func() {
  Context("NewEventMsg(:Event)", func() {
    It("Should use provided events reflected type name and strip 'Event' suffix for Type attribute", func() {

      /* arrange */
      providedEvent := models.NewLogEntryEmittedEvent(
        *new(string),
        *new(time.Time),
        *new(string),
        *new(string),
      )

      expectedType := "NewLogEntryEmitted"

      /* act */
      resultFromNewEventMsg := NewEventMsg(providedEvent)

      /* assert */
      Expect(resultFromNewEventMsg.Type).To(Equal(expectedType))

    })
  })
})
