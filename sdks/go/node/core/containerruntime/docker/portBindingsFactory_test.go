package docker

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("portBindingsFactory", func() {
	Context("Construct", func() {
		It("should return expected result", func() {
			/* arrange */
			providedPorts := map[string]string{
				"6060/udp":  "6060",
				"8080-8081": "9090-9091",
			}
			expectedPortBindings := nat.PortMap{}
			for containerPort, hostPort := range providedPorts {
				portMappings, _ := nat.ParsePortSpec(fmt.Sprintf("%v:%v", hostPort, containerPort))
				for _, portMapping := range portMappings {
					if _, ok := expectedPortBindings[portMapping.Port]; ok {
						expectedPortBindings[portMapping.Port] = append(expectedPortBindings[portMapping.Port], portMapping.Binding)
					} else {
						expectedPortBindings[portMapping.Port] = []nat.PortBinding{portMapping.Binding}
					}
				}
			}

			objectUnderTest := _portBindingsFactory{}

			/* act */
			actualPorts, actualErr := objectUnderTest.Construct(
				providedPorts,
			)

			/* assert */
			Expect(actualPorts).To(Equal(expectedPortBindings))
			Expect(actualErr).To(BeNil())
		})
	})
})
