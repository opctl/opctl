package apireachabilityensurer

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

var _ = Context("apiReachabilityEnsurer", func() {
	Context("Ensure", func() {
		Context("nodeProvider.ListNodes errs", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeNodeProvider := new(nodeprovider.Fake)
				expectedError := errors.New("dummyError")
				fakeNodeProvider.ListNodesReturns(nil, expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _apiReachabilityEnsurer{
					cliExiter:    fakeCliExiter,
					nodeProvider: fakeNodeProvider,
				}

				/* act */
				objectUnderTest.Ensure()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
