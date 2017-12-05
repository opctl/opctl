package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("core", func() {
	Context("KillOp", func() {
		It("should call opKiller.Kill w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpId: "dummyOpId"}

			fakeOpKiller := new(fakeOpKiller)

			objectUnderTest := _core{
				opKiller: fakeOpKiller,
			}

			/* act */
			objectUnderTest.KillOp(providedReq)

			/* assert */
			Expect(fakeOpKiller.KillArgsForCall(0)).To(Equal(providedReq))
		})
	})
})
