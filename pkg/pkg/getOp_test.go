package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Describe("_getOp", func() {

	Context("Execute", func() {

		It("should call opViewFactory.Construct with expected inputs", func() {

			/* arrange */

			providedOpPackagePath := "/dummy/path"

			fakeOpViewFactory := new(fakeOpViewFactory)

			objectUnderTest := &pkg{
				opViewFactory: fakeOpViewFactory,
			}

			/* act */
			objectUnderTest.GetOp(
				providedOpPackagePath,
			)

			/* assert */
			Expect(fakeOpViewFactory.ConstructArgsForCall(0)).To(Equal(providedOpPackagePath))

		})

		It("should return result of opViewFactory.Construct", func() {

			/* arrange */
			expectedOpView := model.OpView{
				Description: "dummyDescription",
				Inputs:      map[string]*model.Param{},
				Outputs:     map[string]*model.Param{},
				Name:        "dummyName",
				Run: &model.Scg{
					Op: &model.ScgOpCall{
						Ref: "dummyOpPkgRef",
					},
				},
				Version: "",
			}
			expectedError := errors.New("ConstructError")

			fakeOpViewFactory := new(fakeOpViewFactory)
			fakeOpViewFactory.ConstructReturns(expectedOpView, expectedError)

			objectUnderTest := &pkg{
				opViewFactory: fakeOpViewFactory,
			}

			/* act */
			actualOpView, actualError := objectUnderTest.GetOp("/dummy/path")

			/* assert */
			Expect(actualOpView).To(Equal(expectedOpView))
			Expect(actualError).To(Equal(expectedError))

		})

	})

})
