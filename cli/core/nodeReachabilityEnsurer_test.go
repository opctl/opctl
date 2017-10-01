package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/util/cliexiter"
)

var _ = Context("nodeReachabilityEnsurer", func() {
	Context("EnsureNodeReachable", func() {
		Context("nodeProvider.ListNodes errs", func() {
			It("should call exiter w/ expected args", func() {
				/* arrange */
				fakeNodeProvider := new(nodeprovider.Fake)
				expectedError := errors.New("dummyError")
				fakeNodeProvider.ListNodesReturns(nil, expectedError)

				fakeCliExiter := new(cliexiter.Fake)

				objectUnderTest := _nodeReachabilityEnsurer{
					cliExiter:    fakeCliExiter,
					nodeProvider: fakeNodeProvider,
				}

				/* act */
				objectUnderTest.EnsureNodeReachable()

				/* assert */
				Expect(fakeCliExiter.ExitArgsForCall(0)).
					To(Equal(cliexiter.ExitReq{Message: expectedError.Error(), Code: 1}))
			})
		})
	})
})
