package cliparamsatisfier

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/clioutput"
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/engineclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
)

var _ = Context("parameterSatisfier", func() {
	Context("Satisfy", func() {
		Context("op has params", func() {
			Context("args provided explicitly w/ values", func() {
				It("should return them in the argMap as provided", func() {
					/* arrange */
					param1Name := "DUMMY_PARAM1_NAME"
					param1Value := &model.Data{String: "dummyParam1Value"}

					objectUnderTest := New(
						new(colorer.Fake),
						new(cliexiter.Fake),
						new(clioutput.Fake),
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
			Context("args provided explicitly w/out values", func() {
				It("should return them in the argMap w/ values from env", func() {
					/* arrange */
					param1Name := "DUMMY_PARAM1_NAME"
					param1Value := &model.Data{String: "dummyParam1Value"}

					fakeVos := new(vos.Fake)
					fakeVos.GetenvReturns(param1Value.String)

					objectUnderTest := New(
						new(colorer.Fake),
						new(cliexiter.Fake),
						new(clioutput.Fake),
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
			Context("args not provided", func() {
				Context("params don't have defaults", func() {
					It("should return them in the argMap w/ values from env", func() {
						/* arrange */
						param1Name := "DUMMY_PARAM1_NAME"
						param1ValueFromEnv := &model.Data{String: "dummyParam1Value"}

						fakeVos := new(vos.Fake)
						fakeVos.GetenvReturns(param1ValueFromEnv.String)

						fakeEngineClient := new(engineclient.Fake)
						fakeEngineClient.StartOpReturns("dummyOpId", errors.New(""))

						objectUnderTest := New(
							new(colorer.Fake),
							new(cliexiter.Fake),
							new(clioutput.Fake),
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
				Context("params have defaults", func() {
					It("should not return them in the argMap", func() {
						/* arrange */
						objectUnderTest := New(
							new(colorer.Fake),
							new(cliexiter.Fake),
							new(clioutput.Fake),
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
		Context("no params", func() {
			It("should return an empty argMap", func() {
				/* arrange */
				expectedResult := map[string]*model.Data{}

				objectUnderTest := New(
					new(colorer.Fake),
					new(cliexiter.Fake),
					new(clioutput.Fake),
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
