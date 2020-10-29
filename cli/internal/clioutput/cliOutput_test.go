package clioutput

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	clicolorerFakes "github.com/opctl/opctl/cli/internal/clicolorer/fakes"
	"github.com/opctl/opctl/sdks/go/model"
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
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(
				fmt.Sprintln(
					_cliColorer.Attention(providedFormat),
				),
			)

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Attention(providedFormat)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				To(Equal(expectedWriteArg))
		})
	})
	Context("Error", func() {
		providedFormat := "dummyFormat %v %v"
		It("should call errWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(
				fmt.Sprintln(
					_cliColorer.Error(providedFormat),
				),
			)

			fakeErrWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				fakeErrWriter,
				new(fakeWriter),
			)

			/* act */
			objectUnderTest.Error(providedFormat)

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
					ContainerExited: &model.ContainerExited{
						ContainerID: "dummyContainerID",
						OpRef:       "dummyOpRef",
						ExitCode:    1,
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(
					fmt.Sprintln(
						_cliColorer.Info(
							fmt.Sprintf(
								"ContainerExited Id='%v' OpRef='%v' ExitCode='%v' Timestamp='%v'\n",
								providedEvent.ContainerExited.ContainerID,
								providedEvent.ContainerExited.OpRef,
								providedEvent.ContainerExited.ExitCode,
								providedEvent.Timestamp.Format(time.RFC3339),
							),
						),
					),
				)

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
					ContainerStarted: &model.ContainerStarted{
						ContainerID: "dummyContainerID",
						OpRef:       "dummyOpRef",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(
					fmt.Sprintln(
						_cliColorer.Info(
							fmt.Sprintf(
								"ContainerStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
								providedEvent.ContainerStarted.ContainerID,
								providedEvent.ContainerStarted.OpRef,
								providedEvent.Timestamp.Format(time.RFC3339),
							),
						),
					),
				)

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
					ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenTo{
						Data: []byte("dummyData"),
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(string(providedEvent.ContainerStdErrWrittenTo.Data))

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
					ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
						Data: []byte("dummyData"),
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(string(providedEvent.ContainerStdOutWrittenTo.Data))

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
		Context("CallEnded", func() {
			Context("Outcome==FAILED", func() {
				It("should call errWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						CallEnded: &model.CallEnded{
							Error: &model.CallEndedError{
								Message: "message",
							},
							CallID:   "dummyOpID",
							CallType: model.CallTypeOp,
							Ref:      "dummyOpRef",
							Outcome:  "FAILED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(
						fmt.Sprintln(
							_cliColorer.Error(
								fmt.Sprintf(
									"OpEnded Id='%v' OpRef='%v' Outcome='%v'%v Timestamp='%v'\n",
									providedEvent.CallEnded.CallID,
									providedEvent.CallEnded.Ref,
									providedEvent.CallEnded.Outcome,
									fmt.Sprintf(" Error='%v'", providedEvent.CallEnded.Error.Message),
									providedEvent.Timestamp.Format(time.RFC3339),
								),
							),
						),
					)

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
			Context("Outcome==SUCCEEDED", func() {
				It("should call stdWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						CallEnded: &model.CallEnded{
							CallID:   "dummyOpID",
							CallType: model.CallTypeOp,
							Ref:      "dummyOpRef",
							Outcome:  "SUCCEEDED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(
						fmt.Sprintln(
							_cliColorer.Success(
								fmt.Sprintf(
									"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
									providedEvent.CallEnded.CallID,
									providedEvent.CallEnded.Ref,
									providedEvent.CallEnded.Outcome,
									providedEvent.Timestamp.Format(time.RFC3339),
								),
							),
						),
					)

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
						CallEnded: &model.CallEnded{
							CallID:   "dummyOpID",
							CallType: model.CallTypeOp,
							Ref:      "dummyOpRef",
							Outcome:  "KILLED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(
						fmt.Sprintln(
							_cliColorer.Info(
								fmt.Sprintf(
									"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
									providedEvent.CallEnded.CallID,
									providedEvent.CallEnded.Ref,
									providedEvent.CallEnded.Outcome,
									providedEvent.Timestamp.Format(time.RFC3339),
								),
							),
						),
					)

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
						CallEnded: &model.CallEnded{
							CallID:   "dummyOpID",
							CallType: model.CallTypeOp,
							Ref:      "dummyOpRef",
							Outcome:  "FAILED",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(
						fmt.Sprintln(
							_cliColorer.Error(
								fmt.Sprintf(
									"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
									providedEvent.CallEnded.CallID,
									providedEvent.CallEnded.Ref,
									providedEvent.CallEnded.Outcome,
									providedEvent.Timestamp.Format(time.RFC3339),
								),
							),
						),
					)

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
					CallStarted: &model.CallStarted{
						CallID:   "dummyOpID",
						CallType: model.CallTypeOp,
						OpRef:    "dummyOpRef",
					},
					Timestamp: time.Now(),
				}
				expectedWriteArg := []byte(
					fmt.Sprintln(
						_cliColorer.Info(
							fmt.Sprintf(
								"OpStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
								providedEvent.CallStarted.CallID,
								providedEvent.CallStarted.OpRef,
								providedEvent.Timestamp.Format(time.RFC3339),
							),
						),
					),
				)

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
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(
				fmt.Sprintln(
					_cliColorer.Info(providedFormat),
				),
			)

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Info(providedFormat)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				To(Equal(expectedWriteArg))
		})
	})
	Context("Success", func() {
		providedFormat := "dummyFormat %v %v"
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(
				fmt.Sprintln(
					_cliColorer.Success(providedFormat),
				),
			)

			fakeStdWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				new(fakeWriter),
				fakeStdWriter,
			)

			/* act */
			objectUnderTest.Success(providedFormat)

			/* assert */
			Expect(fakeStdWriter.WriteArgsForCall(0)).
				To(Equal(expectedWriteArg))
		})
	})
})
