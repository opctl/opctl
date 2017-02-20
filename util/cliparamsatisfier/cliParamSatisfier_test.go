package cliparamsatisfier

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/opctl/util/clicolorer"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/clioutput"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/apiclient"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/pkg/validate"
	"path/filepath"
)

var _ = Context("parameterSatisfier", func() {
	Context("Satisfy", func() {
		Context("op has params", func() {
			Context("args provided explicitly w/ values", func() {
				Context("an arg is invalid", func() {
					It("should return it in the argMap w/ value from env", func() {
						/* arrange */
						param1Name := "DUMMY_PARAM1_NAME"
						param1Value := &model.Data{String: "dummyParam1Value"}

						fakeVos := new(vos.Fake)
						fakeVos.GetenvReturns(param1Value.String)

						objectUnderTest := New(
							new(clicolorer.Fake),
							new(cliexiter.Fake),
							new(clioutput.Fake),
							validate.New(),
							fakeVos,
						)

						expectedResult := map[string]*model.Data{param1Name: param1Value}
						providedArgs := []string{fmt.Sprintf("%v=%v", param1Name, "invalid")}
						providedParams := map[string]*model.Param{
							param1Name: {
								String: &model.StringParam{
									Constraints: &model.StringConstraints{
										Enum: []string{param1Value.String},
									},
								},
							},
						}

						/* act */
						actualResult := objectUnderTest.Satisfy(providedArgs, providedParams)

						/* assert */
						Expect(actualResult).To(Equal(expectedResult))
					})
				})
				Context("all args valid", func() {
					It("should return them in the argMap as provided", func() {
						/* arrange */
						param1Name := "DUMMY_PARAM1_NAME"
						param1Value := &model.Data{String: "dummyParam1Value"}

						objectUnderTest := New(
							new(clicolorer.Fake),
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

			})
			Context("args provided explicitly w/out values", func() {
				It("should return them in the argMap w/ values from env", func() {
					/* arrange */
					param1Name := "DUMMY_PARAM1_NAME"
					param1Value := &model.Data{String: "dummyParam1Value"}

					fakeVos := new(vos.Fake)
					fakeVos.GetenvReturns(param1Value.String)

					objectUnderTest := New(
						new(clicolorer.Fake),
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

						fakeApiClient := new(apiclient.Fake)
						fakeApiClient.StartOpReturns("dummyOpId", errors.New(""))

						objectUnderTest := New(
							new(clicolorer.Fake),
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
						wdReturnedFromGetwd := "dummyWorkDir"

						fakeVos := new(vos.Fake)
						fakeVos.GetwdReturns(wdReturnedFromGetwd, nil)

						objectUnderTest := New(
							new(clicolorer.Fake),
							new(cliexiter.Fake),
							new(clioutput.Fake),
							new(validate.Fake),
							fakeVos,
						)

						providedParam1Name := "dummyParam1Name"
						providedParam2Name := "dummyParam2Name"
						providedParam3Name := "dummyParam3Name"
						providedParams := map[string]*model.Param{
							providedParam1Name: {
								String: &model.StringParam{
									Default: "dummyParam1Default",
								},
							},
							providedParam2Name: {
								File: &model.FileParam{
									Default: "dummyParam2Default",
								},
							},
							providedParam3Name: {
								Dir: &model.DirParam{
									Default: "dummyParam3Default",
								},
							},
						}

						expectedResult := map[string]*model.Data{
							providedParam1Name: {
								String: providedParams[providedParam1Name].String.Default,
							},
							providedParam2Name: {
								File: filepath.Join(wdReturnedFromGetwd, providedParams[providedParam2Name].File.Default),
							},
							providedParam3Name: {
								Dir: filepath.Join(wdReturnedFromGetwd, providedParams[providedParam3Name].Dir.Default),
							},
						}

						/* act */
						actualResult := objectUnderTest.Satisfy([]string{}, providedParams)

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
					new(clicolorer.Fake),
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
