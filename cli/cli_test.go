package main

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	nodeCreateCmd "github.com/opctl/opctl/cli/cmds/node/create"
	"github.com/opctl/opctl/cli/core"
	"github.com/opctl/opctl/cli/types"
	"github.com/opctl/opctl/cli/util/clicolorer"
)

var _ = Context("cli", func() {
	Context("Run", func() {

		Context("--no-color", func() {
			It("should set color.NoColor", func() {
				/* arrange */
				fakeCliColorer := new(clicolorer.Fake)

				objectUnderTest := newCli(
					new(core.Fake),
					fakeCliColorer,
					new(nodeCreateCmd.FakeInvoker),
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "--no-color", "ls"})

				/* assert */
				Expect(fakeCliColorer.DisableCallCount()).To(Equal(1))
			})
		})

		Context("events", func() {
			It("should call core.Events w/ expected args", func() {
				/* arrange */
				providedCtx := context.Background()
				fakeCore := new(core.Fake)

				objectUnderTest := newCli(
					fakeCore,
					new(clicolorer.Fake),
					new(nodeCreateCmd.FakeInvoker),
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "events"})

				/* assert */
				actualCtx := fakeCore.EventsArgsForCall(0)
				Expect(actualCtx).To(Equal(providedCtx))
			})
		})

		Context("ls", func() {
			Context("w/ dirRef", func() {

				It("should call core.Ls w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedDirRef := "dummyPath"
					objectUnderTest := newCli(
						fakeCore,
						new(clicolorer.Fake),
						new(nodeCreateCmd.FakeInvoker),
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls", expectedDirRef})

					/* assert */
					actualCtx,
						actualDirRef := fakeCore.LsArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualDirRef).To(Equal(expectedDirRef))
				})
			})
			Context("w/out dirRef", func() {

				It("should call core.Ls w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedDirRef := ".opspec"
					objectUnderTest := newCli(
						fakeCore,
						new(clicolorer.Fake),
						new(nodeCreateCmd.FakeInvoker),
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls"})

					/* assert */
					actualCtx,
						actualDirRef := fakeCore.LsArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualDirRef).To(Equal(expectedDirRef))
				})
			})
		})

		Context("node", func() {

			Context("create", func() {

				It("should call core.NodeCreate w/ expected args", func() {
					/* arrange */
					fakeNodeCreateCmdInvoker := new(nodeCreateCmd.FakeInvoker)

					objectUnderTest := newCli(
						new(core.Fake),
						new(clicolorer.Fake),
						fakeNodeCreateCmdInvoker,
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "node", "create"})

					/* assert */
					Expect(fakeNodeCreateCmdInvoker.InvokeCallCount()).To(Equal(1))
				})

			})

			Context("kill", func() {

				It("should call core.NodeKill w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					objectUnderTest := newCli(
						fakeCore,
						new(clicolorer.Fake),
						new(nodeCreateCmd.FakeInvoker),
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "node", "kill"})

					/* assert */
					Expect(fakeCore.NodeKillCallCount()).To(Equal(1))
				})

			})
		})

		Context("op", func() {

			Context("create", func() {
				Context("w/ path", func() {
					It("should call core.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedOpName := "dummyOpName"
						expectedPath := "dummyPath"

						objectUnderTest := newCli(
							fakeCore,
							new(clicolorer.Fake),
							new(nodeCreateCmd.FakeInvoker),
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", "--path", expectedPath, expectedOpName})

						/* assert */
						Expect(fakeCore.OpCreateCallCount()).To(Equal(1))
						actualPath, actualOpDescription, actualOpName := fakeCore.OpCreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})

				Context("w/out path", func() {
					It("should call core.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedOpName := "dummyOpName"
						expectedPath := ".opspec"

						objectUnderTest := newCli(
							fakeCore,
							new(clicolorer.Fake),
							new(nodeCreateCmd.FakeInvoker),
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", expectedOpName})

						/* assert */
						Expect(fakeCore.OpCreateCallCount()).To(Equal(1))
						actualPath, actualOpDescription, actualOpName := fakeCore.OpCreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})
				Context("w/ description", func() {
					It("should call core.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedOpName := "dummyOpName"
						expectedOpDescription := "dummyOpDescription"
						expectedPath := ".opspec"

						objectUnderTest := newCli(
							fakeCore,
							new(clicolorer.Fake),
							new(nodeCreateCmd.FakeInvoker),
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", "-d", expectedOpDescription, expectedOpName})

						/* assert */
						Expect(fakeCore.OpCreateCallCount()).To(Equal(1))
						actualPath, actualOpDescription, actualOpName := fakeCore.OpCreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(Equal(expectedOpDescription))
						Expect(actualPath).To(Equal(expectedPath))
					})
				})

				Context("w/out description", func() {
					It("should call core.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(core.Fake)

						expectedName := "dummyOpName"
						expectedPath := ".opspec"

						objectUnderTest := newCli(
							fakeCore,
							new(clicolorer.Fake),
							new(nodeCreateCmd.FakeInvoker),
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", expectedName})

						/* assert */
						Expect(fakeCore.OpCreateCallCount()).To(Equal(1))
						actualPath, actualOpDescription, actualOpName := fakeCore.OpCreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})
			})

			Context("install", func() {
				It("should call core.Install w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedPath := "dummyPath"
					expectedOpRef := "dummyOpRef"
					expectedUsername := "dummyUsername"
					expectedPassword := "dummyPassword"

					objectUnderTest := newCli(
						fakeCore,
						new(clicolorer.Fake),
						new(nodeCreateCmd.FakeInvoker),
					)

					/* act */
					objectUnderTest.Run([]string{
						"opctl",
						"op",
						"install",
						"--path",
						expectedPath,
						"-u",
						expectedUsername,
						"-p",
						expectedPassword,
						expectedOpRef,
					})

					/* assert */
					actualCtx,
						actualPath,
						actualOpRef,
						actualUsername,
						actualPassword := fakeCore.OpInstallArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualPath).To(Equal(expectedPath))
					Expect(actualOpRef).To(Equal(expectedOpRef))
					Expect(actualUsername).To(Equal(expectedUsername))
					Expect(actualPassword).To(Equal(expectedPassword))
				})
			})

			Context("kill", func() {
				It("should call core.OpKill w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedOpID := "dummyOpID"

					objectUnderTest := newCli(
						fakeCore,
						new(clicolorer.Fake),
						new(nodeCreateCmd.FakeInvoker),
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "op", "kill", expectedOpID})

					/* assert */
					Expect(fakeCore.OpKillCallCount()).To(Equal(1))

					actualCtx,
						actualOpID := fakeCore.OpKillArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpID).To(Equal(expectedOpID))
				})
			})

			Context("validate", func() {

				It("should call core.OpValidate w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					opRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						fakeCore,
						new(clicolorer.Fake),
						new(nodeCreateCmd.FakeInvoker),
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "op", "validate", opRef})

					/* assert */
					actualCtx,
						actualOpRef := fakeCore.OpValidateArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpRef).To(Equal(opRef))
				})

			})

		})

		Context("run", func() {
			Context("with two op run args & an arg-file", func() {
				It("should call core.Run w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedRunOpts := &types.RunOpts{
						Args:    []string{"arg1Name=arg1Value", "arg2Name=arg2Value"},
						ArgFile: "dummyArgFile",
					}
					expectedOpRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						fakeCore,
						new(clicolorer.Fake),
						new(nodeCreateCmd.FakeInvoker),
					)

					/* act */
					objectUnderTest.Run([]string{
						"opctl",
						"run",
						"-a",
						expectedRunOpts.Args[0],
						"-a",
						expectedRunOpts.Args[1],
						"--arg-file",
						expectedRunOpts.ArgFile,
						expectedOpRef,
					})

					/* assert */
					actualCtx,
						actualOpUrl,
						actualRunOpts := fakeCore.RunArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpUrl).To(Equal(expectedOpRef))
					Expect(actualRunOpts).To(Equal(expectedRunOpts))
				})
			})

			Context("with zero op run args", func() {
				It("should call core.Run w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedOpRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						fakeCore,
						new(clicolorer.Fake),
						new(nodeCreateCmd.FakeInvoker),
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "run", expectedOpRef})

					/* assert */
					actualCtx,
						actualOpRef,
						actualRunOpts := fakeCore.RunArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpRef).To(Equal(expectedOpRef))
					Expect(actualRunOpts.Args).To(BeEmpty())
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

				objectUnderTest := newCli(
					fakeCore,
					new(clicolorer.Fake),
					new(nodeCreateCmd.FakeInvoker),
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "self-update", "-c", expectedChannel})

				/* assert */
				Expect(fakeCore.SelfUpdateCallCount()).To(Equal(1))

				actualChannel := fakeCore.SelfUpdateArgsForCall(0)
				Expect(actualChannel).To(Equal(expectedChannel))
			})
		})

		Context("without channel flag", func() {

			It("should call core.SelfUpdate with expected releaseChannel", func() {
				/* arrange */
				expectedChannel := "stable"

				fakeCore := new(core.Fake)

				objectUnderTest := newCli(
					fakeCore,
					new(clicolorer.Fake),
					new(nodeCreateCmd.FakeInvoker),
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "self-update"})

				/* assert */
				Expect(fakeCore.SelfUpdateCallCount()).To(Equal(1))

				actualChannel := fakeCore.SelfUpdateArgsForCall(0)
				Expect(actualChannel).To(Equal(expectedChannel))
			})
		})

	})

})
