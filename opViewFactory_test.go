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

    PContext("when opfile.Run.Container is nil and opfile.Run.SubOps is empty", func() {

      It("should return err", func() {

        /* arrange */
        fakeFilesystem := new(FakeFilesystem)

        fakeYamlCodec := new(fakeYamlCodec)
        fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {

          reflect.ValueOf(out).Elem().Set(reflect.ValueOf(models.OpFile{}))
          return
        }

        objectUnderTest := newOpViewFactory(
          fakeFilesystem,
          fakeYamlCodec,
        )

        /* act */
        _, err := objectUnderTest.Construct("/dummy/op/path")

        /* assert */
        Expect(err).To(Not(BeNil()))

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
      expectedInputs := []models.Parameter{
        *models.NewParameter("dummyName", "dummyDefault", "dummyDescription", false),
      }

      expectedOutputs := []models.Parameter{
        *models.NewParameter("dummyName", "dummyDefault", "dummyDescription", false),
      }

      expectedContainer := &models.Container{}

      expectedOpView := *models.NewOpView(
        "dummyDescription",
        expectedInputs,
        "dummyName",
        expectedOutputs,
        models.NewContainerRunInstruction(expectedContainer),
        "dummyVersion",
      )

      fakeFilesystem := new(FakeFilesystem)

      fakeYamlCodec := new(fakeYamlCodec)
      fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {

        stubbedOpFile := models.OpFile{
          Name:expectedOpView.Name,
          Description:expectedOpView.Description,
          Inputs:expectedInputs,
          Run:models.OpFileRunInstruction{
            Container:expectedContainer,
          },
          Outputs:expectedOutputs,
          Version:expectedOpView.Version,
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

    Context("when opFile.Run.Container is not nil", func() {
      It("should return expected opView.Run", func() {

        /* arrange */

        expectedRunInstruction := models.NewContainerRunInstruction(&models.Container{})

        fakeFilesystem := new(FakeFilesystem)

        fakeYamlCodec := new(fakeYamlCodec)
        fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {

          stubbedOpFile := models.OpFile{
            Run:models.OpFileRunInstruction{
              Container:expectedRunInstruction.Container,
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
        Expect(actualOpView.Run).To(Equal(expectedRunInstruction))

      })
    })
    Context("when opFile.Run.Container is nil", func() {
      It("should return expected opView.Run", func() {

        /* arrange */
        expectedRunInstruction :=
          models.NewSubOpsRunInstruction(
            []models.SubOpRunInstruction{
              {
                Url:"dummyUrl1",
                IsParallel:true,
              },
              {
                Url:"dummyUrl2",
                IsParallel:false,
              },
            })

        fakeFilesystem := new(FakeFilesystem)

        fakeYamlCodec := new(fakeYamlCodec)
        fakeYamlCodec.FromYamlStub = func(in []byte, out interface{}) (err error) {

          stubbedOpFile := models.OpFile{
            Run:models.OpFileRunInstruction{
              SubOps:expectedRunInstruction.SubOps,
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
        Expect(actualOpView.Run).To(Equal(expectedRunInstruction))

      })

    })
  })

})
