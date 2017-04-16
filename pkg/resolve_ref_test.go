package pkg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResolveRef", func() {
	It("should call refResolver.Resolve w/ expected inputs", func() {
		/* arrange */
		providedPkgRef := "/dummy/pkg/ref"

		fakeRefResolver := new(fakeRefResolver)

		objectUnderTest := &pkg{
			refResolver: fakeRefResolver,
		}

		/* act */
		objectUnderTest.ResolveRef(providedPkgRef)

		/* assert */
		Expect(fakeRefResolver.ResolveArgsForCall(0)).To(Equal(providedPkgRef))

	})

	It("should return result of refResolver.Resolve", func() {

		/* arrange */
		expectedResult := "dummyPkgRef"

		fakeRefResolver := new(fakeRefResolver)
		fakeRefResolver.ResolveReturns(expectedResult)

		objectUnderTest := &pkg{
			refResolver: fakeRefResolver,
		}

		/* act */
		actualResult := objectUnderTest.ResolveRef("")

		/* assert */
		Expect(actualResult).To(Equal(expectedResult))

	})
})
