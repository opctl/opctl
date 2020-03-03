package clioutput

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	clicolorerFakes "github.com/opctl/opctl/cli/internal/clicolorer/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	"time"
)

var _ = Context("output", func() {
	Context("newOutput", func() {
		It("should return output", func() {
			/* arrange/act/assert */
			Expect(New(
				new(clicolorerFakes.FakeCliColorer),
				new(fakeWriter),
				new(fakeWriter),
			)).To(Not(BeNil()))
		})
	})
	_cliColorer := clicolorer.New()
	Context("Attention", func() {
		providedFormat := "dummyFormat %v %v"
		providedValues := []interface{}{"v1", "v2"}
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(fmt.Sprintln(_cliColorer.Attention(providedFormat, providedValues...)))

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Attention(providedFormat, providedValues...)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				To(Equal(expectedWriteArg))
		})
	})
	Context("Error", func() {
		providedFormat := "dummyFormat %v %v"
		providedValues := []interface{}{"v1", "v2"}
		It("should call errWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(fmt.Sprintln(_cliColorer.Error(providedFormat, providedValues...)))

			fakeErrWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				fakeErrWriter,
				new(fakeWriter),
			)

			/* act */
			objectUnderTest.Error(providedFormat, providedValues...)

			/* assert */
			Expect(fakeErrWriter.WriteArgsForCall(0)).
				To(Equal(expectedWriteArg))
		})
	})
	Context("Event", func() {
		Context("ContainerExited", func() {
			It("should call stdWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					ContainerExited: &model.ContainerExitedEvent{
						ContainerID: "dummyContainerID",
						OpRef:       "dummyOpRef",
						ExitCode:    1,
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_cliColorer.Info(
						"ContainerExited Id='%v' OpRef='%v' ExitCode='%v' Timestamp='%v'\n",
						providedEvent.ContainerExited.ContainerID,
						providedEvent.ContainerExited.OpRef,
						providedEvent.ContainerExited.ExitCode,
						providedEvent.Timestamp.Format(time.RFC3339),
					),
				))

				fakeStdWriter := new(fakeWriter)
				objectUnderTest := New(
					_cliColorer,
					new(fakeWriter),
					fakeStdWriter,
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeStdWriter.WriteArgsForCall(0)).
					To(Equal(expectedWriteArg))
			})
		})
		Context("ContainerStarted", func() {
			It("should call stdWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					ContainerStarted: &model.ContainerStartedEvent{
						ContainerID: "dummyContainerID",
						OpRef:       "dummyOpRef",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_cliColorer.Info(
						"ContainerStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
						providedEvent.ContainerStarted.ContainerID,
						providedEvent.ContainerStarted.OpRef,
						providedEvent.Timestamp.Format(time.RFC3339),
					),
				))

				fakeStdWriter := new(fakeWriter)
				objectUnderTest := New(
					_cliColorer,
					new(fakeWriter),
					fakeStdWriter,
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeStdWriter.WriteArgsForCall(0)).
					To(Equal(expectedWriteArg))
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
				expectedWriteArg := []byte(fmt.Sprint(string(providedEvent.ContainerStdErrWrittenTo.Data)))

				fakeErrWriter := new(fakeWriter)
				objectUnderTest := New(
					_cliColorer,
					fakeErrWriter,
					new(fakeWriter),
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeErrWriter.WriteArgsForCall(0)).
					To(Equal(expectedWriteArg))
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
				expectedWriteArg := []byte(fmt.Sprint(string(providedEvent.ContainerStdOutWrittenTo.Data)))

				fakeStdWriter := new(fakeWriter)
				objectUnderTest := New(
					_cliColorer,
					new(fakeWriter),
					fakeStdWriter,
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeStdWriter.WriteArgsForCall(0)).
					To(Equal(expectedWriteArg))
			})
		})
		Context("OpErred", func() {
			It("should call errWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					OpErred: &model.OpErredEvent{
						OpID:  "dummyOpID",
						OpRef: "dummyOpRef",
						Msg:   "dummyMsg",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_cliColorer.Error(
						"OpErred Id='%v' OpRef='%v' Timestamp='%v' Msg='%v'\n",
						providedEvent.OpErred.OpID,
						providedEvent.OpErred.OpRef,
						providedEvent.Timestamp.Format(time.RFC3339),
						providedEvent.OpErred.Msg,
					),
				))

				fakeErrWriter := new(fakeWriter)
				objectUnderTest := New(
					_cliColorer,
					fakeErrWriter,
					new(fakeWriter),
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeErrWriter.WriteArgsForCall(0)).
					To(Equal(expectedWriteArg))
			})
		})
		Context("OpEnded", func() {
			Context("Outcome==SUCCEEDED", func() {
				It("should call stdWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						OpEnded: &model.OpEndedEvent{
							OpID:    "dummyOpID",
							OpRef:   "dummyOpRef",
							Outcome: "SUCCEEDED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_cliColorer.Success(
							"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpID,
							providedEvent.OpEnded.OpRef,
							providedEvent.OpEnded.Outcome,
							providedEvent.Timestamp.Format(time.RFC3339),
						),
					))

					fakeStdWriter := new(fakeWriter)
					objectUnderTest := New(
						_cliColorer,
						new(fakeWriter),
						fakeStdWriter,
					)

					/* act */
					objectUnderTest.Event(providedEvent)

					/* assert */
					Expect(fakeStdWriter.WriteArgsForCall(0)).
						To(Equal(expectedWriteArg))
				})
			})
			Context("Outcome==KILLED", func() {
				It("should call stdWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						OpEnded: &model.OpEndedEvent{
							OpID:    "dummyOpID",
							OpRef:   "dummyOpRef",
							Outcome: "KILLED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_cliColorer.Info(
							"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpID,
							providedEvent.OpEnded.OpRef,
							providedEvent.OpEnded.Outcome,
							providedEvent.Timestamp.Format(time.RFC3339),
						),
					))

					fakeStdWriter := new(fakeWriter)
					objectUnderTest := New(
						_cliColorer,
						new(fakeWriter),
						fakeStdWriter,
					)

					/* act */
					objectUnderTest.Event(providedEvent)

					/* assert */
					Expect(fakeStdWriter.WriteArgsForCall(0)).
						To(Equal(expectedWriteArg))
				})
			})
			Context("Outcome==FAILED", func() {
				It("should call errWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						OpEnded: &model.OpEndedEvent{
							OpID:    "dummyOpID",
							OpRef:   "dummyOpRef",
							Outcome: "FAILED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_cliColorer.Error(
							"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpID,
							providedEvent.OpEnded.OpRef,
							providedEvent.OpEnded.Outcome,
							providedEvent.Timestamp.Format(time.RFC3339),
						),
					))

					fakeErrWriter := new(fakeWriter)
					objectUnderTest := New(
						_cliColorer,
						fakeErrWriter,
						new(fakeWriter),
					)

					/* act */
					objectUnderTest.Event(providedEvent)

					/* assert */
					Expect(fakeErrWriter.WriteArgsForCall(0)).
						To(Equal(expectedWriteArg))
				})
			})
		})
		Context("OpStarted", func() {
			It("should call stdWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					OpStarted: &model.OpStartedEvent{
						OpID:  "dummyOpID",
						OpRef: "dummyOpRef",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_cliColorer.Info(
						"OpStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
						providedEvent.OpStarted.OpID,
						providedEvent.OpStarted.OpRef,
						providedEvent.Timestamp.Format(time.RFC3339),
					),
				))

				fakeStdWriter := new(fakeWriter)
				objectUnderTest := New(
					_cliColorer,
					new(fakeWriter),
					fakeStdWriter,
				)

				/* act */
				objectUnderTest.Event(providedEvent)

				/* assert */
				Expect(fakeStdWriter.WriteArgsForCall(0)).
					To(Equal(expectedWriteArg))
			})
		})
	})
	Context("Info", func() {
		providedFormat := "dummyFormat %v %v"
		providedValues := []interface{}{"v1", "v2"}
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(fmt.Sprintln(
				_cliColorer.Info(providedFormat, providedValues...),
			))

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Info(providedFormat, providedValues...)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				To(Equal(expectedWriteArg))
		})
	})
	Context("Success", func() {
		providedFormat := "dummyFormat %v %v"
		providedValues := []interface{}{"v1", "v2"}
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(fmt.Sprintln(
				_cliColorer.Success(providedFormat, providedValues...),
			))

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Success(providedFormat, providedValues...)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				To(Equal(expectedWriteArg))
		})
	})
})
