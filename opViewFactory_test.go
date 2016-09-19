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
          new(fakeFormat),
        )

        /* act */
        _, actualError := objectUnderTest.Construct("/dummy/path")

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })

    })

    Context("when YamlFormat.From returns an error", func() {
      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("FromError")

        fakeYamlFormat := new(fakeFormat)
        fakeYamlFormat.ToReturns(expectedError)

        objectUnderTest := newOpViewFactory(
          new(FakeFilesystem),
          fakeYamlFormat,
        )

        /* act */
        _, actualError := objectUnderTest.Construct("/dummy/path")

        /* assert */
        Expect(actualError).To(Equal(expectedError))

      })
    })

    It("should call YamlFormat.From with expected bytes", func() {

      /* arrange */
      expectedBytes := []byte{0, 8, 10}

      fakeFilesystem := new(FakeFilesystem)
      fakeFilesystem.GetBytesOfFileReturns(expectedBytes, nil)

      fakeYamlFormat := new(fakeFormat)

      objectUnderTest := newOpViewFactory(
        fakeFilesystem,
        fakeYamlFormat,
      )

      /* act */
      objectUnderTest.Construct("/dummy/path")

      /* assert */
      actualBytes, _ := fakeYamlFormat.ToArgsForCall(0)
      Expect(actualBytes).To(Equal(expectedBytes))

    })

    It("should return expected opView", func() {

      /* arrange */
      expectedInputs := []models.Param{
        {
          Name:"dummyName",
          Description:"dummyDescription",
          IsSecret:false,
          String: &models.StringParam{
            Default:"dummyDefault",
          },
        },
      }

      expectedRunDeclaration := &models.RunDeclaration{Op:"dummyOpRef"}

      expectedOpView := *models.NewOpView(
        "dummyDescription",
        expectedInputs,
        "dummyName",
        expectedRunDeclaration,
        "dummyVersion",
      )

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlFormat := new(fakeFormat)
      fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

        stubbedOpManifest := models.OpManifest{
          Manifest:models.Manifest{
            Name:expectedOpView.Name,
            Description:expectedOpView.Description,
            Version:expectedOpView.Version,
          },
          Inputs:expectedInputs,
          Run:expectedRunDeclaration,
        }

        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
        return
      }

      objectUnderTest := newOpViewFactory(
        fakeFilesystem,
        fakeYamlFormat,
      )

      /* act */
      actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

      /* assert */
      Expect(actualOpView).To(Equal(expectedOpView))

    })

    Context("when opManifest.Run.Parallel is not nil", func() {
      It("should return expected opView.Run", func() {

        /* arrange */

        expectedRunDeclaration := &models.RunDeclaration{
          Parallel: &models.ParallelRunDeclaration{},
        }

        fakeFilesystem := new(FakeFilesystem)

        fakeYamlFormat := new(fakeFormat)
        fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

          stubbedOpManifest := models.OpManifest{
            Run:expectedRunDeclaration,
          }

          reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
          return
        }

        objectUnderTest := newOpViewFactory(
          fakeFilesystem,
          fakeYamlFormat,
        )

        /* act */
        actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

        /* assert */
        Expect(actualOpView.Run).To(Equal(expectedRunDeclaration))

      })
    })
    Context("when opManifest.Run.Parallel is nil", func() {
      It("should return expected opView.Run", func() {

        /* arrange */
        expectedRunDeclaration := &models.RunDeclaration{Op:"dummyOpRef"}

        fakeFilesystem := new(FakeFilesystem)

        fakeYamlFormat := new(fakeFormat)
        fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

          stubbedOpManifest := models.OpManifest{
            Run:expectedRunDeclaration,
          }

          reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
          return
        }

        objectUnderTest := newOpViewFactory(
          fakeFilesystem,
          fakeYamlFormat,
        )

        /* act */
        actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

        /* assert */
        Expect(actualOpView.Run).To(Equal(expectedRunDeclaration))

      })

    })
  })

})
