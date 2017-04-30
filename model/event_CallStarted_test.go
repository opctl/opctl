package model

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("CallStartedEvent", func() {
	Context("Container", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallStarted: &CallStartedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Container: &ContainerCallStartedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ContainerCallStarted",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallStarted.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallStarted.CallID),
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
			CallStarted: &CallStartedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Op: &OpCallStartedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"OpCallStarted",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallStarted.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallStarted.CallID),
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
			CallStarted: &CallStartedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Parallel: &ParallelCallStartedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ParallelCallStarted",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallStarted.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallStarted.CallID),
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
			CallStarted: &CallStartedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Serial: &SerialCallStartedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"SerialCallStarted",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallStarted.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallStarted.CallID),
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
