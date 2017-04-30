package model

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("CallCancelledEvent", func() {
	Context("Container", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallCancelled: &CallCancelledEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Container: &ContainerCallCancelledEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ContainerCallCancelled",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallCancelled.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallCancelled.CallID),
						fmt.Sprintf("Timestamp='%v'", providedEvent.Timestamp.Format(time.RFC3339)),
					},
					" ",
				)

				/* act */
				actualText := fmt.Sprint(providedEvent)

				/* assert */
				Expect(actualText).To(Equal(expectedText))

			})
		})
	})
	Context("Op", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallCancelled: &CallCancelledEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Op: &OpCallCancelledEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"OpCallCancelled",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallCancelled.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallCancelled.CallID),
						fmt.Sprintf("Timestamp='%v'", providedEvent.Timestamp.Format(time.RFC3339)),
					},
					" ",
				)

				/* act */
				actualText := fmt.Sprint(providedEvent)

				/* assert */
				Expect(actualText).To(Equal(expectedText))

			})
		})
	})
	Context("Parallel", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallCancelled: &CallCancelledEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Parallel: &ParallelCallCancelledEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ParallelCallCancelled",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallCancelled.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallCancelled.CallID),
						fmt.Sprintf("Timestamp='%v'", providedEvent.Timestamp.Format(time.RFC3339)),
					},
					" ",
				)

				/* act */
				actualText := fmt.Sprint(providedEvent)

				/* assert */
				Expect(actualText).To(Equal(expectedText))

			})
		})
	})
	Context("Serial", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallCancelled: &CallCancelledEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Serial: &SerialCallCancelledEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"SerialCallCancelled",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallCancelled.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallCancelled.CallID),
						fmt.Sprintf("Timestamp='%v'", providedEvent.Timestamp.Format(time.RFC3339)),
					},
					" ",
				)

				/* act */
				actualText := fmt.Sprint(providedEvent)

				/* assert */
				Expect(actualText).To(Equal(expectedText))

			})
		})
	})
})
