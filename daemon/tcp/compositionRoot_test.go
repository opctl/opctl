package tcp

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/engine/daemon/core"
)

var _ = Describe("compositionRoot", func() {

	Context("GetEventBusHandler()", func() {
		It("should not return nil", func() {

			/* arrange */
			objectUnderTest := newCompositionRoot(
				new(core.FakeCore),
			)

			/* act */
			actualGetEventBusHandler := objectUnderTest.GetEventBusHandler()

			/* assert */
			Expect(actualGetEventBusHandler).NotTo(BeNil())

		})
	})

	Context("GetLivenessHandler()", func() {
		It("should not return nil", func() {

			/* arrange */
			objectUnderTest := newCompositionRoot(
				new(core.FakeCore),
			)

			/* act */
			actualGetLivenessHandler := objectUnderTest.GetLivenessHandler()

			/* assert */
			Expect(actualGetLivenessHandler).NotTo(BeNil())

		})
	})

	Context("KillOpHandler()", func() {
		It("should not return nil", func() {

			/* arrange */
			objectUnderTest := newCompositionRoot(
				new(core.FakeCore),
			)

			/* act */
			actualKillOpHandler := objectUnderTest.KillOpHandler()

			/* assert */
			Expect(actualKillOpHandler).NotTo(BeNil())

		})
	})

	Context("StartOpHandler()", func() {
		It("should not return nil", func() {

			/* arrange */
			objectUnderTest := newCompositionRoot(
				new(core.FakeCore),
			)

			/* act */
			actualStartOpHandler := objectUnderTest.StartOpHandler()

			/* assert */
			Expect(actualStartOpHandler).NotTo(BeNil())

		})
	})

})
