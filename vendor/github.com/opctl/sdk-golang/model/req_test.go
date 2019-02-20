package model

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("StartOpReq", func() {
	Context("json.Unmarshal", func() {
		Context("deprecated pkg property present", func() {
			It("should return expected result", func() {
				/* arrange */
				providedPkgRef := "pkgRef"
				providedPullUsername := "pullUsername"
				providedPullPassword := "pullPassword"

				providedString := fmt.Sprintf(
					`{"op":{"ref": "%v", "pullCreds": {"username": "%v", "password": "%v"}}}`,
					providedPkgRef,
					providedPullUsername,
					providedPullPassword,
				)

				expectedResult := StartOpReq{
					Op: StartOpReqOp{
						Ref: providedPkgRef,
						PullCreds: &PullCreds{
							Username: providedPullUsername,
							Password: providedPullPassword,
						},
					},
				}

				/* act */
				actualResult := StartOpReq{}
				err := json.Unmarshal([]byte(providedString), &actualResult)
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))

			})
		})
		Context("deprecated pkg property not present", func() {
			It("should return expected result", func() {
				/* arrange */
				providedPkgRef := "pkgRef"
				providedPullUsername := "pullUsername"
				providedPullPassword := "pullPassword"

				providedString := fmt.Sprintf(
					`{"pkg":{"ref": "%v", "pullCreds": {"username": "%v", "password": "%v"}}}`,
					providedPkgRef,
					providedPullUsername,
					providedPullPassword,
				)

				expectedResult := StartOpReq{
					Op: StartOpReqOp{
						Ref: providedPkgRef,
						PullCreds: &PullCreds{
							Username: providedPullUsername,
							Password: providedPullPassword,
						},
					},
				}

				/* act */
				actualResult := StartOpReq{}
				err := json.Unmarshal([]byte(providedString), &actualResult)
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))

			})
		})
	})
})
