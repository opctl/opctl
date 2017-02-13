package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/vos"
)

var _ = Context("_exiter", func() {
	Context("newExiter", func() {
		It("should return Exiter", func() {
			/* arrange/act/assert */
			Expect(newExiter(
				new(fakeOutput),
				new(vos.Fake)),
			).Should(Not(BeNil()))
		})
	})
	Context("Exit", func() {
		Context("req.Code > 0", func() {
			It("should call output w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Code:    3,
					Message: "dummyMessage",
				}

				fakeOutput := new(fakeOutput)
				objectUnderTest := newExiter(
					fakeOutput,
					new(vos.Fake),
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeOutput.ErrorArgsForCall(0)).
					Should(Equal(providedExitReq.Message))
			})
			It("should call vos.Exit w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Code: 3,
				}
				fakeVos := new(vos.Fake)
				objectUnderTest := newExiter(
					new(fakeOutput),
					fakeVos,
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeVos.ExitArgsForCall(0)).Should(Equal(providedExitReq.Code))
			})
		})
		Context("req.Code <= 0", func() {
			It("should call output w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Message: "dummyMessage",
				}

				fakeOutput := new(fakeOutput)
				objectUnderTest := newExiter(
					fakeOutput,
					new(vos.Fake),
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeOutput.SuccessArgsForCall(0)).
					Should(Equal(providedExitReq.Message))
			})
			It("should call vos.Exit w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Code: 0,
				}
				fakeVos := new(vos.Fake)
				objectUnderTest := newExiter(
					new(fakeOutput),
					fakeVos,
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeVos.ExitArgsForCall(0)).Should(Equal(0))
			})
		})
	})
})
