package data

import (
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdk/go/model"
	"net/url"
)

var _ = Context("providerFactory", func() {
	Context("NewFSProvider", func() {
		It("should return expected Provider", func() {
			/* arrange */
			providedBasePaths := []string{"dummyBasePath"}

			objectUnderTest := _providerFactory{}

			/* act */
			actualProvider := objectUnderTest.NewFSProvider(providedBasePaths...)

			/* assert */
			Expect(actualProvider).To(Equal(fsProvider{
				os:        ios.New(),
				basePaths: providedBasePaths,
			}))
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
			Expect(actualProvider).To(Equal(gitProvider{
				localFSProvider: fsProvider{
					os:        ios.New(),
					basePaths: []string{providedBasePath},
				},
				basePath:  providedBasePath,
				puller:    newPuller(),
				pullCreds: providedPullCreds,
			}))
		})
	})
	Context("NewNodeProvider", func() {
		It("should return nodeProvider", func() {
			/* arrange */
			providedAPIBaseURL, err := url.Parse("dummyAPIBaseURL")
			if nil != err {
				panic(err)
			}

			objectUnderTest := _providerFactory{}

			/* act */
			actualProvider := objectUnderTest.NewNodeProvider(
				*providedAPIBaseURL,
				&model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"},
			)

			/* assert */
			// can't check for equality because of nested private props so we check what we can
			if _, ok := actualProvider.(nodeProvider); !ok {
				Fail("actualProvider wrong type")
			}
			Expect(actualProvider).NotTo(BeNil())
		})
	})
})
