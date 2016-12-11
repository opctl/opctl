package bundle

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "errors"
  "reflect"
  "github.com/opspec-io/sdk-golang/util/fs"
  "github.com/opspec-io/sdk-golang/util/format"
)

var _ = Describe("_opViewFactory", func() {

  Context("Construct", func() {

    Context("when FileSystem.GetBytesOfFile returns an error", func() {

      It("should be returned", func() {

        /* arrange */
        expectedError := errors.New("GetBytesOfFileError")

        fakeFileSystem := new(fs.FakeFileSystem)
        fakeFileSystem.GetBytesOfFileReturns(nil, expectedError)

        objectUnderTest := newOpViewFactory(
          fakeFileSystem,
          new(format.FakeFormat),
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

        fakeYamlFormat := new(format.FakeFormat)
        fakeYamlFormat.ToReturns(expectedError)

        objectUnderTest := newOpViewFactory(
          new(fs.FakeFileSystem),
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

      fakeFileSystem := new(fs.FakeFileSystem)
      fakeFileSystem.GetBytesOfFileReturns(expectedBytes, nil)

      fakeYamlFormat := new(format.FakeFormat)

      objectUnderTest := newOpViewFactory(
        fakeFileSystem,
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
      expectedInputs := []model.Param{
        {
          String: &model.StringParam{
            Default:"dummyDefault",
            Description:"dummyDescription",
            MinLength:0,
            MaxLength:1000,
            Name:"dummyName",
            Pattern:".*",
            IsSecret:true,
          },
        },
      }

      expectedCallGraphDeclaration := &model.CallGraphDeclaration{
        Op:&model.OpCallDeclaration{
          Ref:"dummyOpRef",
        },
      }

      expectedOpView := model.OpView{
        Description: "dummyDescription",
        Inputs: expectedInputs,
        Name: "dummyName",
        Run: expectedCallGraphDeclaration,
        Version: "dummyVersion",
      }

      fakeFileSystem := new(fs.FakeFileSystem)

      fakeYamlFormat := new(format.FakeFormat)
      fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

        stubbedOpManifest := model.OpManifest{
          Manifest:model.Manifest{
            Name:expectedOpView.Name,
            Description:expectedOpView.Description,
            Version:expectedOpView.Version,
          },
          Inputs:expectedInputs,
          Run:expectedCallGraphDeclaration,
        }

        reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
        return
      }

      objectUnderTest := newOpViewFactory(
        fakeFileSystem,
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

        expectedCallGraphDeclaration := &model.CallGraphDeclaration{
          Parallel: &model.ParallelCallDeclaration{},
        }

        fakeFileSystem := new(fs.FakeFileSystem)

        fakeYamlFormat := new(format.FakeFormat)
        fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

          stubbedOpManifest := model.OpManifest{
            Run:expectedCallGraphDeclaration,
          }

          reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
          return
        }

        objectUnderTest := newOpViewFactory(
          fakeFileSystem,
          fakeYamlFormat,
        )

        /* act */
        actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

        /* assert */
        Expect(actualOpView.Run).To(Equal(expectedCallGraphDeclaration))

      })
    })
    Context("when opManifest.Run.Parallel is nil", func() {
      It("should return expected opView.Run", func() {

        /* arrange */
        expectedCallGraphDeclaration := &model.CallGraphDeclaration{
          Op:&model.OpCallDeclaration{
            Ref:"dummyOpRef",
          },
        }

        fakeFileSystem := new(fs.FakeFileSystem)

        fakeYamlFormat := new(format.FakeFormat)
        fakeYamlFormat.ToStub = func(in []byte, out interface{}) (err error) {

          stubbedOpManifest := model.OpManifest{
            Run:expectedCallGraphDeclaration,
          }

          reflect.ValueOf(out).Elem().Set(reflect.ValueOf(stubbedOpManifest))
          return
        }

        objectUnderTest := newOpViewFactory(
          fakeFileSystem,
          fakeYamlFormat,
        )

        /* act */
        actualOpView, _ := objectUnderTest.Construct("/dummy/op/path")

        /* assert */
        Expect(actualOpView.Run).To(Equal(expectedCallGraphDeclaration))

      })

    })
  })

})
