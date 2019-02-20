package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
)

var _ = Context("core", func() {
	Context("KillOp", func() {
		It("should call opKiller.Kill w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpID: "dummyOpID"}

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
