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
	It("scaffolds a Go op with the in-repo SDK and runs it", func() {
		runner, err := filepath.Abs("../cmd.sh")
		Expect(err).NotTo(HaveOccurred())
		dir, err := os.MkdirTemp("", "scenario-runner-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		write := func(name, text string) {
			Expect(os.WriteFile(filepath.Join(dir, name), []byte(text), 0700)).To(Succeed())
		}
		write("dockerd", "#!/bin/sh\nexit 0\n")
		write("opctl", `#!/bin/sh
echo "$@" >> "$OPCTL_CALLS"
if [ "$1" = op ]; then mkdir -p "$5/$8" && echo "module $8" > "$5/$8/go.mod"; fi
exit 0
`)
		command := exec.Command("sh", runner)
		command.Env = append(os.Environ(), "PATH="+dir+":"+os.Getenv("PATH"), "authAddGithub=false", "expect=success", "op=op/create/go", "createGo=hello", "createParent="+dir, "sdkPath=/sdk", "OPCTL_CALLS="+filepath.Join(dir, "calls"), "expectedOutputsPath="+filepath.Join(dir, "none.json"))
		output, err := command.CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), string(output))
		calls, err := os.ReadFile(filepath.Join(dir, "calls"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(calls)).To(ContainSubstring("op create --go -p " + dir + " -d Greets the caller. hello\n"))
		Expect(string(calls)).To(HaveSuffix("run --no-progress --arg-file /args.yml " + dir + "/hello\n"))
		goMod, err := os.ReadFile(filepath.Join(dir, "hello", "go.mod"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(goMod)).To(Equal("module hello\n\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => /sdk\n"))
	})
})
