package core

import (
	"context"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/sdks/go/node/core/containerruntime/fakes"
)

var _ = Context("core", func() {
	Context("New", func() {
		It("should return Core", func() {
			/* arrange */
			dataDir, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			/* act/assert */
			Expect(
				New(
					context.Background(),
					new(FakeContainerRuntime),
					dataDir,
				),
			).To(Not(BeNil()))
		})
	})
})
