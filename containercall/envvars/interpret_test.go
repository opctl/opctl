package envvars

import (
	"encoding/json"
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

var _ = Context("EnvVars", func() {
	Context("Interpret", func() {
		Context("implicitly bound", func() {
			Context("name not in scope", func() {
				It("should return expected error", func() {
					/* arrange */
					envVarName := "dummyEnvVarName"
					providedSCGContainerCallEnvVars := map[string]string{
						// implicitly bound
						envVarName: "",
					}

					expectedErr := fmt.Errorf("Unable to bind env var to '%v' via implicit ref; '%v' not in scope", envVarName, envVarName)

					objectUnderTest := _EnvVars{}

					/* act */
					_, actualErr := objectUnderTest.Interpret(
						map[string]*model.Value{},
						providedSCGContainerCallEnvVars,
					)

					/* assert */
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
		})
		It("should return expected dcg.EnvVars", func() {
			/* arrange */
			providedCurrentScopeRef1 := "dummyScopeRef1"
			providedCurrentScopeRef1String := "dummyScopeRef1String"

			providedCurrentScopeRef2 := "dummyScopeRef2"
			providedCurrentScopeRef2Number := float64(2.3)

			providedCurrentScopeRef3 := "dummyScopeRef3"
			providedScopeRef3Object := map[string]interface{}{"dummyProp1Name": "dummyProp1Value"}

			providedCurrentScope := map[string]*model.Value{
				providedCurrentScopeRef1: {String: &providedCurrentScopeRef1String},
				providedCurrentScopeRef2: {Number: &providedCurrentScopeRef2Number},
				providedCurrentScopeRef3: {Object: providedScopeRef3Object},
			}

			providedScopeRef3ObjectJsonBytes, err := json.Marshal(providedScopeRef3Object)
			if nil != err {
				Fail(err.Error())
			}

			fakeJson := new(ijson.Fake)
			fakeJson.MarshalReturns(providedScopeRef3ObjectJsonBytes, err)

			expectedEnvVars := map[string]string{
				providedCurrentScopeRef1: providedCurrentScopeRef1String,
				providedCurrentScopeRef2: strconv.FormatFloat(providedCurrentScopeRef2Number, 'f', -1, 64),
				providedCurrentScopeRef3: string(providedScopeRef3ObjectJsonBytes),
			}

			providedSCGContainerCallEnvVars := map[string]string{
				// implicitly bound to string
				providedCurrentScopeRef1: "",
				// implicitly bound to number
				providedCurrentScopeRef2: "",
				// implicitly bound to object
				providedCurrentScopeRef3: "",
			}

			objectUnderTest := _EnvVars{
				json: fakeJson,
			}

			/* act */
			actualDCGContainerCallEnvVars, _ := objectUnderTest.Interpret(
				providedCurrentScope,
				providedSCGContainerCallEnvVars,
			)

			/* assert */
			Expect(actualDCGContainerCallEnvVars).To(Equal(expectedEnvVars))
		})
	})
})
