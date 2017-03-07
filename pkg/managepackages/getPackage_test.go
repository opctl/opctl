package managepackages

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

var _ = Describe("_getPackage", func() {

	Context("Execute", func() {

		It("should call packageViewFactory.Construct with expected inputs", func() {

			/* arrange */

			providedPackageRef := "/dummy/path"

			fakePackageViewFactory := new(fakePackageViewFactory)

			objectUnderTest := &managePackages{
				packageViewFactory: fakePackageViewFactory,
			}

			/* act */
			objectUnderTest.GetPackage(
				providedPackageRef,
			)

			/* assert */
			Expect(fakePackageViewFactory.ConstructArgsForCall(0)).To(Equal(providedPackageRef))

		})

		It("should return result of packageViewFactory.Construct", func() {

			/* arrange */
			expectedPackageView := model.PackageView{
				Description: "dummyDescription",
				Inputs:      map[string]*model.Param{},
				Outputs:     map[string]*model.Param{},
				Name:        "dummyName",
				Run: &model.Scg{
					Op: &model.ScgOpCall{
						Ref: "dummyPkgRef",
					},
				},
				Version: "",
			}
			expectedError := errors.New("ConstructError")

			fakePackageViewFactory := new(fakePackageViewFactory)
			fakePackageViewFactory.ConstructReturns(expectedPackageView, expectedError)

			objectUnderTest := &managePackages{
				packageViewFactory: fakePackageViewFactory,
			}

			/* act */
			actualPackageView, actualError := objectUnderTest.GetPackage("/dummy/path")

			/* assert */
			Expect(actualPackageView).To(Equal(expectedPackageView))
			Expect(actualError).To(Equal(expectedError))

		})

	})

})
