package types

import (
	"fmt"
	"github.com/ghodss/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("SCGOpCall", func() {
	Context("yaml.Unmarshal", func() {
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
