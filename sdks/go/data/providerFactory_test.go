package data

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	clientFakes "github.com/opctl/opctl/sdks/go/node/api/client/fakes"
)

var _ = Context("providerFactory", func() {
	Context("NewFSProvider", func() {
		It("should return Provider", func() {
			/* arrange */
			providedBasePaths := []string{"dummyBasePath"}

			objectUnderTest := _providerFactory{}

			/* act */
			actualProvider := objectUnderTest.NewFSProvider(providedBasePaths...)

			/* assert */
			Expect(actualProvider).To(Not(BeNil()))
		})
	})
	Context("NewGitProvider", func() {
		It("should return expected Provider", func() {
			/* arrange */
			providedBasePath := "dummyBasePath"
			providedPullCreds := &model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"}

			objectUnderTest := _providerFactory{}

			/* act */
			actualProvider := objectUnderTest.NewGitProvider(
				providedBasePath,
				providedPullCreds,
			)

			/* assert */
			Expect(actualProvider).To(Not(BeNil()))
		})
	})
	Context("NewNodeProvider", func() {
		It("should return nodeProvider", func() {
			/* arrange */
			fakeAPIClient := new(clientFakes.FakeClient)
			providedPullCreds := &model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"}

			objectUnderTest := _providerFactory{}

			/* act */
			actualProvider := objectUnderTest.NewNodeProvider(
				fakeAPIClient,
				providedPullCreds,
			)

			/* assert */
			Expect(actualProvider).To(Not(BeNil()))
		})
	})
})
