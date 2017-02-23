package pscanary

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Context("psCanary", func() {
	Context("New", func() {
		It("should return Exiter", func() {
			/* arrange/act/assert */
			Expect(New()).Should(Not(BeNil()))
		})
	})
	Context("IsAlive", func() {
		Context("process exists", func() {
			Context("process is alive", func() {
				It("should return true", func() {
					/* arrange */

					objectUnderTest := New()
					expectedResult := true

					// use current PID; we know it's running : )
					providedProcessId := os.Getpid()

					/* act */
					actualResult := objectUnderTest.IsAlive(providedProcessId)

					/* assert */
					Expect(actualResult).
						Should(Equal(expectedResult))
				})
			})
		})
		Context("process doesn't exist", func() {
			It("should return false", func() {
				/* arrange */

				objectUnderTest := New()
				expectedResult := false

				// use ridiculously large PID so we've reasonable confidence it won't exist
				providedProcessId := int(1e9)

				/* act */
				actualResult := objectUnderTest.IsAlive(providedProcessId)

				/* assert */
				Expect(actualResult).
					Should(Equal(expectedResult))
			})
		})
	})
})
