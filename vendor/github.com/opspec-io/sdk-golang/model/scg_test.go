package model

import (
	"fmt"
	"github.com/ghodss/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("SCGOpCall", func() {
	Context("yaml.Unmarshal", func() {
		Context("deprecated pkg property present", func() {
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

				expectedResult := SCGOpCall{
					Ref: providedPkgRef,
					PullCreds: &SCGPullCreds{
						Username: providedPullUsername,
						Password: providedPullPassword,
					},
				}

				/* act */
				actualResult := SCGOpCall{}
				err := yaml.Unmarshal([]byte(providedString), &actualResult)
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
					`{"ref": "%v", "pullCreds": {"username": "%v", "password": "%v"}}`,
					providedPkgRef,
					providedPullUsername,
					providedPullPassword,
				)

				expectedResult := SCGOpCall{
					Ref: providedPkgRef,
					PullCreds: &SCGPullCreds{
						Username: providedPullUsername,
						Password: providedPullPassword,
					},
				}

				/* act */
				actualResult := SCGOpCall{}
				err := yaml.Unmarshal([]byte(providedString), &actualResult)
				if nil != err {
					panic(err)
				}

				/* assert */
				Expect(actualResult).To(Equal(expectedResult))

			})
		})
	})
})
