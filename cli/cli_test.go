package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/cli/core"
	"github.com/opspec-io/opctl/util/clicolorer"
)

var _ = Context("cli", func() {
	Context("Run", func() {

		Context("--no-color", func() {
			It("should set color.NoColor", func() {
				/* arrange */
				fakeCliColorer := new(clicolorer.Fake)

				objectUnderTest := newCli(new(core.Fake), fakeCliColorer)

				/* act */
				objectUnderTest.Run([]string{"opctl", "--no-color", "-v"})

				/* assert */
				Expect(fakeCliColorer.DisableCallCount()).Should(Equal(1))
			})
		})

		Context("collection", func() {

			Context("create", func() {

				Context("with description", func() {
					It("should call core.CreateCollection w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedCollectionName := "dummyCollectionName"
						expectedCollectionDescription := "dummyCollectionDescription"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "collection", "create", "-d", expectedCollectionDescription, expectedCollectionName})

						/* assert */
						Expect(fakeCore.CreateCollectionCallCount()).Should(Equal(1))
						actualCollectionDescription, actualCollectionName := fakeCore.CreateCollectionArgsForCall(0)
						Expect(actualCollectionName).Should(Equal(expectedCollectionName))
						Expect(actualCollectionDescription).Should(Equal(expectedCollectionDescription))
					})
				})

				Context("with no description", func() {
					It("should call core.CreateCollectionUseCase w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedCollectionName := "dummyCollectionName"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "collection", "create", expectedCollectionName})

						/* assert */
						Expect(fakeCore.CreateCollectionCallCount()).Should(Equal(1))
						actualCollectionDescription, actualCollectionName := fakeCore.CreateCollectionArgsForCall(0)
						Expect(actualCollectionName).Should(Equal(expectedCollectionName))
						Expect(actualCollectionDescription).Should(BeEmpty())
					})
				})
			})

			Context("set", func() {

				Context("description", func() {
					It("should call core.SetCollectionDescription w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedCollectionDescription := "dummyCollectionDescription"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "collection", "set", "description", expectedCollectionDescription})

						/* assert */
						Expect(fakeCore.SetCollectionDescriptionCallCount()).Should(Equal(1))
						Expect(fakeCore.SetCollectionDescriptionArgsForCall(0)).Should(Equal(expectedCollectionDescription))
					})
				})

			})

		})

		Context("events", func() {
			It("should call core.StreamEvents w/ expected args", func() {
				/* arrange */
				fakeCore := new(core.Fake)

				objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

				/* act */
				objectUnderTest.Run([]string{"opctl", "events"})

				/* assert */
				Expect(fakeCore.StreamEventsCallCount()).Should(Equal(1))
			})
		})

		Context("kill", func() {
			It("should call core.KillOp w/ expected args", func() {
				/* arrange */
				fakeCore := new(core.Fake)

				expectedOpId := "dummyOpId"

				objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

				/* act */
				objectUnderTest.Run([]string{"opctl", "kill", expectedOpId})

				/* assert */
				Expect(fakeCore.KillOpCallCount()).Should(Equal(1))
				Expect(fakeCore.KillOpArgsForCall(0)).Should(Equal(expectedOpId))
			})
		})

		Context("ls", func() {
			Context("with collection", func() {

				It("should call core.ListOpsInCollection w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedCollection := "dummyCollection"
					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls", "-c", expectedCollection})

					/* assert */
					Expect(fakeCore.ListOpsInCollectionCallCount()).Should(Equal(1))
					actualCollection := fakeCore.ListOpsInCollectionArgsForCall(0)
					Expect(actualCollection).Should(Equal(expectedCollection))
				})
			})
			Context("without collection", func() {

				It("should call core.ListOpsInCollection w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedCollection := ".opspec"
					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls"})

					/* assert */
					Expect(fakeCore.ListOpsInCollectionCallCount()).Should(Equal(1))
					actualCollection := fakeCore.ListOpsInCollectionArgsForCall(0)
					Expect(actualCollection).Should(Equal(expectedCollection))
				})
			})
		})

		Context("node", func() {

			Context("create", func() {

				It("should call core.NodeCreate w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "node", "create"})

					/* assert */
					Expect(fakeCore.NodeCreateCallCount()).Should(Equal(1))
				})

			})

			Context("kill", func() {

				It("should call core.NodeKill w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "node", "kill"})

					/* assert */
					Expect(fakeCore.NodeKillCallCount()).Should(Equal(1))
				})

			})
		})

		Context("op", func() {

			Context("create", func() {
				Context("with collection", func() {
					It("should call core.CreateOp w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedOpName := "dummyOpName"
						expectedCollection := "dummyCollection"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", "-c", expectedCollection, expectedOpName})

						/* assert */
						Expect(fakeCore.CreateOpCallCount()).Should(Equal(1))
						actualCollection, actualOpDescription, actualOpName := fakeCore.CreateOpArgsForCall(0)
						Expect(actualOpName).Should(Equal(expectedOpName))
						Expect(actualOpDescription).Should(BeEmpty())
						Expect(actualCollection).Should(Equal(expectedCollection))
					})
				})

				Context("with no collection", func() {
					It("should call core.CreateOp w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedOpName := "dummyOpName"
						expectedCollection := ".opspec"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", expectedOpName})

						/* assert */
						Expect(fakeCore.CreateOpCallCount()).Should(Equal(1))
						actualCollection, actualOpDescription, actualOpName := fakeCore.CreateOpArgsForCall(0)
						Expect(actualOpName).Should(Equal(expectedOpName))
						Expect(actualOpDescription).Should(BeEmpty())
						Expect(actualCollection).Should(Equal(expectedCollection))
					})
				})
				Context("with description", func() {
					It("should call core.CreateOp w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedOpName := "dummyOpName"
						expectedOpDescription := "dummyOpDescription"
						expectedCollection := ".opspec"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", "-d", expectedOpDescription, expectedOpName})

						/* assert */
						Expect(fakeCore.CreateOpCallCount()).Should(Equal(1))
						actualCollection, actualOpDescription, actualOpName := fakeCore.CreateOpArgsForCall(0)
						Expect(actualOpName).Should(Equal(expectedOpName))
						Expect(actualOpDescription).Should(Equal(expectedOpDescription))
						Expect(actualCollection).Should(Equal(expectedCollection))
					})
				})

				Context("with no description", func() {
					It("should call core.CreateOp w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedName := "dummyOpName"
						expectedCollection := ".opspec"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", expectedName})

						/* assert */
						Expect(fakeCore.CreateOpCallCount()).Should(Equal(1))
						actualCollection, actualOpDescription, actualOpName := fakeCore.CreateOpArgsForCall(0)
						Expect(actualOpName).Should(Equal(expectedName))
						Expect(actualOpDescription).Should(BeEmpty())
						Expect(actualCollection).Should(Equal(expectedCollection))
					})
				})
			})

			Context("set", func() {

				Context("description", func() {
					Context("with collection", func() {
						It("should call core.SetOpDescription w/ expected args", func() {
							/* arrange */
							fakeCore := new(core.Fake)

							expectedName := "dummyOpName"
							expectedDescription := "dummyOpDescription"
							expectedCollection := "dummyCollection"

							objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

							/* act */
							objectUnderTest.Run([]string{"opctl", "op", "set", "description", "-c", expectedCollection, expectedDescription, expectedName})

							/* assert */
							Expect(fakeCore.SetOpDescriptionCallCount()).Should(Equal(1))
							actualCollection, actualOpDescription, actualOpName := fakeCore.SetOpDescriptionArgsForCall(0)
							Expect(actualOpName).Should(Equal(expectedName))
							Expect(actualOpDescription).Should(Equal(expectedDescription))
							Expect(actualCollection).Should(Equal(expectedCollection))
						})
					})

					Context("with no collection", func() {
						It("should call core.SetOpDescription w/ expected args", func() {
							/* arrange */
							fakeCore := new(core.Fake)

							expectedName := "dummyOpName"
							expectedDescription := "dummyOpDescription"
							expectedCollection := ".opspec"

							objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

							/* act */
							objectUnderTest.Run([]string{"opctl", "op", "set", "description", expectedDescription, expectedName})

							/* assert */
							Expect(fakeCore.SetOpDescriptionCallCount()).Should(Equal(1))
							actualCollection, actualOpDescription, actualOpName := fakeCore.SetOpDescriptionArgsForCall(0)
							Expect(actualOpName).Should(Equal(expectedName))
							Expect(actualOpDescription).Should(Equal(expectedDescription))
							Expect(actualCollection).Should(Equal(expectedCollection))
						})
					})
				})

			})

		})

		Context("run", func() {
			Context("with collection", func() {
				It("should call core.RunOp w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedCollection := "dummyCollection"
					expectedOpUrl := "dummyOpUrl"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "run", "-c", expectedCollection, expectedOpUrl})

					/* assert */
					Expect(fakeCore.RunOpCallCount()).Should(Equal(1))
					actualOpArgs, actualCollection, actualOpUrl := fakeCore.RunOpArgsForCall(0)
					Expect(actualOpUrl).Should(Equal(expectedOpUrl))
					Expect(actualCollection).Should(Equal(expectedCollection))
					Expect(actualOpArgs).Should(BeEmpty())
				})
			})
			Context("with two op run args", func() {
				It("should call core.RunOp w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedOpArgs := []string{"arg1Name=arg1Value", "arg2Name=arg2Value"}
					expectedCollection := ".opspec"
					expectedOpUrl := "dummyOpUrl"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "run", "-a", expectedOpArgs[0], "-a", expectedOpArgs[1], expectedOpUrl})

					/* assert */
					Expect(fakeCore.RunOpCallCount()).Should(Equal(1))
					actualOpArgs, actualCollection, actualOpUrl := fakeCore.RunOpArgsForCall(0)
					Expect(actualOpUrl).Should(Equal(expectedOpUrl))
					Expect(actualCollection).Should(Equal(expectedCollection))
					Expect(actualOpArgs).Should(Equal(expectedOpArgs))
				})
			})

			Context("with zero op run args", func() {
				It("should call core.RunOp w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedOpUrl := "dummyOpUrl"
					expectedCollection := ".opspec"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "run", expectedOpUrl})

					/* assert */
					Expect(fakeCore.RunOpCallCount()).Should(Equal(1))

					actualOpArgs, actualCollection, actualOpUrl := fakeCore.RunOpArgsForCall(0)
					Expect(actualOpUrl).Should(Equal(expectedOpUrl))
					Expect(actualCollection).Should(Equal(expectedCollection))
					Expect(actualOpArgs).Should(BeEmpty())
				})
			})
		})
	})

	Context("self-update", func() {

		Context("with channel flag", func() {

			It("should call core.SelfUpdate with expected releaseChannel", func() {
				/* arrange */
				expectedChannel := "beta"

				fakeCore := new(core.Fake)

				objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

				/* act */
				objectUnderTest.Run([]string{"opctl", "self-update", "-c", expectedChannel})

				/* assert */
				Expect(fakeCore.SelfUpdateCallCount()).Should(Equal(1))

				actualChannel := fakeCore.SelfUpdateArgsForCall(0)
				Expect(actualChannel).Should(Equal(expectedChannel))
			})
		})

		Context("without channel flag", func() {

			It("should call core.SelfUpdate with expected releaseChannel", func() {
				/* arrange */
				expectedChannel := "stable"

				fakeCore := new(core.Fake)

				objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

				/* act */
				objectUnderTest.Run([]string{"opctl", "self-update"})

				/* assert */
				Expect(fakeCore.SelfUpdateCallCount()).Should(Equal(1))

				actualChannel := fakeCore.SelfUpdateArgsForCall(0)
				Expect(actualChannel).Should(Equal(expectedChannel))
			})
		})

	})

})
