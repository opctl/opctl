package pkg

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/model"
)

var _ = Describe("Get", func() {
	It("should call getter.Get w/ expected inputs", func() {
		/* arrange */
		providedGetReq := &GetReq{PkgRef: "/dummy/pkg/ref"}

		fakeGetter := new(fakeGetter)

		objectUnderTest := &pkg{
			getter:    fakeGetter,
			validator: new(fakeValidator),
		}

		/* act */
		_, err := objectUnderTest.Get(providedGetReq)
		if nil != err {
			panic(err)
		}

		/* assert */
		Expect(fakeGetter.GetArgsForCall(0)).To(Equal(providedGetReq))

	})

	It("should return result of getter.Get", func() {

		/* arrange */
		expectedPkgManifest := &model.PkgManifest{
			Description: "dummyDescription",
			Inputs:      map[string]*model.Param{},
			Outputs:     map[string]*model.Param{},
			Name:        "dummyName",
			Run: &model.SCG{
				Op: &model.SCGOpCall{
					Pkg: &model.SCGOpCallPkg{
						Ref: "dummyPkgRef",
					},
				},
			},
			Version: "",
		}
		expectedError := errors.New("UnmarshalError")

		fakeGetter := new(fakeGetter)
		fakeGetter.GetReturns(expectedPkgManifest, expectedError)

		objectUnderTest := &pkg{
			getter:    fakeGetter,
			validator: new(fakeValidator),
		}

		/* act */
		actualPkgManifest, actualError := objectUnderTest.Get(&GetReq{PkgRef: "/dummy/path"})

		/* assert */
		Expect(actualPkgManifest).To(Equal(expectedPkgManifest))
		Expect(actualError).To(Equal(expectedError))

	})
})
