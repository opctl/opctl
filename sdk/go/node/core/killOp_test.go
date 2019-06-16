package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/model"
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
			actualCallID,
				actualRootCallID := fakeCallKiller.KillArgsForCall(0)
			Expect(actualCallID).To(Equal(providedReq.OpID))
			Expect(actualRootCallID).To(Equal(providedReq.OpID))
		})
	})
})
