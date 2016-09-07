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
      expectedInputs := []models.Param{
        {
          String: &models.StringParam{
            Name:"dummyName",
            Default:"dummyDefault",
            Description:"dummyDescription",
            IsSecret:false,
          },
        },
      }

      expectedRunStatement := &models.RunStatement{Op:"dummyOpRef"}

      expectedOpView := *models.NewOpView(
        "dummyDescription",
        expectedInputs,
        "dummyName",
        expectedRunStatement,
        "dummyVersion",
      )

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {

        stubbedOpBundleManifest := models.OpBundleManifest{
          BundleManifest:models.BundleManifest{
            Name:expectedOpView.Name,
            Description:expectedOpView.Description,
            Version:expectedOpView.Version,
          },
          Inputs:expectedInputs,
          Run:expectedRunStatement,
        }

        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpBundleManifest))
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

    Context("when opBundleManifest.Run.Parallel is not nil", func() {
      It("should return expected opView.Run", func() {

        /* arrange */

        expectedRunStatement := &models.RunStatement{
          Parallel: &models.ParallelRunStatement{},
        }

        fakeFilesystem := new(FakeFilesystem)

        fakeYamlCodec := new(fakeYamlCodec)
        fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {

          stubbedOpBundleManifest := models.OpBundleManifest{
            Run:expectedRunStatement,
          }

          reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpBundleManifest))
          return
        }

        objectUnderTest := newOpViewFactory(
          fakeFilesystem,
          fakeYamlCodec,
        )

        /* act */
        actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

        /* assert */
        Expect(actualOpView.Run).To(Equal(expectedRunStatement))

      })
    })
    Context("when opBundleManifest.Run.Parallel is nil", func() {
      It("should return expected opView.Run", func() {

        /* arrange */
        expectedRunStatement := &models.RunStatement{Op:"dummyOpRef"}

        fakeFilesystem := new(FakeFilesystem)

        fakeYamlCodec := new(fakeYamlCodec)
        fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {

          stubbedOpBundleManifest := models.OpBundleManifest{
            Run:expectedRunStatement,
          }

          reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpBundleManifest))
          return
        }

        objectUnderTest := newOpViewFactory(
          fakeFilesystem,
          fakeYamlCodec,
        )

        /* act */
        actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

        /* assert */
        Expect(actualOpView.Run).To(Equal(expectedRunStatement))

      })

    })
  })

})
