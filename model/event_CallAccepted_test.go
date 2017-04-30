package model

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("CallAcceptedEvent", func() {
	Context("Container", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallAccepted: &CallAcceptedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Container: &ContainerCallAcceptedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ContainerCallAccepted",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallAccepted.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallAccepted.CallID),
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
			CallAccepted: &CallAcceptedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Op: &OpCallAcceptedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"OpCallAccepted",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallAccepted.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallAccepted.CallID),
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
			CallAccepted: &CallAcceptedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Parallel: &ParallelCallAcceptedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ParallelCallAccepted",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallAccepted.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallAccepted.CallID),
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
			CallAccepted: &CallAcceptedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Serial: &SerialCallAcceptedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"SerialCallAccepted",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallAccepted.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallAccepted.CallID),
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
