package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/core"
	"github.com/opctl/opctl/util/clicolorer"
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
				Expect(fakeCliColorer.DisableCallCount()).To(Equal(1))
			})
		})

		Context("events", func() {
			It("should call core.Events w/ expected args", func() {
				/* arrange */
				fakeCore := new(core.Fake)

				objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

				/* act */
				objectUnderTest.Run([]string{"opctl", "events"})

				/* assert */
				Expect(fakeCore.EventsCallCount()).To(Equal(1))
			})
		})

		Context("ls", func() {
			Context("w/ path", func() {

				It("should call core.PkgLs w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedPath := "dummyPath"
					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls", "-c", expectedPath})

					/* assert */
					Expect(fakeCore.PkgLsCallCount()).To(Equal(1))
					actualPath := fakeCore.PkgLsArgsForCall(0)
					Expect(actualPath).To(Equal(expectedPath))
				})
			})
			Context("w/out path", func() {

				It("should call core.PkgLs w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedPath := ".opspec"
					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "ls"})

					/* assert */
					Expect(fakeCore.PkgLsCallCount()).To(Equal(1))
					actualPath := fakeCore.PkgLsArgsForCall(0)
					Expect(actualPath).To(Equal(expectedPath))
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
					Expect(fakeCore.NodeCreateCallCount()).To(Equal(1))
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
					Expect(fakeCore.NodeKillCallCount()).To(Equal(1))
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
					Expect(fakeCore.OpKillCallCount()).To(Equal(1))
					Expect(fakeCore.OpKillArgsForCall(0)).To(Equal(expectedOpId))
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
						Expect(fakeCore.PkgCreateCallCount()).To(Equal(1))
						actualPath, actualPkgDescription, actualPkgName := fakeCore.PkgCreateArgsForCall(0)
						Expect(actualPkgName).To(Equal(expectedPkgName))
						Expect(actualPkgDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
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
						Expect(fakeCore.PkgCreateCallCount()).To(Equal(1))
						actualPath, actualPkgDescription, actualPkgName := fakeCore.PkgCreateArgsForCall(0)
						Expect(actualPkgName).To(Equal(expectedPkgName))
						Expect(actualPkgDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
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
						Expect(fakeCore.PkgCreateCallCount()).To(Equal(1))
						actualPath, actualPkgDescription, actualPkgName := fakeCore.PkgCreateArgsForCall(0)
						Expect(actualPkgName).To(Equal(expectedPkgName))
						Expect(actualPkgDescription).To(Equal(expectedPkgDescription))
						Expect(actualPath).To(Equal(expectedPath))
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
						Expect(fakeCore.PkgCreateCallCount()).To(Equal(1))
						actualPath, actualPkgDescription, actualPkgName := fakeCore.PkgCreateArgsForCall(0)
						Expect(actualPkgName).To(Equal(expectedName))
						Expect(actualPkgDescription).To(BeEmpty())
						Expect(actualPath).To(Equal(expectedPath))
					})
				})
			})

			Context("pull", func() {
				It("should call core.Pull w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedPkgRef := "dummyPkgRef"
					expectedUsername := "dummyUsername"
					expectedPassword := "dummyPassword"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{
						"opctl",
						"pkg",
						"pull",
						"-u",
						expectedUsername,
						"-p",
						expectedPassword,
						expectedPkgRef,
					})

					/* assert */
					actualPkgRef, actualUsername, actualPassword := fakeCore.PkgPullArgsForCall(0)
					Expect(actualPkgRef).To(Equal(expectedPkgRef))
					Expect(actualUsername).To(Equal(expectedUsername))
					Expect(actualPassword).To(Equal(expectedPassword))
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
					actualPkgRef := fakeCore.PkgValidateArgsForCall(0)
					Expect(actualPkgRef).To(Equal(pkgRef))
				})

			})

		})

		Context("run", func() {
			Context("with two op run args & an arg-file", func() {
				It("should call core.Run w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedRunOpts := &core.RunOpts{
						Args:    []string{"arg1Name=arg1Value", "arg2Name=arg2Value"},
						ArgFile: "dummyArgFile",
					}
					expectedPkgRef := ".opspec/dummyPkgName"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

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
						expectedPkgRef,
					})

					/* assert */
					Expect(fakeCore.RunCallCount()).To(Equal(1))
					actualOpUrl, actualRunOpts := fakeCore.RunArgsForCall(0)
					Expect(actualOpUrl).To(Equal(expectedPkgRef))
					Expect(actualRunOpts).To(Equal(expectedRunOpts))
				})
			})

			Context("with zero op run args", func() {
				It("should call core.Run w/ expected args", func() {
					/* arrange */
					fakeCore := new(core.Fake)

					expectedPkgRef := ".opspec/dummyPkgName"

					objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

					/* act */
					objectUnderTest.Run([]string{"opctl", "run", expectedPkgRef})

					/* assert */
					Expect(fakeCore.RunCallCount()).To(Equal(1))

					actualPkgRef, actualRunOpts := fakeCore.RunArgsForCall(0)
					Expect(actualPkgRef).To(Equal(expectedPkgRef))
					Expect(actualRunOpts).To(BeEmpty())
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

				objectUnderTest := newCli(fakeCore, new(clicolorer.Fake))

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
