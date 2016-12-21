package core

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

var _ = Describe("createOp", func() {

	fakeWorkDirPathGetter := new(fakeWorkDirPathGetter)
	workDirPath := "/dummyWorkDirPath"
	fakeWorkDirPathGetter.GetReturns(workDirPath)

	Context("Execute", func() {
		It("should invoke bundle.CreateOp with expected args", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)

			providedCollection := "dummyCollection"
			providedName := "dummyName"

			expectedReq := model.CreateOpReq{
				Path:        path.Join(workDirPath, providedCollection, providedName),
				Name:        providedName,
				Description: "dummyOpDescription",
			}

			objectUnderTest := _core{
				bundle:            fakeBundle,
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.CreateOp(providedCollection, expectedReq.Description, expectedReq.Name)

			/* assert */

			Expect(fakeBundle.CreateOpArgsForCall(0)).Should(Equal(expectedReq))

		})
		It("should return error from bundle.CreateOp", func() {
			/* arrange */
			fakeBundle := new(bundle.FakeBundle)
			expectedError := errors.New("dummyError")
			fakeBundle.CreateOpReturns(expectedError)

			fakeExiter := new(fakeExiter)

			objectUnderTest := _core{
				bundle:            fakeBundle,
				exiter:            fakeExiter,
				workDirPathGetter: fakeWorkDirPathGetter,
			}

			/* act */
			objectUnderTest.CreateOp("", "", "")

			/* assert */
			Expect(fakeExiter.ExitArgsForCall(0)).
				Should(Equal(ExitReq{Message: expectedError.Error(), Code: 1}))

		})
	})
})
