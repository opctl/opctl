package cliexiter

import (
	"github.com/golang-interfaces/vos"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/util/clioutput"
)

var _ = Context("cliExiter", func() {
	Context("New", func() {
		It("should return Exiter", func() {
			/* arrange/act/assert */
			Expect(New(
				new(clioutput.Fake),
				new(vos.Fake)),
			).To(Not(BeNil()))
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

				fakeOutput := new(clioutput.Fake)
				objectUnderTest := New(
					fakeOutput,
					new(vos.Fake),
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeOutput.ErrorArgsForCall(0)).
					To(Equal(providedExitReq.Message))
			})
			It("should call vos.Exit w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Code: 3,
				}
				fakeVOS := new(vos.Fake)
				objectUnderTest := New(
					new(clioutput.Fake),
					fakeVOS,
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeVOS.ExitArgsForCall(0)).To(Equal(providedExitReq.Code))
			})
		})
		Context("req.Code <= 0", func() {
			It("should call output w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Message: "dummyMessage",
				}

				fakeOutput := new(clioutput.Fake)
				objectUnderTest := New(
					fakeOutput,
					new(vos.Fake),
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeOutput.SuccessArgsForCall(0)).
					To(Equal(providedExitReq.Message))
			})
			It("should call vos.Exit w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Code: 0,
				}
				fakeVOS := new(vos.Fake)
				objectUnderTest := New(
					new(clioutput.Fake),
					fakeVOS,
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeVOS.ExitArgsForCall(0)).To(Equal(0))
			})
		})
	})
})
