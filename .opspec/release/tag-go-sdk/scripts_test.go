package taggosdk

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const commit = "0123456789abcdef0123456789abcdef01234567"

// fakeCurl answers the tag lookup with FAKE_STATUS and FAKE_BODY and records each call.
const fakeCurl = `#!/bin/sh
echo "$@" >> "$CALLS"
out=
while [ $# -gt 0 ]; do
  if [ "$1" = -o ]; then out=$2; shift; fi
  shift
done
if [ -n "$out" ]; then printf '%s' "$FAKE_BODY" > "$out"; printf '%s' "$FAKE_STATUS"; fi
`

// fakeGo prints the JSON that go mod download prints for a module version.
const fakeGo = `#!/bin/sh
echo "$@" >> "$CALLS"
cat <<'JSON'
{
	"Path": "github.com/opctl/opctl/sdks/go/opspec/sdk",
	"Version": "v0.1.77",
	"Sum": "h1:zip+sum/abc=",
	"GoModSum": "h1:mod+sum/def="
}
JSON
`

var _ = Describe("tagging the Go op SDK", func() {
	var dir string
	BeforeEach(func() {
		var err error
		dir, err = os.MkdirTemp("", "tag-go-sdk-")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(filepath.Join(dir, "curl"), []byte(fakeCurl), 0700)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(dir, "go"), []byte(fakeGo), 0700)).To(Succeed())
	})
	AfterEach(func() { Expect(os.RemoveAll(dir)).To(Succeed()) })
	run := func(script, status, body string) (string, string, error) {
		path, err := filepath.Abs(script)
		Expect(err).NotTo(HaveOccurred())
		cmd := exec.Command("sh", path)
		cmd.Env = append(os.Environ(), "PATH="+dir+":"+os.Getenv("PATH"), "CALLS="+filepath.Join(dir, "calls"), "FAKE_STATUS="+status, "FAKE_BODY="+body,
			"version=0.1.77", "commit="+commit, "githubAccessToken=token", "outDir="+dir, "TMPDIR="+dir)
		output, err := cmd.CombinedOutput()
		calls, _ := os.ReadFile(filepath.Join(dir, "calls"))
		return string(output), string(calls), err
	}

	It("creates a missing tag at the release commit", func() {
		output, calls, err := run("tag.sh", "404", `{"message":"Not Found"}`)
		Expect(err).NotTo(HaveOccurred(), output)
		Expect(calls).To(ContainSubstring("https://api.github.com/repos/opctl/opctl/git/ref/tags/sdks/go/opspec/sdk/v0.1.77"))
		Expect(calls).To(ContainSubstring(`-X POST`))
		Expect(calls).To(ContainSubstring(`{"ref":"refs/tags/sdks/go/opspec/sdk/v0.1.77","sha":"` + commit + `"}`))
	})

	It("keeps a tag that already points at the release commit", func() {
		output, calls, err := run("tag.sh", "200", `{"ref":"refs/tags/sdks/go/opspec/sdk/v0.1.77","object":{"sha": "`+commit+`","type":"commit"}}`)
		Expect(err).NotTo(HaveOccurred(), output)
		Expect(calls).NotTo(ContainSubstring("POST"))
	})

	It("fails when the tag points elsewhere or can't be read", func() {
		for status, body := range map[string]string{
			"200": `{"object":{"sha":"fedcba9876543210fedcba9876543210fedcba98","type":"commit"}}`,
			"500": `{"message":"error"}`,
		} {
			Expect(os.RemoveAll(filepath.Join(dir, "calls"))).To(Succeed())
			output, calls, err := run("tag.sh", status, body)
			Expect(err).To(HaveOccurred(), output)
			Expect(calls).NotTo(ContainSubstring("POST"))
		}
	})

	It("reads the released module's go.sum hashes", func() {
		output, calls, err := run("sums.sh", "", "")
		Expect(err).NotTo(HaveOccurred(), output)
		Expect(calls).To(Equal("mod download -json github.com/opctl/opctl/sdks/go/opspec/sdk@v0.1.77\n"))
		for name, sum := range map[string]string{"sdkSum": "h1:zip+sum/abc=", "sdkGoModSum": "h1:mod+sum/def="} {
			data, err := os.ReadFile(filepath.Join(dir, name))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(Equal(sum))
		}
	})
})
