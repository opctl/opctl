package main

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cliColorerFakes "github.com/opctl/opctl/cli/internal/clicolorer/fakes"
	coreFakes "github.com/opctl/opctl/cli/internal/core/fakes"
	nodeFakes "github.com/opctl/opctl/cli/internal/core/node/fakes"
	opFakes "github.com/opctl/opctl/cli/internal/core/op/fakes"
	"github.com/opctl/opctl/cli/internal/model"
)

var _ = Context("cli", func() {
	Context("Run", func() {

		Context("--no-color", func() {
			It("should set color.NoColor", func() {
				/* arrange */
				fakeCliColorer := new(cliColorerFakes.FakeCliColorer)

				objectUnderTest := newCli(
					fakeCliColorer,
					new(coreFakes.FakeCore),
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "--no-color", "ls"})

				/* assert */
				Expect(fakeCliColorer.DisableCallCount()).To(Equal(1))
			})
		})

		Context("events", func() {
			It("should call coreFakes.Events w/ expected args", func() {
				/* arrange */
				providedCtx := context.Background()
				fakeCore := new(coreFakes.FakeCore)

				objectUnderTest := newCli(
					new(cliColorerFakes.FakeCliColorer),
					fakeCore,
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

				It("should call coreFakes.Ls w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					expectedDirRef := "dummyPath"
					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
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

				It("should call coreFakes.Ls w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					expectedDirRef := ".opspec"
					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
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

				It("should call nodeCreateCmd.Invoke w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					fakeNode := new(nodeFakes.FakeNode)
					fakeCore.NodeReturns(fakeNode)

					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "node", "create"})

					/* assert */
					Expect(fakeNode.CreateCallCount()).To(Equal(1))
				})

			})

			Context("kill", func() {

				It("should call nodeKillCmd.Invoke w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					fakeNode := new(nodeFakes.FakeNode)
					fakeCore.NodeReturns(fakeNode)

					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "node", "kill"})

					/* assert */
					Expect(fakeNode.KillCallCount()).To(Equal(1))
				})

			})
		})

		Context("op", func() {

			Context("create", func() {
				Context("w/ path", func() {
					It("should call coreFakes.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(coreFakes.FakeCore)

						fakeOp := new(opFakes.FakeOp)
						fakeCore.OpReturns(fakeOp)

						expectedOpName := "dummyOpName"
						expectedPath := "dummyPath"

						objectUnderTest := newCli(
							new(cliColorerFakes.FakeCliColorer),
							fakeCore,
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", "--path", expectedPath, expectedOpName})

						/* assert */
						actualPath, actualOpDescription, actualOpName := fakeOp.CreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})

				Context("w/out path", func() {
					It("should call coreFakes.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(coreFakes.FakeCore)

						fakeOp := new(opFakes.FakeOp)
						fakeCore.OpReturns(fakeOp)

						expectedOpName := "dummyOpName"
						expectedPath := ".opspec"

						objectUnderTest := newCli(
							new(cliColorerFakes.FakeCliColorer),
							fakeCore,
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", expectedOpName})

						/* assert */
						actualPath, actualOpDescription, actualOpName := fakeOp.CreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})
				Context("w/ description", func() {
					It("should call coreFakes.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(coreFakes.FakeCore)

						fakeOp := new(opFakes.FakeOp)
						fakeCore.OpReturns(fakeOp)

						expectedOpName := "dummyOpName"
						expectedOpDescription := "dummyOpDescription"
						expectedPath := ".opspec"

						objectUnderTest := newCli(
							new(cliColorerFakes.FakeCliColorer),
							fakeCore,
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", "-d", expectedOpDescription, expectedOpName})

						/* assert */
						actualPath, actualOpDescription, actualOpName := fakeOp.CreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedOpName))
						Expect(actualOpDescription).To(Equal(expectedOpDescription))
						Expect(actualPath).To(Equal(expectedPath))
					})
				})

				Context("w/out description", func() {
					It("should call coreFakes.Create w/ expected args", func() {
						/* arrange */
						fakeCore := new(coreFakes.FakeCore)

						fakeOp := new(opFakes.FakeOp)
						fakeCore.OpReturns(fakeOp)

						expectedName := "dummyOpName"
						expectedPath := ".opspec"

						objectUnderTest := newCli(
							new(cliColorerFakes.FakeCliColorer),
							fakeCore,
						)

						/* act */
						objectUnderTest.Run([]string{"opctl", "op", "create", expectedName})

						/* assert */
						actualPath, actualOpDescription, actualOpName := fakeOp.CreateArgsForCall(0)
						Expect(actualOpName).To(Equal(expectedName))
						Expect(actualOpDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})
			})

			Context("install", func() {
				It("should call coreFakes.Install w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					fakeOp := new(opFakes.FakeOp)
					fakeCore.OpReturns(fakeOp)

					expectedPath := "dummyPath"
					expectedOpRef := "dummyOpRef"
					expectedUsername := "dummyUsername"
					expectedPassword := "dummyPassword"

					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
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
						actualPassword := fakeOp.InstallArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualPath).To(Equal(expectedPath))
					Expect(actualOpRef).To(Equal(expectedOpRef))
					Expect(actualUsername).To(Equal(expectedUsername))
					Expect(actualPassword).To(Equal(expectedPassword))
				})
			})

			Context("kill", func() {
				It("should call coreFakes.OpKill w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					fakeOp := new(opFakes.FakeOp)
					fakeCore.OpReturns(fakeOp)

					expectedOpID := "dummyOpID"

					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "op", "kill", expectedOpID})

					/* assert */
					actualCtx,
						actualOpID := fakeOp.KillArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpID).To(Equal(expectedOpID))
				})
			})

			Context("validate", func() {

				It("should call coreFakes.OpValidate w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					fakeOp := new(opFakes.FakeOp)
					fakeCore.OpReturns(fakeOp)

					opRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
					)

					/* act */
					objectUnderTest.Run([]string{"opctl", "op", "validate", opRef})

					/* assert */
					actualCtx,
						actualOpRef := fakeOp.ValidateArgsForCall(0)

					Expect(actualCtx).To(Equal(context.TODO()))
					Expect(actualOpRef).To(Equal(opRef))
				})

			})

		})

		Context("run", func() {
			Context("with two op run args & an arg-file", func() {
				It("should call coreFakes.Run w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					expectedRunOpts := &model.RunOpts{
						Args:    []string{"arg1Name=arg1Value", "arg2Name=arg2Value"},
						ArgFile: "dummyArgFile",
					}
					expectedOpRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
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
				It("should call coreFakes.Run w/ expected args", func() {
					/* arrange */
					fakeCore := new(coreFakes.FakeCore)

					expectedOpRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						new(cliColorerFakes.FakeCliColorer),
						fakeCore,
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

			It("should call coreFakes.SelfUpdate with expected releaseChannel", func() {
				/* arrange */
				expectedChannel := "beta"

				fakeCore := new(coreFakes.FakeCore)

				objectUnderTest := newCli(
					new(cliColorerFakes.FakeCliColorer),
					fakeCore,
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "self-update", "-c", expectedChannel})

				/* assert */
				actualChannel := fakeCore.SelfUpdateArgsForCall(0)
				Expect(actualChannel).To(Equal(expectedChannel))
			})

		})

		Context("without channel flag", func() {

			It("should call coreFakes.SelfUpdate with expected releaseChannel", func() {
				/* arrange */
				expectedChannel := "stable"

				fakeCore := new(coreFakes.FakeCore)

				objectUnderTest := newCli(
					new(cliColorerFakes.FakeCliColorer),
					fakeCore,
				)

				/* act */
				objectUnderTest.Run([]string{"opctl", "self-update"})

				/* assert */
				actualChannel := fakeCore.SelfUpdateArgsForCall(0)
				Expect(actualChannel).To(Equal(expectedChannel))
			})
		})

	})

})
