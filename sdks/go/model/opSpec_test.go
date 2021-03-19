package model

import (
	"fmt"

	"github.com/ghodss/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("OpCallSpec", func() {
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

			expectedResult := OpCallSpec{
				Ref: providedPkgRef,
				PullCreds: &CredsSpec{
					Username: providedPullUsername,
					Password: providedPullPassword,
				},
			}

			/* act */
			actualResult := OpCallSpec{}
			err := yaml.Unmarshal([]byte(providedString), &actualResult)
			if err != nil {
				panic(err)
			}

			/* assert */
			Expect(actualResult).To(Equal(expectedResult))

		})
	})
})
