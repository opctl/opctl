package pkg

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("_get", func() {

	Context("Execute", func() {
		It("should call validate w/ expected inputs", func() {
			/* arrange */
			providedPkgRef := "/dummy/path"

			fakeValidator := new(fakeValidator)

			// err to cause immediate return
			fakeValidator.ValidateReturns([]error{errors.New("dummyError")})

			objectUnderTest := pkg{
				validator: fakeValidator,
			}

			/* act */
			objectUnderTest.Get(providedPkgRef)

			/* assert */
			Expect(fakeValidator.ValidateArgsForCall(0)).To(Equal(providedPkgRef))
		})
		Context("validator.Validate returns errors", func() {
			It("should return the expected error", func() {
				/* arrange */

				errs := []error{errors.New("dummyErr1"), errors.New("dummyErr2")}
				expectedErr := fmt.Errorf(`
-
  Error(s):
    - %v
    - %v
-`, errs[0], errs[1])

				fakeValidator := new(fakeValidator)

				// err to cause immediate return
				fakeValidator.ValidateReturns(errs)

				objectUnderTest := pkg{
					validator: fakeValidator,
				}

				/* act */
				_, actualError := objectUnderTest.Get("")

				/* assert */
				Expect(actualError).To(Equal(expectedErr))
			})
		})
		Context("validator.Validate doesn't return errors", func() {

			It("should call packageViewFactory.Construct with expected inputs", func() {

				/* arrange */

				providedPkgRef := "/dummy/path"

				fakePackageViewFactory := new(fakePackageViewFactory)

				objectUnderTest := &pkg{
					packageViewFactory: fakePackageViewFactory,
					validator:          new(fakeValidator),
				}

				/* act */
				objectUnderTest.Get(
					providedPkgRef,
				)

				/* assert */
				Expect(fakePackageViewFactory.ConstructArgsForCall(0)).To(Equal(providedPkgRef))

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

				objectUnderTest := &pkg{
					packageViewFactory: fakePackageViewFactory,
					validator:          new(fakeValidator),
				}

				/* act */
				actualPackageView, actualError := objectUnderTest.Get("/dummy/path")

				/* assert */
				Expect(actualPackageView).To(Equal(expectedPackageView))
				Expect(actualError).To(Equal(expectedError))

			})
		})

	})

})
