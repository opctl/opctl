package opspec

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/models"
  "errors"
  "reflect"
)

var _ = Describe("_opViewFactory", func() {

  Context("Construct", func() {

    Context("when Filesystem.GetBytesOfFile returns an error", func() {

      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("GetBytesOfFileError")

        fakeFilesystem := new(FakeFilesystem)
        fakeFilesystem.GetBytesOfFileReturns(nil, expectedError)

        objectUnderTest := newOpViewFactory(
          fakeFilesystem,
          new(fakeYamlCodec),
        )

        /* act */
        _, actualError := objectUnderTest.Construct("/dummy/path")

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })

    })

    Context("when YamlCodec.FromYaml returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("FromYamlError")

        fakeYamlCodec := new(fakeYamlCodec)
        fakeYamlCodec.FromYamlReturns(expectedError)

        objectUnderTest := newOpViewFactory(
          new(FakeFilesystem),
          fakeYamlCodec,
        )

        /* act */
        _, actualError := objectUnderTest.Construct("/dummy/path")

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlCodec.FromYaml with expected bytes", func() {

      /* arrange */
      expectedBytes := []byte{0, 8, 10}

      fakeFilesystem := new(FakeFilesystem)
      fakeFilesystem.GetBytesOfFileReturns(expectedBytes, nil)

      fakeYamlCodec := new(fakeYamlCodec)

      objectUnderTest := newOpViewFactory(
        fakeFilesystem,
        fakeYamlCodec,
      )

      /* act */
      objectUnderTest.Construct("/dummy/path")

      /* assert */
      actualBytes, _ := fakeYamlCodec.FromYamlArgsForCall(0)
      Expect(actualBytes).To(Equal(expectedBytes))

    })

    It("should return expected opView", func() {

      /* arrange */
      expectedSubOpUrl := "subOpName"

      expectedOpView := *models.NewOpView(
        "dummyDescription",
        "dummyName",
        []models.OpSummaryView{
          *models.NewOpSummaryView(expectedSubOpUrl),
        },
      )

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {

        stubbedOpFile := models.OpFile{
          Name:expectedOpView.Name,
          Description:expectedOpView.Description,
          SubOps:[]models.OpFileSubOp{
            models.OpFileSubOp{
              Url:expectedSubOpUrl,
            },
          },
        }

        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpFile))
        return
      }

      objectUnderTest := newOpViewFactory(
        fakeFilesystem,
        fakeYamlCodec,
      )

      /* act */
      actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

      /* assert */
      Expect(actualOpView).To(Equal(expectedOpView))

    })

  })

})
