package euid0

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("sudoArgs", func() {
	providedEnviron := []string{
		"HOME=/home/user",
		"PATH=/usr/bin",
		"USER=user",
		"LOGNAME=user",
		"SECRET=a=b,c",
		"WEIRD,NAME=x",
		"=C:=C:\\",
	}
	providedBin := "/usr/local/bin/opctl"
	providedArgs := []string{"run", "-a", "x=y", "op"}

	Context("classic sudo", func() {
		It("should pass -E", func() {
			/* arrange/act/assert */
			Expect(sudoArgs(false, providedEnviron, providedBin, providedArgs)).
				To(Equal([]string{"-E", "/usr/local/bin/opctl", "run", "-a", "x=y", "op"}))
		})
	})
	Context("sudo-rs", func() {
		It("should name each preservable env var without its value", func() {
			/* arrange/act/assert */
			Expect(sudoArgs(true, providedEnviron, providedBin, providedArgs)).
				To(Equal([]string{"--preserve-env=HOME,SECRET", "/usr/local/bin/opctl", "run", "-a", "x=y", "op"}))
		})
	})
})
