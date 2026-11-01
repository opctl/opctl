package main

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("scenario runner output assertions", func() {
	It("preserves CLI status and checks outputs only after expected success", func() {
		runner, err := filepath.Abs("../cmd.sh")
		Expect(err).NotTo(HaveOccurred())
		for _, tc := range []struct {
			name, expected, status, expectation, assertStatus string
			success, called                                   bool
		}{
			{"old success", "success", "0", "", "0", true, false},
			{"empty object assertion", "success", "0", "{}", "0", true, true},
			{"assert failure", "success", "0", "{}", "1", false, true},
			{"expected CLI failure", "failure", "1", "{}", "0", true, false},
			{"unexpected CLI failure", "success", "1", "{}", "0", false, false},
			{"unexpected success", "failure", "0", "", "0", false, false},
		} {
			dir, err := os.MkdirTemp("", "scenario-runner-")
			Expect(err).NotTo(HaveOccurred())
			write := func(name, text string) {
				Expect(os.WriteFile(filepath.Join(dir, name), []byte(text), 0700)).To(Succeed())
			}
			write("dockerd", "#!/bin/sh\nexit 0\n")
			write("opctl", "#!/bin/sh\nif [ \"$1\" = run ]; then exit \"$FAKE_RUN_STATUS\"; fi\nexit 0\n")
			write("assert-outputs", "#!/bin/sh\nprintf called > \"$ASSERT_MARKER\"\nexit \"$FAKE_ASSERT_STATUS\"\n")
			write("expected.json", tc.expectation)
			command := exec.Command("sh", runner)
			command.Env = append(os.Environ(), "PATH="+dir+":"+os.Getenv("PATH"), "authAddGithub=false", "expect="+tc.expected, "op=test", "opRef=/test", "FAKE_RUN_STATUS="+tc.status, "FAKE_ASSERT_STATUS="+tc.assertStatus, "ASSERT_MARKER="+filepath.Join(dir, "called"), "expectedOutputsPath="+filepath.Join(dir, "expected.json"), "assertOutputsPath="+filepath.Join(dir, "assert-outputs"))
			output, err := command.CombinedOutput()
			Expect(err == nil).To(Equal(tc.success), "%s: %s", tc.name, output)
			_, err = os.Stat(filepath.Join(dir, "called"))
			Expect(err == nil).To(Equal(tc.called), tc.name)
			Expect(os.RemoveAll(dir)).To(Succeed())
		}
	})
})
