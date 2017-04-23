package pkg

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resolve", func() {
	It("should call resolver.Resolve w/ expected inputs", func() {
		/* arrange */
		providedBasePath := "dummyBasePath"
		providedPkgRef := "dummyPkgRef"

		expectedBasePath := providedBasePath
		expectedPkgRef := providedPkgRef

		fakeResolver := new(fakeResolver)

		objectUnderTest := &pkg{
			resolver:  fakeResolver,
			validator: new(fakeValidator),
		}

		/* act */
		objectUnderTest.Resolve(providedBasePath, providedPkgRef)

		/* assert */
		actualBasePath, actualPkgRef := fakeResolver.ResolveArgsForCall(0)
		Expect(actualBasePath).To(Equal(expectedBasePath))
		Expect(actualPkgRef).To(Equal(expectedPkgRef))

	})

	It("should return result of resolver.Resolve", func() {

		/* arrange */
		expectedAbsPath := "dummyAbsPath"
		expectedOk := true

		fakeResolver := new(fakeResolver)
		fakeResolver.ResolveReturns(expectedAbsPath, expectedOk)

		objectUnderTest := &pkg{
			resolver: fakeResolver,
		}

		/* act */
		actualAbsPath, actualOk := objectUnderTest.Resolve("", "")

		/* assert */
		Expect(actualAbsPath).To(Equal(expectedAbsPath))
		Expect(actualOk).To(Equal(expectedOk))

	})
})
