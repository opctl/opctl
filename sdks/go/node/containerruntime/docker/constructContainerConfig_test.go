package docker

import (
	"fmt"
	"os"
	"sort"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("constructContainerConfig", func() {
	It("should propagate proxy env vars from the node env without overriding op values", func() {
		/* arrange */
		os.Setenv("HTTP_PROXY", "http://node-proxy:3128")
		os.Setenv("NO_PROXY", "localhost")
		defer os.Unsetenv("HTTP_PROXY")
		defer os.Unsetenv("NO_PROXY")

		/* act */
		actualResult := constructContainerConfig(
			[]string{"dummyCmd"},
			map[string]string{
				// op explicitly sets HTTP_PROXY; it must win over the node env
				"HTTP_PROXY": "http://op-proxy:8080",
			},
			"dummyImageRef",
			nat.PortMap{},
			"dummyWorkDir",
		)

		/* assert */
		Expect(actualResult.Env).To(ContainElement("HTTP_PROXY=http://op-proxy:8080"))
		Expect(actualResult.Env).NotTo(ContainElement("HTTP_PROXY=http://node-proxy:3128"))
		Expect(actualResult.Env).To(ContainElement("NO_PROXY=localhost"))
	})

	It("should return expected result", func() {

		/* arrange */
		// ensure the host's proxy env can't perturb the deterministic expectation
		for _, name := range []string{
			"HTTP_PROXY", "http_proxy",
			"HTTPS_PROXY", "https_proxy",
			"NO_PROXY", "no_proxy",
			"ALL_PROXY", "all_proxy",
		} {
			if prev, had := os.LookupEnv(name); had {
				os.Unsetenv(name)
				defer os.Setenv(name, prev)
			}
		}

		providedCmd := []string{
			"dummyCmd",
		}
		providedEnvVars := map[string]string{
			"envVar1Name": "envVar1Value",
			"envVar2Name": "envVar2Value",
			"envVar3Name": "envVar3Value",
		}
		providedImageRef := "dummyImageRef"
		providedPortBindings := nat.PortMap{
			"80/tcp":   []nat.PortBinding{},
			"6060/udp": []nat.PortBinding{},
		}
		providedWorkDir := "dummyWorkDir"

		expectedResult := &container.Config{
			Entrypoint:   providedCmd,
			Env:          []string{},
			ExposedPorts: nat.PortSet{},
			Image:        providedImageRef,
			WorkingDir:   providedWorkDir,
			Tty:          true,
		}

		for port := range providedPortBindings {
			expectedResult.ExposedPorts[port] = struct{}{}
		}
		for envVarName, envVarValue := range providedEnvVars {
			expectedResult.Env = append(expectedResult.Env, fmt.Sprintf("%v=%v", envVarName, envVarValue))
		}
		sort.Strings(expectedResult.Env)

		/* act */
		actualResult := constructContainerConfig(
			providedCmd,
			providedEnvVars,
			providedImageRef,
			providedPortBindings,
			providedWorkDir,
		)

		/* assert */
		Expect(actualResult).To(Equal(expectedResult))
	})
})
