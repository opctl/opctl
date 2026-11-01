package goop

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("container stream contract", func() {
	It("accepts the shared SDK stream fixture through the host parser", func() {
		data, err := os.ReadFile("../sdk/testdata/plan.streams.json")
		Expect(err).NotTo(HaveOccurred())
		plan, err := ParsePlan(data)
		Expect(err).NotTo(HaveOccurred())
		producer := (*plan.Run.Serial)[0].Container
		Expect(producer.Stdout).To(Equal("$(opctl_string_1)"))
		Expect(producer.Stderr).To(Equal("$(opctl_string_2)"))
	})
})
