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

		Context("ls", func() {
			Context("w/ path", func() {

				It("should call core.ListPackages w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedPath := "dummyPath"
					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls", "-c", expectedPath})

					/* assert */
					Expect(fakeCore.ListPackagesCallCount()).Should(Equal(1))
					actualPath := fakeCore.ListPackagesArgsForCall(0)
					Expect(actualPath).Should(Equal(expectedPath))
				})
			})
			Context("w/out path", func() {

				It("should call core.ListPackages w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedPath := ".opspec"
					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls"})

					/* assert */
					Expect(fakeCore.ListPackagesCallCount()).Should(Equal(1))
					actualPath := fakeCore.ListPackagesArgsForCall(0)
					Expect(actualPath).Should(Equal(expectedPath))
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

			Context("kill", func() {
				It("should call core.OpKill w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedOpId := "dummyOpId"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "op", "kill", expectedOpId})

					/* assert */
					Expect(fakeCore.OpKillCallCount()).Should(Equal(1))
					Expect(fakeCore.OpKillArgsForCall(0)).Should(Equal(expectedOpId))
				})
			})

		})

		Context("pkg", func() {

			Context("create", func() {
				Context("w/ path", func() {
					It("should call core.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedPkgName := "dummyPkgName"
						expectedPath := "dummyPath"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "pkg", "create", "-c", expectedPath, expectedPkgName})

						/* assert */
						Expect(fakeCore.CreateCallCount()).Should(Equal(1))
						actualPath, actualPkgDescription, actualPkgName := fakeCore.CreateArgsForCall(0)
						Expect(actualPkgName).Should(Equal(expectedPkgName))
						Expect(actualPkgDescription).Should(BeEmpty())
						Expect(actualPath).Should(Equal(expectedPath))
					})
				})

				Context("w/out path", func() {
					It("should call core.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedPkgName := "dummyPkgName"
						expectedPath := ".opspec"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "pkg", "create", expectedPkgName})

						/* assert */
						Expect(fakeCore.CreateCallCount()).Should(Equal(1))
						actualPath, actualPkgDescription, actualPkgName := fakeCore.CreateArgsForCall(0)
						Expect(actualPkgName).Should(Equal(expectedPkgName))
						Expect(actualPkgDescription).Should(BeEmpty())
						Expect(actualPath).Should(Equal(expectedPath))
					})
				})
				Context("w/ description", func() {
					It("should call core.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedPkgName := "dummyPkgName"
						expectedPkgDescription := "dummyPkgDescription"
						expectedPath := ".opspec"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "pkg", "create", "-d", expectedPkgDescription, expectedPkgName})

						/* assert */
						Expect(fakeCore.CreateCallCount()).Should(Equal(1))
						actualPath, actualPkgDescription, actualPkgName := fakeCore.CreateArgsForCall(0)
						Expect(actualPkgName).Should(Equal(expectedPkgName))
						Expect(actualPkgDescription).Should(Equal(expectedPkgDescription))
						Expect(actualPath).Should(Equal(expectedPath))
					})
				})

				Context("w/out description", func() {
					It("should call core.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedName := "dummyPkgName"
						expectedPath := ".opspec"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "pkg", "create", expectedName})

						/* assert */
						Expect(fakeCore.CreateCallCount()).Should(Equal(1))
						actualPath, actualPkgDescription, actualPkgName := fakeCore.CreateArgsForCall(0)
						Expect(actualPkgName).Should(Equal(expectedName))
						Expect(actualPkgDescription).Should(BeEmpty())
						Expect(actualPath).Should(Equal(expectedPath))
					})
				})
			})

			Context("set", func() {

				Context("description", func() {

					It("should call core.PkgSetDescription w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						pkgRef := ".opspec/dummyPkgName"
						expectedDescription := "dummyPkgDescription"

						objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

						/* act */
						objectUnderTest.Run([]string{"opctl", "pkg", "set", "description", expectedDescription, pkgRef})

						/* assert */
						Expect(fakeCore.PkgSetDescriptionCallCount()).Should(Equal(1))
						actualPkgDescription, actualPkgRef := fakeCore.PkgSetDescriptionArgsForCall(0)
						Expect(actualPkgRef).Should(Equal(pkgRef))
						Expect(actualPkgDescription).Should(Equal(expectedDescription))
					})

				})

			})

			Context("validate", func() {

				It("should call core.PkgValidate w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					pkgRef := ".opspec/dummyPkgName"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "pkg", "validate", pkgRef})

					/* assert */
					Expect(fakeCore.PkgSetDescriptionCallCount()).Should(Equal(1))
					actualPkgRef := fakeCore.PkgValidateArgsForCall(0)
					Expect(actualPkgRef).Should(Equal(pkgRef))
				})

			})

		})

		Context("run", func() {
			Context("with two op run args", func() {
				It("should call core.RunOp w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedOpArgs := []string{"arg1Name=arg1Value", "arg2Name=arg2Value"}
					expectedPkgRef := ".opspec/dummyPkgName"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "run", "-a", expectedOpArgs[0], "-a", expectedOpArgs[1], expectedPkgRef})

					/* assert */
					Expect(fakeCore.RunOpCallCount()).Should(Equal(1))
					actualOpArgs, actualOpUrl := fakeCore.RunOpArgsForCall(0)
					Expect(actualOpUrl).Should(Equal(expectedPkgRef))
					Expect(actualOpArgs).Should(Equal(expectedOpArgs))
				})
			})

			Context("with zero op run args", func() {
				It("should call core.RunOp w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedPkgRef := ".opspec/dummyPkgName"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "run", expectedPkgRef})

					/* assert */
					Expect(fakeCore.RunOpCallCount()).Should(Equal(1))

					actualOpArgs, actualPkgRef := fakeCore.RunOpArgsForCall(0)
					Expect(actualPkgRef).Should(Equal(expectedPkgRef))
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
