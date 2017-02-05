package core

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/engineclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
)

var _ = Context("parameterSatisfier", func() {
	Context("Satisfy", func() {
		Context("op has params", func() {
			Context("op args provided explicitly w/ values", func() {
				It("should return them in the argMap as provided", func() {
					/* arrange */
					param1Name := "DUMMY_PARAM1_NAME"
					param1Value := &model.Data{String: "dummyParam1Value"}

					objectUnderTest := newParamSatisfier(
						new(colorer.Fake),
						new(fakeExiter),
						new(fakeOutput),
						new(validate.Fake),
						new(vos.Fake),
					)

					expectedResult := map[string]*model.Data{param1Name: param1Value}
					providedArgs := []string{fmt.Sprintf("%v=%v", param1Name, param1Value.String)}
					providedParams := map[string]*model.Param{
						param1Name: {
							String: &model.StringParam{},
						},
					}

					/* act */
					actualResult := objectUnderTest.Satisfy(providedArgs, providedParams)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
				})
			})
			Context("op args provided explicitly w/out values", func() {
				It("should return them in the argMap w/ values from env", func() {
					/* arrange */
					param1Name := "DUMMY_PARAM1_NAME"
					param1Value := &model.Data{String: "dummyParam1Value"}

					fakeVos := new(vos.Fake)
					fakeVos.GetenvReturns(param1Value.String)

					objectUnderTest := newParamSatisfier(
						new(colorer.Fake),
						new(fakeExiter),
						new(fakeOutput),
						new(validate.Fake),
						fakeVos,
					)

					expectedResult := map[string]*model.Data{param1Name: param1Value}
					providedArgs := []string{param1Name}
					providedParams := map[string]*model.Param{
						param1Name: {
							String: &model.StringParam{},
						},
					}

					/* act */
					actualResult := objectUnderTest.Satisfy(providedArgs, providedParams)

					/* assert */
					Expect(actualResult).To(Equal(expectedResult))
				})
			})
			Context("op args not provided", func() {
				Context("op params don't have defaults", func() {
					It("should return them in the argMap w/ values from env", func() {
						/* arrange */
						param1Name := "DUMMY_PARAM1_NAME"
						param1ValueFromEnv := &model.Data{String: "dummyParam1Value"}

						fakeVos := new(vos.Fake)
						fakeVos.GetenvReturns(param1ValueFromEnv.String)

						fakeEngineClient := new(engineclient.Fake)
						fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

						objectUnderTest := newParamSatisfier(
							new(colorer.Fake),
							new(fakeExiter),
							new(fakeOutput),
							new(validate.Fake),
							fakeVos,
						)

						expectedResult := map[string]*model.Data{param1Name: param1ValueFromEnv}
						providedArgs := []string{}
						providedParams := map[string]*model.Param{
							param1Name: {
								String: &model.StringParam{},
							},
						}

						/* act */
						actualResult := objectUnderTest.Satisfy(providedArgs, providedParams)

						/* assert */
						Expect(actualResult).To(Equal(expectedResult))
					})
				})
				Context("op params have defaults", func() {
					It("should not return them in the argMap", func() {
						/* arrange */
						objectUnderTest := newParamSatisfier(
							new(colorer.Fake),
							new(fakeExiter),
							new(fakeOutput),
							new(validate.Fake),
							new(vos.Fake),
						)

						expectedResult := map[string]*model.Data{}
						providedArgs := []string{}
						providedParams := map[string]*model.Param{
							"dummyParam1Name": {
								String: &model.StringParam{
									Default: "dummyParam1Default",
								},
							},
						}

						/* act */
						actualResult := objectUnderTest.Satisfy(providedArgs, providedParams)

						/* assert */
						Expect(actualResult).To(Equal(expectedResult))
					})
				})
			})
		})
		Context("op doesn't have params", func() {
			It("should return an empty argMap", func() {
				/* arrange */
				expectedResult := map[string]*model.Data{}

				objectUnderTest := newParamSatisfier(
					new(colorer.Fake),
					new(fakeExiter),
					new(fakeOutput),
					new(validate.Fake),
					new(vos.Fake),
				)

				/* act */
				actualResult := objectUnderTest.Satisfy([]string{}, map[string]*model.Param{})

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))
			})
		})
	})
})
