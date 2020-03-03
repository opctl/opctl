package cliexiter

import (
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clioutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
)

var _ = Context("cliExiter", func() {
	Context("New", func() {
		It("should return Exiter", func() {
			/* arrange/act/assert */
			Expect(New(
				new(clioutputFakes.FakeCliOutput),
				new(ios.Fake)),
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

				fakeOutput := new(clioutputFakes.FakeCliOutput)
				objectUnderTest := New(
					fakeOutput,
					new(ios.Fake),
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeOutput.ErrorArgsForCall(0)).
					To(Equal(providedExitReq.Message))
			})
			It("should call ios.Exit w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Code: 3,
				}
				fakeIOS := new(ios.Fake)
				objectUnderTest := New(
					new(clioutputFakes.FakeCliOutput),
					fakeIOS,
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeIOS.ExitArgsForCall(0)).To(Equal(providedExitReq.Code))
			})
		})
		Context("req.Code <= 0", func() {
			It("should call output w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Message: "dummyMessage",
				}

				fakeOutput := new(clioutputFakes.FakeCliOutput)
				objectUnderTest := New(
					fakeOutput,
					new(ios.Fake),
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeOutput.SuccessArgsForCall(0)).
					To(Equal(providedExitReq.Message))
			})
			It("should call ios.Exit w/ expected args", func() {
				/* arrange */
				providedExitReq := ExitReq{
					Code: 0,
				}
				fakeIOS := new(ios.Fake)
				objectUnderTest := New(
					new(clioutputFakes.FakeCliOutput),
					fakeIOS,
				)

				/* act */
				objectUnderTest.Exit(providedExitReq)

				/* assert */
				Expect(fakeIOS.ExitArgsForCall(0)).To(Equal(0))
			})
		})
	})
})
