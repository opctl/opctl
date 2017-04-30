package model

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("CallEndedEvent", func() {
	Context("Container", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallEnded: &CallEndedEvent{
				CallEndedEventBase: &CallEndedEventBase{
					Outcome: CallOutcomeCancelled,
					CallEventBase: &CallEventBase{
						CallID:     "dummyCallID",
						RootCallID: "dummyRootCallID",
					},
				},
				Container: &ContainerCallEndedEvent{
					ExitCode: 1,
				},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ContainerCallEnded",
						fmt.Sprintf("ExitCode='%v'", providedEvent.CallEnded.Container.ExitCode),
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallEnded.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallEnded.CallID),
						fmt.Sprintf("Outcome='%v'", providedEvent.CallEnded.Outcome),
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
			CallEnded: &CallEndedEvent{
				CallEndedEventBase: &CallEndedEventBase{
					Outcome: CallOutcomeCancelled,
					CallEventBase: &CallEventBase{
						CallID:     "dummyCallID",
						RootCallID: "dummyRootCallID",
					},
				},
				Op: &OpCallEndedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"OpCallEnded",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallEnded.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallEnded.CallID),
						fmt.Sprintf("Outcome='%v'", providedEvent.CallEnded.Outcome),
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
			CallEnded: &CallEndedEvent{
				CallEndedEventBase: &CallEndedEventBase{
					Outcome: CallOutcomeCancelled,
					CallEventBase: &CallEventBase{
						CallID:     "dummyCallID",
						RootCallID: "dummyRootCallID",
					},
				},
				Parallel: &ParallelCallEndedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ParallelCallEnded",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallEnded.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallEnded.CallID),
						fmt.Sprintf("Outcome='%v'", providedEvent.CallEnded.Outcome),
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
			CallEnded: &CallEndedEvent{
				CallEndedEventBase: &CallEndedEventBase{
					Outcome: CallOutcomeCancelled,
					CallEventBase: &CallEventBase{
						CallID:     "dummyCallID",
						RootCallID: "dummyRootCallID",
					},
				},
				Serial: &SerialCallEndedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"SerialCallEnded",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallEnded.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallEnded.CallID),
						fmt.Sprintf("Outcome='%v'", providedEvent.CallEnded.Outcome),
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
