package model

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/util/format"
	"time"
)

var _ = Describe("Event", func() {
	json := format.NewJsonFormat()

	Context("when formatting to/from json", func() {

		Context("with non-nil $.containerExited", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedEvent := Event{
					ContainerExited: &ContainerExitedEvent{
						ContainerRef: "dummyContainerRef",
						ExitCode:     1,
						RootOpId:     "dummyRootOpId",
						ContainerId:  "dummyContainerId",
						OpPkgRef:     "dummyOpPkgRef",
					},
					Timestamp: time.Now().UTC(),
				}

				/* act */
				providedJson, err := json.From(expectedEvent)
				if nil != err {
					panic(err)
				}

				actualEvent := Event{}
				json.To(providedJson, &actualEvent)

				/* assert */
				Expect(actualEvent).To(Equal(expectedEvent))

			})

		})

		Context("with non-nil $.containerStarted", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedEvent := Event{
					ContainerStarted: &ContainerStartedEvent{
						ContainerRef: "dummyContainerRef",
						OpPkgRef:     "dummyOpPkgRef",
						ContainerId:  "dummyContainerId",
						RootOpId:     "dummyRootOpId",
					},
					Timestamp: time.Now().UTC(),
				}

				/* act */
				providedJson, err := json.From(expectedEvent)
				if nil != err {
					panic(err)
				}

				actualEvent := Event{}
				json.To(providedJson, &actualEvent)

				/* assert */
				Expect(actualEvent).To(Equal(expectedEvent))

			})

		})

		Context("with non-nil $.containerStdErrWrittenTo", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedEvent := Event{
					ContainerStdErrWrittenTo: &ContainerStdErrWrittenToEvent{
						Data:        []byte("dummyData"),
						OpPkgRef:    "dummyOpPkgRef",
						ContainerId: "dummyContainerId",
						RootOpId:    "dummyRootOpId",
					},
					Timestamp: time.Now().UTC(),
				}

				/* act */
				providedJson, err := json.From(expectedEvent)
				if nil != err {
					panic(err)
				}

				actualEvent := Event{}
				json.To(providedJson, &actualEvent)

				/* assert */
				Expect(actualEvent).To(Equal(expectedEvent))

			})

		})

		Context("with non-nil $.containerStdOutWrittenTo", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedEvent := Event{
					ContainerStdOutWrittenTo: &ContainerStdOutWrittenToEvent{
						Data:        []byte("dummyData"),
						OpPkgRef:    "dummyOpPkgRef",
						ContainerId: "dummyContainerId",
						RootOpId:    "dummyRootOpId",
					},
					Timestamp: time.Now().UTC(),
				}

				/* act */
				providedJson, out := json.From(expectedEvent)
				if nil != out {
					panic(out)
				}

				actualEvent := Event{}
				json.To(providedJson, &actualEvent)

				/* assert */
				Expect(actualEvent).To(Equal(expectedEvent))

			})

		})

		Context("with non-nil $.opEnded", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedEvent := Event{
					OpEnded: &OpEndedEvent{
						OpPkgRef: "dummyOpPkgRef",
						OpId:     "dummyOpId",
						Outcome:  "dummyOutcome",
						RootOpId: "dummyRootOpId",
					},
					Timestamp: time.Now().UTC(),
				}

				/* act */
				providedJson, err := json.From(expectedEvent)
				if nil != err {
					panic(err)
				}

				actualEvent := Event{}
				json.To(providedJson, &actualEvent)

				/* assert */
				Expect(actualEvent).To(Equal(expectedEvent))

			})

		})

		Context("with non-nil $.opStarted", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedEvent := Event{
					OpStarted: &OpStartedEvent{
						OpPkgRef: "dummyOpPkgRef",
						OpId:     "dummyOpId",
						RootOpId: "dummyRootOpId",
					},
					Timestamp: time.Now().UTC(),
				}

				/* act */
				providedJson, err := json.From(expectedEvent)
				if nil != err {
					panic(err)
				}

				actualEvent := Event{}
				json.To(providedJson, &actualEvent)

				/* assert */
				Expect(actualEvent).To(Equal(expectedEvent))

			})

		})

		Context("with non-nil $.opEncounteredError", func() {

			It("should have expected attributes", func() {

				/* arrange */
				expectedEvent := Event{
					OpEncounteredError: &OpEncounteredErrorEvent{
						Msg:      "dummyMsg",
						OpPkgRef: "dummyOpPkgRef",
						OpId:     "dummyOpId",
						RootOpId: "dummyRootOpId",
					},
					Timestamp: time.Now().UTC(),
				}

				/* act */
				providedJson, err := json.From(expectedEvent)
				if nil != err {
					panic(err)
				}

				actualEvent := Event{}
				json.To(providedJson, &actualEvent)

				/* assert */
				Expect(actualEvent).To(Equal(expectedEvent))

			})

		})
	})
})
