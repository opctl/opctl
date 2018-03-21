package dotyml

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"gopkg.in/yaml.v2"
)

var _ = Context("unmarshaller", func() {

	Context("Unmarshal", func() {

		It("should call validate w/ expected inputs", func() {
			/* arrange */
			providedBytes := []byte("dummyBytes")

			fakeValidator := new(fakeValidator)

			// err to cause immediate return
			fakeValidator.ValidateReturns([]error{errors.New("dummyError")})

			objectUnderTest := _unmarshaller{
				validator: fakeValidator,
			}

			/* act */
			objectUnderTest.Unmarshal(providedBytes)

			/* assert */
			Expect(fakeValidator.ValidateArgsForCall(0)).To(Equal(providedBytes))
		})
		Context("Validator.Validate returns errors", func() {
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

				objectUnderTest := _unmarshaller{
					validator: fakeValidator,
				}

				/* act */
				_, actualError := objectUnderTest.Unmarshal(nil)

				/* assert */
				Expect(actualError).To(Equal(expectedErr))
			})
		})
		Context("Validator.Validate doesn't return errors", func() {

			It("should return expected pkgManifest", func() {

				/* arrange */
				paramDefault := "dummyDefault"
				dummyParams := map[string]*model.Param{
					"dummyName": {
						String: &model.StringParam{
							Constraints: &model.StringConstraints{
								MinLength: 0,
								MaxLength: 1000,
								Pattern:   "dummyPattern",
								Format:    "dummyFormat",
								Enum:      []string{"dummyEnumItem1"},
							},
							Default:     &paramDefault,
							Description: "dummyDescription",
							IsSecret:    true,
						},
					},
				}

				expectedPkgManifest := &model.PkgManifest{
					Description: "dummyDescription",
					Inputs:      dummyParams,
					Name:        "dummyName",
					Outputs:     dummyParams,
					Run: &model.SCG{
						Op: &model.SCGOpCall{
							Pkg: &model.SCGOpCallPkg{
								Ref: "dummyOpRef",
							},
						},
					},
					Version: "dummyVersion",
				}
				providedBytes, err := yaml.Marshal(expectedPkgManifest)
				if nil != err {
					panic(err.Error())
				}

				objectUnderTest := _unmarshaller{
					validator: new(fakeValidator),
				}

				/* act */
				actualPkgManifest, _ := objectUnderTest.Unmarshal(providedBytes)

				/* assert */
				Expect(actualPkgManifest).To(Equal(expectedPkgManifest))

			})
		})
	})
})
