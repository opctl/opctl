package bundle

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "errors"
)

var _ = Describe("_getOp", func() {

  Context("Execute", func() {

    It("should call opViewFactory.Construct with expected args", func() {

      /* arrange */

      providedOpBundlePath := "/dummy/path"

      fakeOpViewFactory := new(fakeOpViewFactory)

      objectUnderTest := &_bundle{
        opViewFactory:fakeOpViewFactory,
      }

      /* act */
      objectUnderTest.GetOp(
        providedOpBundlePath,
      )

      /* assert */
      Expect(fakeOpViewFactory.ConstructArgsForCall(0)).To(Equal(providedOpBundlePath))

    })

    It("should return result of opViewFactory.Construct", func() {

      /* arrange */
      expectedOpView := model.OpView{
        Description: "dummyDescription",
        Inputs:[]model.Param{},
        Name: "dummyName",
        Run: &model.RunDeclaration{
          Op:&model.OpRunDeclaration{
            Ref:"dummyOpRef",
          },
        },
        Version: "",
      }
      expectedError := errors.New("ConstructError")

      fakeOpViewFactory := new(fakeOpViewFactory)
      fakeOpViewFactory.ConstructReturns(expectedOpView, expectedError)

      objectUnderTest := &_bundle{
        opViewFactory:fakeOpViewFactory,
      }

      /* act */
      actualOpView, actualError := objectUnderTest.GetOp("/dummy/path")

      /* assert */
      Expect(actualOpView).To(Equal(expectedOpView))
      Expect(actualError).To(Equal(expectedError))

    })

  })

})
