package model

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("CallRequestedEvent", func() {
	Context("Container", func() {
		providedEvent := Event{
			Timestamp: time.Now().UTC(),
			CallRequested: &CallRequestedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Container: &ContainerCallRequestedEvent{
					ImageRef: "dummyImageRef",
				},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ContainerCallRequested",
						fmt.Sprintf("ImageRef='%v'", providedEvent.CallRequested.Container.ImageRef),
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallRequested.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallRequested.CallID),
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
			CallRequested: &CallRequestedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Op: &OpCallRequestedEvent{
					PkgRef: "dummyPkgRef",
				},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"OpCallRequested",
						fmt.Sprintf("PkgRef='%v'", providedEvent.CallRequested.Op.PkgRef),
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallRequested.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallRequested.CallID),
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
			CallRequested: &CallRequestedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Parallel: &ParallelCallRequestedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"ParallelCallRequested",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallRequested.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallRequested.CallID),
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
			CallRequested: &CallRequestedEvent{
				CallEventBase: &CallEventBase{
					CallID:     "dummyCallID",
					RootCallID: "dummyRootCallID",
				},
				Serial: &SerialCallRequestedEvent{},
			},
		}
		Context("text format", func() {
			It("should match expected text", func() {
				/* arrange */
				expectedText := strings.Join(
					[]string{
						"SerialCallRequested",
						fmt.Sprintf("RootCallId='%v'", providedEvent.CallRequested.RootCallID),
						fmt.Sprintf("CallId='%v'", providedEvent.CallRequested.CallID),
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
