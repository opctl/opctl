package containerd

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

// clearProxyEnv removes proxy env vars for the duration of a test so
// containerruntime.ProxyEnvVars (which reads the process env) doesn't make
// constructCreateArgs results depend on the ambient environment.
func clearProxyEnv() func() {
	names := []string{
		"HTTP_PROXY", "http_proxy",
		"HTTPS_PROXY", "https_proxy",
		"NO_PROXY", "no_proxy",
		"ALL_PROXY", "all_proxy",
	}
	saved := map[string]string{}
	for _, n := range names {
		if v, ok := os.LookupEnv(n); ok {
			saved[n] = v
			os.Unsetenv(n)
		}
	}
	return func() {
		for n, v := range saved {
			os.Setenv(n, v)
		}
	}
}

var _ = Context("constructCreateArgs", func() {
	strPtr := func(s string) *string { return &s }

	It("builds a deterministic, fully-populated arg vector", func() {
		restore := clearProxyEnv()
		defer restore()

		req := &model.ContainerCall{
			ContainerID: "cid",
			Name:        strPtr("myname"),
			WorkDir:     "/work",
			Cmd:         []string{"echo", "hi"},
			EnvVars:     map[string]string{"B": "2", "A": "1"},
			Files:       map[string]string{"/f2": "/hostf2", "/f1": "/hostf1"},
			Dirs:        map[string]string{"/d": "/hostd"},
			// only the unix socket (value containing a path separator) is mounted
			Sockets: map[string]string{"/sock": "/var/run/x.sock", "/tcp": "tcp:1234"},
			Ports:   map[string]string{"8080": "9090"},
			Image:   &model.ContainerCallImage{Ref: strPtr("alpine:3")},
		}

		args, err := constructCreateArgs(req, getContainerName(req.ContainerID))

		Expect(err).To(BeNil())
		Expect(args).To(Equal([]string{
			"create",
			"--name", "opctl_cid",
			"--network", "opctl",
			"--privileged",
			"--network-alias", "myname",
			"--workdir", "/work",
			"--env", "A=1",
			"--env", "B=2",
			"--volume", "/hostf1:/f1",
			"--volume", "/hostf2:/f2",
			"--volume", "/hostd:/d",
			"--volume", "/var/run/x.sock:/sock",
			"--publish", "9090:8080",
			"alpine:3",
			"echo", "hi",
		}))
	})

	It("omits optional flags when not set", func() {
		restore := clearProxyEnv()
		defer restore()

		req := &model.ContainerCall{
			ContainerID: "cid",
			Image:       &model.ContainerCallImage{Ref: strPtr("alpine")},
		}

		args, err := constructCreateArgs(req, getContainerName(req.ContainerID))

		Expect(err).To(BeNil())
		Expect(args).To(Equal([]string{
			"create",
			"--name", "opctl_cid",
			"--network", "opctl",
			"--privileged",
			"alpine",
		}))
	})

	It("propagates proxy env vars without clobbering op-declared ones", func() {
		restore := clearProxyEnv()
		defer restore()
		os.Setenv("HTTPS_PROXY", "http://proxy:3128")
		os.Setenv("NO_PROXY", "from-env")

		req := &model.ContainerCall{
			ContainerID: "cid",
			// op explicitly sets NO_PROXY; env value must not override it
			EnvVars: map[string]string{"NO_PROXY": "from-op"},
			Image:   &model.ContainerCallImage{Ref: strPtr("alpine")},
		}

		args, err := constructCreateArgs(req, getContainerName(req.ContainerID))

		Expect(err).To(BeNil())
		Expect(args).To(ContainElement("HTTPS_PROXY=http://proxy:3128"))
		Expect(args).To(ContainElement("NO_PROXY=from-op"))
		Expect(args).NotTo(ContainElement("NO_PROXY=from-env"))
	})

	It("errors when no image ref is provided", func() {
		req := &model.ContainerCall{ContainerID: "cid", Image: &model.ContainerCallImage{}}
		_, err := constructCreateArgs(req, getContainerName(req.ContainerID))
		Expect(err).To(HaveOccurred())
	})
})

var _ = Context("getContainerName", func() {
	It("prefixes the container id", func() {
		Expect(getContainerName("abc")).To(Equal("opctl_abc"))
	})
})

var _ = Context("lastNonEmptyLine", func() {
	It("returns the last non-empty trimmed line", func() {
		Expect(lastNonEmptyLine("garbage\n0\n\n")).To(Equal("0"))
		Expect(lastNonEmptyLine("137\n")).To(Equal("137"))
	})
})

var _ = Context("isNotFound", func() {
	It("detects nerdctl not-found output", func() {
		Expect(isNotFound("Error: no such container: opctl_x")).To(BeTrue())
		Expect(isNotFound("nerdctl: not found")).To(BeTrue())
		Expect(isNotFound("some other error")).To(BeFalse())
	})
})
