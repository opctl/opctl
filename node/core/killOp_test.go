package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/model"
)

var _ = Context("core", func() {
	Context("KillOp", func() {
		It("should call callKiller.Kill w/ expected args", func() {
			/* arrange */
			providedReq := model.KillOpReq{OpID: "dummyOpID"}

			fakeCallKiller := new(fakeCallKiller)

			objectUnderTest := _core{
				callKiller: fakeCallKiller,
			}

			/* act */
			objectUnderTest.KillOp(providedReq)

			/* assert */
			Expect(fakeCallKiller.KillArgsForCall(0)).To(Equal(providedReq.OpID))
		})
	})
})
