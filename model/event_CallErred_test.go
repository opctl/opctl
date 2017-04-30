package model

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("CallErredEvent", func() {
	Context("Container", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallErred: &CallErredEvent{
				CallErredEventBase: &CallErredEventBase{
					Msg: "dummyMsg",
					CallEventBase: &CallEventBase{
						CallID:     "dummyCallID",
						RootCallID: "dummyRootCallID",
					},
				},
				Container: &ContainerCallErredEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ContainerCallErred",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallErred.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallErred.CallID),
						fmt.Sprintf("Msg='%v'", providedEvent.CallErred.Msg),
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
			CallErred: &CallErredEvent{
				CallErredEventBase: &CallErredEventBase{
					Msg: "dummyMsg",
					CallEventBase: &CallEventBase{
						CallID:     "dummyCallID",
						RootCallID: "dummyRootCallID",
					},
				},
				Op: &OpCallErredEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"OpCallErred",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallErred.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallErred.CallID),
						fmt.Sprintf("Msg='%v'", providedEvent.CallErred.Msg),
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
			CallErred: &CallErredEvent{
				CallErredEventBase: &CallErredEventBase{
					Msg: "dummyMsg",
					CallEventBase: &CallEventBase{
						CallID:     "dummyCallID",
						RootCallID: "dummyRootCallID",
					},
				},
				Parallel: &ParallelCallErredEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ParallelCallErred",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallErred.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallErred.CallID),
						fmt.Sprintf("Msg='%v'", providedEvent.CallErred.Msg),
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
			CallErred: &CallErredEvent{
				CallErredEventBase: &CallErredEventBase{
					Msg: "dummyMsg",
					CallEventBase: &CallEventBase{
						CallID:     "dummyCallID",
						RootCallID: "dummyRootCallID",
					},
				},
				Serial: &SerialCallErredEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"SerialCallErred",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallErred.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallErred.CallID),
						fmt.Sprintf("Msg='%v'", providedEvent.CallErred.Msg),
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
