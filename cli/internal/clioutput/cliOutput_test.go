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
	Context("Warning", func() {
		providedFormat := "dummyFormat %v %v"
		It("should call stdWriter w/ expected args", func() {
			/* arrange */
			expectedWriteArg := []byte(
				fmt.Sprintln(
					_cliColorer.Error(providedFormat),
				),
			)

			fakeErrWriter := new(fakeWriter)
			objectUnderTest := New(
				_cliColorer,
				new(fakeWriter),
				fakeErrWriter,
			)

			/* act */
			objectUnderTest.Warning(providedFormat)

			/* assert */
			Expect(fakeErrWriter.WriteArgsForCall(0)).
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
			Context("Call.Container truthy", func() {
				It("should call stdWriter w/ expected args", func() {
					/* arrange */
					imageRef := "imageRef"
					providedEvent := &model.Event{
						CallEnded: &model.CallEnded{
							Call: model.Call{
								Container: &model.ContainerCall{
									Image: &model.ContainerCallImage{
										Ref: &imageRef,
									},
								},
								ID: "ID",
							},
							Outcome: model.OpOutcomeSucceeded,
							Ref:     "ref",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(
						fmt.Sprintln(
							_cliColorer.Success(
								fmt.Sprintf(
									"ContainerExited Id='%v'%v Outcome='%v'%v Timestamp='%v'\n",
									providedEvent.CallEnded.Call.ID,
									fmt.Sprintf(" ImageRef='%v'", imageRef),
									providedEvent.CallEnded.Outcome,
									"",
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
			Context("Call.Op truthy", func() {

				Context("Outcome==FAILED", func() {
					It("should call errWriter w/ expected args", func() {
						/* arrange */
						providedEvent := &model.Event{
							CallEnded: &model.CallEnded{
								Call: model.Call{
									ID: "ID",
									Op: &model.OpCall{
										BaseCall: model.BaseCall{
											OpPath: "opPath",
										},
									},
								},
								Error: &model.CallEndedError{
									Message: "message",
								},
								Ref:     "ref",
								Outcome: "FAILED",
							},
							Timestamp: time.Now(),
						}
						expectedWriteArg := []byte(
							fmt.Sprintln(
								_cliColorer.Error(
									fmt.Sprintf(
										"OpEnded Id='%v' OpRef='%v' Outcome='%v'%v Timestamp='%v'\n",
										providedEvent.CallEnded.Call.ID,
										providedEvent.CallEnded.Call.Op.OpPath,
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
								Call: model.Call{
									ID: "ID",
									Op: &model.OpCall{
										BaseCall: model.BaseCall{
											OpPath: "opPath",
										},
									},
								},
								Ref:     "ref",
								Outcome: "SUCCEEDED",
							},
							Timestamp: time.Now(),
						}
						expectedWriteArg := []byte(
							fmt.Sprintln(
								_cliColorer.Success(
									fmt.Sprintf(
										"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
										providedEvent.CallEnded.Call.ID,
										providedEvent.CallEnded.Call.Op.OpPath,
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
								Call: model.Call{
									ID: "ID",
									Op: &model.OpCall{
										BaseCall: model.BaseCall{
											OpPath: "opPath",
										},
									},
								},
								Ref:     "ref",
								Outcome: "KILLED",
							},
							Timestamp: time.Now(),
						}
						expectedWriteArg := []byte(
							fmt.Sprintln(
								_cliColorer.Info(
									fmt.Sprintf(
										"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
										providedEvent.CallEnded.Call.ID,
										providedEvent.CallEnded.Call.Op.OpPath,
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
			})
			Context("Error truthy", func() {
				It("should call errWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						CallEnded: &model.CallEnded{
							Call: model.Call{
								ID: "ID",
							},
							Error: &model.CallEndedError{
								Message: "message",
							},
							Ref: "ref",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(
						fmt.Sprintln(
							_cliColorer.Error(
								fmt.Sprintf(
									"Error='%v' Id='%v' OpRef='%v' Timestamp='%v'\n",
									providedEvent.CallEnded.Error.Message,
									providedEvent.CallEnded.Call.ID,
									providedEvent.CallEnded.Ref,
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
		Context("CallStarted", func() {
			Context("Call.Container truthy", func() {
				It("should call stdWriter w/ expected args", func() {
					/* arrange */
					imageRef := "imageRef"
					providedEvent := &model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								Container: &model.ContainerCall{
									Image: &model.ContainerCallImage{
										Ref: &imageRef,
									},
								},
							},
							Ref: "ref",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(
						fmt.Sprintln(
							_cliColorer.Info(
								fmt.Sprintf(
									"ContainerStarted Id='%v' OpRef='%v' ImageRef='%v' Timestamp='%v'\n",
									providedEvent.CallStarted.Call.ID,
									providedEvent.CallStarted.Ref,
									imageRef,
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
			Describe("Call.Op truthy", func() {
				It("should call stdWriter w/ expected args", func() {
					/* arrange */
					providedEvent := &model.Event{
						CallStarted: &model.CallStarted{
							Call: model.Call{
								ID: "ID",
								Op: &model.OpCall{
									BaseCall: model.BaseCall{
										OpPath: "opPath",
									},
								},
							},
							Ref: "ref",
						},
						Timestamp: time.Now(),
					}
					expectedWriteArg := []byte(
						fmt.Sprintln(
							_cliColorer.Info(
								fmt.Sprintf(
									"OpStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
									providedEvent.CallStarted.Call.ID,
									providedEvent.CallStarted.Call.Op.OpPath,
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
	Context("DisableColor", func() {
		It("should call colorer to disable color", func() {
			/* arrange */
			fakeCliColorer := new(clicolorerFakes.FakeCliColorer)

			objectUnderTest := New(
				fakeCliColorer,
				new(fakeWriter),
				new(fakeWriter),
			)

			/* act */
			objectUnderTest.DisableColor()

			/* assert */
			Expect(fakeCliColorer.DisableColorCallCount()).To(Equal(1))
		})
	})
})
