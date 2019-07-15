package bracketed

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/sdk-golang/data/coerce"
	"github.com/opctl/sdk-golang/model"
)

var _ = Context("coerceToArrayOrObjecter", func() {
	Context("newCoerceToArrayOrObjecter", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(newCoerceToArrayOrObjecter()).Should(Not(BeNil()))
		})
	})
	Context("CoerceToArrayOrObject", func() {
		It("should call coerce.ToArray w/ expected args", func() {
			/* arrange */
			providedData := model.Value{String: new(string)}

			fakeCoerce := new(coerce.Fake)
			fakeCoerce.ToArrayReturns(nil, errors.New("dummyErr"))

			objectUnderTest := _coerceToArrayOrObjecter{
				coerce: fakeCoerce,
			}

			/* act */
			objectUnderTest.CoerceToArrayOrObject(&providedData)

			/* assert */
			actualValue := fakeCoerce.ToArrayArgsForCall(0)

			Expect(*actualValue).To(Equal(providedData))
		})
	})
})
