package clioutput

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/opspec-io/sdk-golang/model"
	"time"
)

var _ = Context("output", func() {
	Context("newOutput", func() {
		It("should return output", func() {
			/* arrange/act/assert */
			Expect(New(
				new(clicolorer.Fake),
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
						ContainerId: "dummyContainerId",
						PkgRef:      "dummyPkgRef",
						ExitCode:    1,
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_cliColorer.Info(
						"ContainerExited Id='%v' PkgRef='%v' ExitCode='%v' Timestamp='%v'\n",
						providedEvent.ContainerExited.ContainerId,
						providedEvent.ContainerExited.PkgRef,
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
						ContainerId: "dummyContainerId",
						PkgRef:      "dummyPkgRef",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_cliColorer.Info(
						"ContainerStarted Id='%v' PkgRef='%v' Timestamp='%v'\n",
						providedEvent.ContainerStarted.ContainerId,
						providedEvent.ContainerStarted.PkgRef,
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
				expectedWriteArg := []byte(fmt.Sprintln(string(providedEvent.ContainerStdErrWrittenTo.Data)))

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
				expectedWriteArg := []byte(fmt.Sprintln(string(providedEvent.ContainerStdOutWrittenTo.Data)))

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
		Context("OpEncounteredError", func() {
			It("should call errWriter w/ expected args", func() {
				/* arrange */
				providedEvent := &model.Event{
					OpEncounteredError: &model.OpEncounteredErrorEvent{
						OpId:   "dummyOpId",
						PkgRef: "dummyPkgRef",
						Msg:    "dummyMsg",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_cliColorer.Error(
						"OpEncounteredError Id='%v' PkgRef='%v' Timestamp='%v' Msg='%v'\n",
						providedEvent.OpEncounteredError.OpId,
						providedEvent.OpEncounteredError.PkgRef,
						providedEvent.Timestamp.Format(time.RFC3339),
						providedEvent.OpEncounteredError.Msg,
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
							OpId:    "dummyOpId",
							PkgRef:  "dummyPkgRef",
							Outcome: "SUCCEEDED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_cliColorer.Success(
							"OpEnded Id='%v' PkgRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpId,
							providedEvent.OpEnded.PkgRef,
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
							OpId:    "dummyOpId",
							PkgRef:  "dummyPkgRef",
							Outcome: "KILLED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_cliColorer.Info(
							"OpEnded Id='%v' PkgRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpId,
							providedEvent.OpEnded.PkgRef,
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
							OpId:    "dummyOpId",
							PkgRef:  "dummyPkgRef",
							Outcome: "FAILED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(fmt.Sprintln(
						_cliColorer.Error(
							"OpEnded Id='%v' PkgRef='%v' Outcome='%v' Timestamp='%v'\n",
							providedEvent.OpEnded.OpId,
							providedEvent.OpEnded.PkgRef,
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
						OpId:   "dummyOpId",
						PkgRef: "dummyPkgRef",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(fmt.Sprintln(
					_cliColorer.Info(
						"OpStarted Id='%v' PkgRef='%v' Timestamp='%v'\n",
						providedEvent.OpStarted.OpId,
						providedEvent.OpStarted.PkgRef,
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
