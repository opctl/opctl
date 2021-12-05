package docker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("constructRegistryAuth", func() {
	It("should return expected RegistryAuth", func() {
		/* arrange */

		expectedRegistryAuth := "eyJ1c2VybmFtZSI6ImR1bW15UHVsbElkZW50aXR5IiwicGFzc3dvcmQiOiJkdW1teVB1bGxTZWNyZXQifQ=="

		/* act */
		actualRegistryAuth, _ := constructRegistryAuth("dummyPullIdentity", "dummyPullSecret")

		/* assert */
		Expect(actualRegistryAuth).To(Equal(expectedRegistryAuth))
	})
})
