package updater

import (
	"github.com/equinox-io/equinox"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/opctl/opctl/cli/internal/updater/internal/fakes"
	"github.com/pkg/errors"
)

var _ = Context("updater", func() {
	Context("New", func() {
		It("should return Updater", func() {
			/* arrange/act/assert */
			Expect(New()).To(Not(BeNil()))
		})
	})
	Context("_new", func() {
		It("should return Updater", func() {
			/* arrange/act/assert */
			Expect(_new(new(FakeEquinoxClient))).To(Not(BeNil()))
		})
	})
	Context("GetUpdateIfExists", func() {
		Context("update exists", func() {
			It("should return expected update & no error", func() {
				/* arrange */
				equinoxResponse := equinox.Response{
					ReleaseVersion: "dummyReleaseVersion",
				}
				fakeEquinoxClient := new(FakeEquinoxClient)
				fakeEquinoxClient.CheckReturns(equinoxResponse, nil)
				objectUnderTest := _new(fakeEquinoxClient)
				expectedUpdate := Update{
					equinoxResponse: &equinoxResponse,
					Version:         equinoxResponse.ReleaseVersion,
				}

				/* act */
				actualUpdate, actualError := objectUnderTest.GetUpdateIfExists("")

				/* assert */
				// deep equal doesn't work on non-exported fields
				Expect(actualUpdate.equinoxResponse).To(Equal(expectedUpdate.equinoxResponse))
				Expect(actualUpdate.Version).To(Equal(expectedUpdate.Version))
				Expect(actualError).To(BeNil())
			})
		})
		Context("error occurs checking for update", func() {
			It("should return expected error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")
				fakeEquinoxClient := new(FakeEquinoxClient)
				fakeEquinoxClient.CheckReturns(equinox.Response{}, expectedError)
				objectUnderTest := _new(fakeEquinoxClient)

				/* act */
				_, actualError := objectUnderTest.GetUpdateIfExists("")

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("update doesn't exist", func() {
			It("should return nil & no error", func() {
				/* arrange */
				fakeEquinoxClient := new(FakeEquinoxClient)
				fakeEquinoxClient.CheckReturns(equinox.Response{}, equinox.NotAvailableErr)
				objectUnderTest := _new(fakeEquinoxClient)

				/* act */
				actualUpdate, actualError := objectUnderTest.GetUpdateIfExists("")

				/* assert */
				Expect(actualUpdate).To(BeNil())
				Expect(actualError).To(BeNil())
			})
		})
	})
	Context("ApplyUpdate", func() {
		providedUpdate := &Update{
			equinoxResponse: &equinox.Response{},
		}
		Context("Error occurs", func() {
			It("should return the error", func() {
				/* arrange */
				expectedError := errors.New("dummyError")
				fakeEquinoxClient := new(FakeEquinoxClient)
				fakeEquinoxClient.ApplyReturns(expectedError)
				objectUnderTest := _new(fakeEquinoxClient)

				/* act */
				actualError := objectUnderTest.ApplyUpdate(providedUpdate)

				/* assert */
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("No error occurs", func() {
			It("should call update.equinoxResponse.Apply()", func() {
				/* arrange */
				fakeEquinoxClient := new(FakeEquinoxClient)
				objectUnderTest := _new(fakeEquinoxClient)

				/* act */
				objectUnderTest.ApplyUpdate(providedUpdate)

				/* assert */
				Expect(fakeEquinoxClient.ApplyCallCount()).To(Equal(1))
			})
		})
	})
})
