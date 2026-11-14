package euid0

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("sudoArgs", func() {
	Context("classic sudo", func() {
		It("should pass -E", func() {
			/* arrange */
			objectUnderTest := sudoArgs
			providedEnviron := []string{"HOME=/home/user", "SECRET=s3cret"}
			providedArgs := []string{"run", "-a", "x=y", "op"}

			expectedArgs := []string{"-E", "/usr/local/bin/opctl", "run", "-a", "x=y", "op"}

			/* act */
			actualArgs := objectUnderTest(false, providedEnviron, "/usr/local/bin/opctl", providedArgs)

			/* assert */
			Expect(actualArgs).To(Equal(expectedArgs))
		})
	})
	Context("sudo-rs", func() {
		It("should name each preservable env var without its value", func() {
			/* arrange */
			objectUnderTest := sudoArgs
			providedEnviron := []string{
				"HOME=/home/user",
				"PATH=/usr/bin",
				"USER=user",
				"LOGNAME=user",
				"SUDO_USER=root",
				"SUDO_UID=0",
				"SUDO_GID=0",
				"SUDO_HOME=/root",
				"SUDO_COMMAND=/usr/bin/sh",
				"SUDO_ASKPASS=/usr/bin/askpass",
				"SECRET=a=b,c",
				"WEIRD,NAME=x",
				"BASH_FUNC_f%%=() { echo hi; }",
				"=C:=C:\\",
			}
			providedArgs := []string{"run", "-a", "x=y", "op"}

			expectedArgs := []string{"--preserve-env=HOME,SUDO_ASKPASS,SECRET", "/usr/local/bin/opctl", "run", "-a", "x=y", "op"}

			/* act */
			actualArgs := objectUnderTest(true, providedEnviron, "/usr/local/bin/opctl", providedArgs)

			/* assert */
			Expect(actualArgs).To(Equal(expectedArgs))
		})
	})
})
