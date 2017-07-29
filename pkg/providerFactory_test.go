package pkg

import (
	"github.com/golang-interfaces/ios"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Context("providerFactory", func() {
	Context("NewGitProvider", func() {
		It("should return expected Provider", func() {
			/* arrange */
			providedBasePath := "dummyBasePath"
			providedPullCreds := &model.PullCreds{Username: "dummyUsername", Password: "dummyPassword"}

			objectUnderTest := _ProviderFactory{}

			/* act */
			actualProvider := objectUnderTest.NewGitProvider(
				providedBasePath,
				providedPullCreds,
			)

			/* assert */
			Expect(actualProvider).To(Equal(gitProvider{
				localFSProvider: localFSProvider{
					os:        ios.New(),
					basePaths: []string{providedBasePath},
				},
				basePath:  providedBasePath,
				puller:    newPuller(),
				pullCreds: providedPullCreds,
			}))
		})
	})
	Context("NewLocalFSProvider", func() {
		It("should return expected Provider", func() {
			/* arrange */
			providedBasePaths := []string{"dummyBasePath"}

			objectUnderTest := _ProviderFactory{}

			/* act */
			actualProvider := objectUnderTest.NewLocalFSProvider(providedBasePaths...)

			/* assert */
			Expect(actualProvider).To(Equal(localFSProvider{
				os:        ios.New(),
				basePaths: providedBasePaths,
			}))
		})
	})
})
