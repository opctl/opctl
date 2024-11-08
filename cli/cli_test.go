package main

import (
	"io"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var pathToOpctl string

var _ = BeforeSuite(func() {
	compiledPath, err := gexec.Build("./", "-buildvcs=false")
	if err != nil {
		panic(err)
	}

	pathToOpctl = filepath.Join(compiledPath, "cli")

	// start node
	command := exec.Command(pathToOpctl, "node", "create")
	if _, err := gexec.Start(command, io.Discard, io.Discard); err != nil {
		panic(err)
	}
})

var _ = Context("cli", func() {
	Context("--no-color", func() {
		It("should not err", func() {
			/* arrange */
			command := exec.Command(pathToOpctl, "--no-color", "ls")

			/* act */
			session, actualErr := gexec.Start(command, io.Discard, io.Discard)

			/* assert */
			Expect(actualErr).NotTo(HaveOccurred())
			Eventually(session, 10).Should(gexec.Exit(0))
		})
	})

	Context("auth", func() {

		Context("add", func() {

			It("should not err", func() {
				/* arrange */
				providedResources := "resources"
				providedUsername := "username"
				providedPassword := "password"
				command := exec.Command(pathToOpctl, "auth", "add", providedResources, "-u", providedUsername, "-p", providedPassword)

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})

		})

	})

	Context("events", func() {
		It("should not err", func() {
			/* arrange */
			command := exec.Command(pathToOpctl, "events")

			/* act */
			session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			session.Interrupt()

			/* assert */
			Expect(actualErr).NotTo(HaveOccurred())
			Eventually(session, 10).Should(gexec.Exit(130))
		})
	})

	Context("ls", func() {
		Context("w/ dirRef", func() {
			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "ls", "testdata/ls")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
				Expect(string(session.Out.Contents())).Should(Equal(
					`REF		DESCRIPTION
testdata/ls/op1	A single line description
`))
			})
		})
		Context("w/out dirRef", func() {

			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "ls")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})
		})
	})

	// disable for now since it will kill the running test container
	XContext("node", Label("Ordered"), func() {

		Context("create", func() {

			It("should not err", func() {
				/* arrange */
				// ensure no node running
				err := exec.Command(pathToOpctl, "node", "delete").Run()
				if err != nil {
					panic(err)
				}

				command := exec.Command(pathToOpctl, "node", "create")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				session.Interrupt()

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(130))
			})

		})

		Context("delete", func() {

			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "node", "delete")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})

		})

		Context("kill", Label("Serial"), func() {

			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "node", "kill")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})

		})
	})

	Context("op", func() {

		Context("create", func() {
			Context("w/ path", func() {
				It("should not err", func() {
					/* arrange */
					command := exec.Command(pathToOpctl, "op", "create", "--path", "/tmp", "withPath")

					/* act */
					session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

					/* assert */
					Expect(actualErr).NotTo(HaveOccurred())
					Eventually(session, 10).Should(gexec.Exit(0))
				})
			})

			Context("w/out path", func() {
				It("should not err", func() {
					/* arrange */
					command := exec.Command(pathToOpctl, "op", "create", "withoutPath")

					/* act */
					session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

					/* assert */
					Expect(actualErr).NotTo(HaveOccurred())
					Eventually(session, 10).Should(gexec.Exit(0))
				})
			})
			Context("w/ description", func() {
				It("should not err", func() {
					/* arrange */
					command := exec.Command(pathToOpctl, "op", "create", "--path", "/tmp", "-d", "dummyOpDescription", "withDescription")

					/* act */
					session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

					/* assert */
					Expect(actualErr).NotTo(HaveOccurred())
					Eventually(session, 10).Should(gexec.Exit(0))
				})
			})

			Context("w/out description", func() {
				It("should not err", func() {
					/* arrange */
					command := exec.Command(pathToOpctl, "op", "create", "--path", "/tmp", "withoutDescription")

					/* act */
					session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

					/* assert */
					Expect(actualErr).NotTo(HaveOccurred())
					Eventually(session, 10).Should(gexec.Exit(0))
				})
			})
		})

		Context("install", func() {
			It("should not err", func() {
				/* arrange */
				command := exec.Command(
					pathToOpctl,
					"op",
					"install",
					"--path",
					"/tmp/twoArgsCopy",
					"./testdata/twoArgs",
				)

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})
		})

		Context("kill", func() {
			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "op", "kill", "xxx")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})
		})

		Context("validate", func() {

			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "op", "validate", "./testdata/zeroArgs")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})

		})

	})

	Context("run", func() {
		Context("with two args", func() {
			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "run", "-a", "arg1=value1", "-a", "arg2=value2", "./testdata/twoArgs")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})
		})

		Context("with zero args", func() {
			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "run", "./testdata/zeroArgs")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(0))
			})
		})
	})

	Context("self-update", func() {

		It("should not err", func() {
			/* arrange */
			command := exec.Command(pathToOpctl, "self-update")

			/* act */
			session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

			/* assert */
			Expect(actualErr).NotTo(HaveOccurred())
			Eventually(session, 10).Should(gexec.Exit(1))
		})

	})

	Context("ui", func() {
		Context("w/ mountRef", func() {
			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "ui", "../.opspec/build")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(1))
			})
		})
		Context("w/out mountRef", func() {
			It("should not err", func() {
				/* arrange */
				command := exec.Command(pathToOpctl, "ui")

				/* act */
				session, actualErr := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				/* assert */
				Expect(actualErr).NotTo(HaveOccurred())
				Eventually(session, 10).Should(gexec.Exit(1))
			})
		})
	})

})
