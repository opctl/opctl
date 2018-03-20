package docker

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Context("containerConfigFactory", func() {
	Context("Construct", func() {
		It("should return expected result", func() {

			/* arrange */
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

			objectUnderTest := _containerConfigFactory{}

			/* act */
			actualResult := objectUnderTest.Construct(
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
})
