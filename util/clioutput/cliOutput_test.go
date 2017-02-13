package clioutput

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"time"
)

var _ = Context("output", func() {
	Context("newOutput", func() {
		It("should return output", func() {
			/* arrange/act/assert */
			Expect(New(
				colorer.New(),
				new(fakeWriter),
				new(fakeWriter),
			)).Should(Not(BeNil()))
		})
	})
	_colorer := colorer.New()
	Context("Attention", func() {
		providedFormat := "dummyFormat %v %v"
		providedValues := []interface{}{"v1", "v2"}
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(fmt.Sprintln(_colorer.Attention(providedFormat, providedValues...)))

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_colorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Attention(providedFormat, providedValues...)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				Should(Equal(expectedWriteArg))
		})
	})
	Context("Error", func() {
		providedFormat := "dummyFormat %v %v"
		providedValues := []interface{}{"v1", "v2"}
		It("should call errWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(fmt.Sprintln(_colorer.Error(providedFormat, providedValues...)))

			fakeErrWriter := new(fakeWriter)
			objectUnderTest := New(
				_colorer,
				fakeErrWriter,
				new(fakeWriter),
			)

			/* act */
			objectUnderTest.Error(providedFormat, providedValues...)

			/* assert */
			Expect(fakeErrWriter.WriteArgsForCall(0)).
				Should(Equal(expectedWriteArg))
		})
	})
	Context("Event", func() {
		Context("ContainerExited", func() {
			It("should call stdWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					ContainerExited: &model.ContainerExitedEvent{
						ContainerId: "dummyContainerId",
						OpRef:       "dummyOpRef",
						ExitCode:    1,
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_colorer.Info(
						"ContainerExited Id='%v' OpRef='%v' ExitCode='%v' Timestamp='%v'\n",
						providedEvent.ContainerExited.ContainerId,
						providedEvent.ContainerExited.OpRef,
						providedEvent.ContainerExited.ExitCode,
						providedEvent.Timestamp.Format(time.RFC3339),
					),
				))

				fakeStdWriter := new(fakeWriter)
				objectUnderTest := New(
					_colorer,
					new(fakeWriter),
					fakeStdWriter,
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeStdWriter.WriteArgsForCall(0)).
					Should(Equal(expectedWriteArg))
			})
		})
		Context("ContainerStarted", func() {
			It("should call stdWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					ContainerStarted: &model.ContainerStartedEvent{
						ContainerId: "dummyContainerId",
						OpRef:       "dummyOpRef",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_colorer.Info(
						"ContainerStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
						providedEvent.ContainerStarted.ContainerId,
						providedEvent.ContainerStarted.OpRef,
						providedEvent.Timestamp.Format(time.RFC3339),
					),
				))

				fakeStdWriter := new(fakeWriter)
				objectUnderTest := New(
					_colorer,
					new(fakeWriter),
					fakeStdWriter,
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeStdWriter.WriteArgsForCall(0)).
					Should(Equal(expectedWriteArg))
			})
		})
		Context("ContainerStdErrWrittenTo", func() {
			It("should call stdWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenToEvent{
						Data: []byte("dummyData"),
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(string(providedEvent.ContainerStdErrWrittenTo.Data)))

				fakeErrWriter := new(fakeWriter)
				objectUnderTest := New(
					_colorer,
					fakeErrWriter,
					new(fakeWriter),
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeErrWriter.WriteArgsForCall(0)).
					Should(Equal(expectedWriteArg))
			})
		})
		Context("ContainerOutWrittenTo", func() {
			It("should call stdWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
						Data: []byte("dummyData"),
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(string(providedEvent.ContainerStdOutWrittenTo.Data)))

				fakeStdWriter := new(fakeWriter)
				objectUnderTest := New(
					_colorer,
					new(fakeWriter),
					fakeStdWriter,
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeStdWriter.WriteArgsForCall(0)).
					Should(Equal(expectedWriteArg))
			})
		})
		Context("OpEncounteredError", func() {
			It("should call errWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						OpId:  "dummyOpId",
						OpRef: "dummyOpRef",
						Msg:   "dummyMsg",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_colorer.Error(
						"OpEncounteredError Id='%v' OpRef='%v' Timestamp='%v' Msg='%v'\n",
						providedEvent.OpEncounteredError.OpId,
						providedEvent.OpEncounteredError.OpRef,
						providedEvent.Timestamp.Format(time.RFC3339),
						providedEvent.OpEncounteredError.Msg,
					),
				))

				fakeErrWriter := new(fakeWriter)
				objectUnderTest := New(
					_colorer,
					fakeErrWriter,
					new(fakeWriter),
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeErrWriter.WriteArgsForCall(0)).
					Should(Equal(expectedWriteArg))
			})
		})
		Context("OpEnded", func() {
			Context("Outcome==SUCCEEDED", func() {
				It("should call stdWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						OpEnded: &model.OpEndedEvent{
							OpId:    "dummyOpId",
							OpRef:   "dummyOpRef",
							Outcome: "SUCCEEDED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_colorer.Success(
							"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpId,
							providedEvent.OpEnded.OpRef,
							providedEvent.OpEnded.Outcome,
							providedEvent.Timestamp.Format(time.RFC3339),
						),
					))

					fakeStdWriter := new(fakeWriter)
					objectUnderTest := New(
						_colorer,
						new(fakeWriter),
						fakeStdWriter,
					)

					/* act */
					objectUnderTest.Event(providedEvent)

					/* assert */
					Expect(fakeStdWriter.WriteArgsForCall(0)).
						Should(Equal(expectedWriteArg))
				})
			})
			Context("Outcome==KILLED", func() {
				It("should call stdWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						OpEnded: &model.OpEndedEvent{
							OpId:    "dummyOpId",
							OpRef:   "dummyOpRef",
							Outcome: "KILLED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_colorer.Info(
							"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpId,
							providedEvent.OpEnded.OpRef,
							providedEvent.OpEnded.Outcome,
							providedEvent.Timestamp.Format(time.RFC3339),
						),
					))

					fakeStdWriter := new(fakeWriter)
					objectUnderTest := New(
						_colorer,
						new(fakeWriter),
						fakeStdWriter,
					)

					/* act */
					objectUnderTest.Event(providedEvent)

					/* assert */
					Expect(fakeStdWriter.WriteArgsForCall(0)).
						Should(Equal(expectedWriteArg))
				})
			})
			Context("Outcome==FAILED", func() {
				It("should call errWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						OpEnded: &model.OpEndedEvent{
							OpId:    "dummyOpId",
							OpRef:   "dummyOpRef",
							Outcome: "FAILED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_colorer.Error(
							"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpId,
							providedEvent.OpEnded.OpRef,
							providedEvent.OpEnded.Outcome,
							providedEvent.Timestamp.Format(time.RFC3339),
						),
					))

					fakeErrWriter := new(fakeWriter)
					objectUnderTest := New(
						_colorer,
						fakeErrWriter,
						new(fakeWriter),
					)

					/* act */
					objectUnderTest.Event(providedEvent)

					/* assert */
					Expect(fakeErrWriter.WriteArgsForCall(0)).
						Should(Equal(expectedWriteArg))
				})
			})
		})
		Context("OpStarted", func() {
			It("should call stdWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					OpStarted: &model.OpStartedEvent{
						OpId:  "dummyOpId",
						OpRef: "dummyOpRef",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_colorer.Info(
						"OpStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
						providedEvent.OpStarted.OpId,
						providedEvent.OpStarted.OpRef,
						providedEvent.Timestamp.Format(time.RFC3339),
					),
				))

				fakeStdWriter := new(fakeWriter)
				objectUnderTest := New(
					_colorer,
					new(fakeWriter),
					fakeStdWriter,
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeStdWriter.WriteArgsForCall(0)).
					Should(Equal(expectedWriteArg))
			})
		})
	})
	Context("Info", func() {
		providedFormat := "dummyFormat %v %v"
		providedValues := []interface{}{"v1", "v2"}
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(fmt.Sprintln(
				_colorer.Info(providedFormat, providedValues...),
			))

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_colorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Info(providedFormat, providedValues...)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				Should(Equal(expectedWriteArg))
		})
	})
	Context("Success", func() {
		providedFormat := "dummyFormat %v %v"
		providedValues := []interface{}{"v1", "v2"}
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(fmt.Sprintln(
				_colorer.Success(providedFormat, providedValues...),
			))

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_colorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Success(providedFormat, providedValues...)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				Should(Equal(expectedWriteArg))
		})
	})
})
