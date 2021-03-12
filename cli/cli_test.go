package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"os"
)

var _ = Context("cli", func() {
	Context("Run", func() {
		// @TODO: the below is not really testing anything but the test scenarios are good.
		// We need to move to remove testModeEnvVar and implement something like gexec (http://onsi.github.io/gomega/#gexec-testing-external-processes)
		// to properly test the CLI otherwise CLI exit codes and stdin/stderr/stdout reads/writes interfere w/ the test harness
		os.Setenv(testModeEnvVar, "")

		cliOutput := clioutput.New(clicolorer.New(), os.Stderr, os.Stdout)

		Context("--no-color", func() {
			It("should not err", func() {
				/* arrange */
				objectUnderTest := newCli(
					cliOutput,
				)

				/* act */
				actualErr := objectUnderTest.Run([]string{"opctl", "--no-color", "ls"})

				/* assert */
				Expect(actualErr).To(BeNil())
			})
		})

		Context("auth", func() {

			Context("add", func() {

				It("should not err", func() {
					/* arrange */
					providedResources := "resources"
					providedUsername := "username"
					providedPassword := "password"

					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "auth", "add", providedResources, "-u", providedUsername, "-p", providedPassword})

					/* assert */
					Expect(actualErr).To(BeNil())
				})

			})

		})

		Context("events", func() {
			It("should not err", func() {
				/* arrange */
				objectUnderTest := newCli(
					cliOutput,
				)

				/* act */
				actualErr := objectUnderTest.Run([]string{"opctl", "events"})

				/* assert */
				Expect(actualErr).To(BeNil())
			})
		})

		Context("ls", func() {
			Context("w/ dirRef", func() {

				It("should not err", func() {
					/* arrange */
					expectedDirRef := "dummyPath"
					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "ls", expectedDirRef})

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})
			Context("w/out dirRef", func() {

				It("should not err", func() {
					/* arrange */
					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "ls"})

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})
		})

		Context("node", func() {

			Context("create", func() {

				It("should not err", func() {
					/* arrange */
					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "node", "create"})

					/* assert */
					Expect(actualErr).To(BeNil())
				})

			})

			Context("kill", func() {

				It("should not err", func() {
					/* arrange */
					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "node", "kill"})

					/* assert */
					Expect(actualErr).To(BeNil())
				})

			})
		})

		Context("op", func() {

			Context("create", func() {
				Context("w/ path", func() {
					It("should not err", func() {
						/* arrange */
						expectedOpName := "dummyOpName"
						expectedPath := "dummyPath"

						objectUnderTest := newCli(
							cliOutput,
						)

						/* act */
						actualErr := objectUnderTest.Run([]string{"opctl", "op", "create", "--path", expectedPath, expectedOpName})

						/* assert */
						Expect(actualErr).To(BeNil())
					})
				})

				Context("w/out path", func() {
					It("should not err", func() {
						/* arrange */
						expectedOpName := "dummyOpName"

						objectUnderTest := newCli(
							cliOutput,
						)

						/* act */
						actualErr := objectUnderTest.Run([]string{"opctl", "op", "create", expectedOpName})

						/* assert */
						Expect(actualErr).To(BeNil())
					})
				})
				Context("w/ description", func() {
					It("should not err", func() {
						/* arrange */
						expectedOpName := "dummyOpName"
						expectedOpDescription := "dummyOpDescription"

						objectUnderTest := newCli(
							cliOutput,
						)

						/* act */
						actualErr := objectUnderTest.Run([]string{"opctl", "op", "create", "-d", expectedOpDescription, expectedOpName})

						/* assert */
						Expect(actualErr).To(BeNil())
					})
				})

				Context("w/out description", func() {
					It("should not err", func() {
						/* arrange */
						expectedName := "dummyOpName"

						objectUnderTest := newCli(
							cliOutput,
						)

						/* act */
						actualErr := objectUnderTest.Run([]string{"opctl", "op", "create", expectedName})

						/* assert */
						Expect(actualErr).To(BeNil())
					})
				})
			})

			Context("install", func() {
				It("should not err", func() {
					/* arrange */
					expectedPath := "dummyPath"
					expectedOpRef := "dummyOpRef"
					expectedUsername := "dummyUsername"
					expectedPassword := "dummyPassword"

					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{
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
					Expect(actualErr).To(BeNil())
				})
			})

			Context("kill", func() {
				It("should not err", func() {
					/* arrange */
					expectedOpID := "dummyOpID"

					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "op", "kill", expectedOpID})

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})

			Context("validate", func() {

				It("should not err", func() {
					/* arrange */
					opRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "op", "validate", opRef})

					/* assert */
					Expect(actualErr).To(BeNil())
				})

			})

		})

		Context("run", func() {
			Context("with two op run args & an arg-file", func() {
				It("should not err", func() {
					/* arrange */
					providedArgs := []string{"arg1Name=arg1Value", "arg2Name=arg2Value"}
					providedArgFile := "dummyArgFile"
					expectedOpRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{
						"opctl",
						"run",
						"-a",
						providedArgs[0],
						"-a",
						providedArgs[1],
						"--arg-file",
						providedArgFile,
						expectedOpRef,
					})

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})

			Context("with zero op run args", func() {
				It("should not err", func() {
					/* arrange */
					expectedOpRef := ".opspec/dummyOpName"

					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "run", expectedOpRef})

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})
		})

		Context("self-update", func() {

			It("should not err", func() {
				/* arrange */
				objectUnderTest := newCli(
					cliOutput,
				)

				/* act */
				actualErr := objectUnderTest.Run([]string{"opctl", "self-update"})

				/* assert */
				Expect(actualErr).To(BeNil())
			})

		})

		Context("ui", func() {
			Context("w/ mountRef", func() {

				It("should not err", func() {
					/* arrange */
					expectedDirRef := "./dummyPath"
					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "ui", expectedDirRef})

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})
			Context("w/out mountRef", func() {

				It("should not err", func() {
					/* arrange */
					objectUnderTest := newCli(
						cliOutput,
					)

					/* act */
					actualErr := objectUnderTest.Run([]string{"opctl", "ui"})

					/* assert */
					Expect(actualErr).To(BeNil())
				})
			})
		})
	})

})
