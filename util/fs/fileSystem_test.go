package fs

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "os"
  "path"
  "io/ioutil"
  "github.com/nu7hatch/gouuid"
)

var _ = Describe("_fileSystem", func() {

  Context("newFileSystem()", func() {
    It("should return an instance convertable to FileSystem", func() {

      /* arrange/act */
      objectUnderTest := NewFileSystem()

      /* assert */
      _, ok := objectUnderTest.(FileSystem)
      if !ok {
        Fail("result not assignable to FileSystem")
      }

    })
  })

  Context("AddDir", func() {

    Context("when passed valid path", func() {

      It("should create dir", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        objectUnderTest := NewFileSystem()

        /* act */
        objectUnderTest.AddDir(providedPath)

        /* assert */
        Expect(providedPath).To(BeADirectory())

      })

      It("should return nil err", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        objectUnderTest := NewFileSystem()

        /* act */
        err = objectUnderTest.AddDir(providedPath)

        /* assert */
        Expect(err).To(BeNil())

      })

    })

  })

  Context("GetBytesOfFile", func() {

    Context("when passed path of existing file", func() {

      It("should return expected bytes", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        expectedBytes := []byte("dummyBytes")
        err = ioutil.WriteFile(providedPath, expectedBytes, 0644)
        if (nil != err) {
          panic(err)
        }

        objectUnderTest := NewFileSystem()

        /* act */
        actualBytes, _ := objectUnderTest.GetBytesOfFile(providedPath)

        /* assert */
        Expect(actualBytes).To(Equal(expectedBytes))

      })

      It("should return nil error", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        err = ioutil.WriteFile(providedPath, []byte("dummyBytes"), 0644)
        if (nil != err) {
          panic(err)
        }

        objectUnderTest := NewFileSystem()

        /* act */
        _, err = objectUnderTest.GetBytesOfFile(providedPath)

        /* assert */
        Expect(err).To(BeNil())

      })

    })

  })

  Context("ListChildFileInfosOfDir", func() {

    Context("when passed path of non-existent dir", func() {

      It("should return non nill err", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        objectUnderTest := NewFileSystem()

        /* act */
        _, err = objectUnderTest.ListChildFileInfosOfDir(providedPath)

        /* assert */
        Expect(err).ToNot(BeNil())

      })

    })

    Context("when passed path of existent dir", func() {

      It("should return expected fileInfos", func() {

        /* arrange */
        existentDir, err := os.Getwd()
        if (nil != err) {
          panic(err)
        }

        existentDirWithContent := path.Dir(existentDir)

        expectedChildFileInfos, err := ioutil.ReadDir(existentDirWithContent)

        objectUnderTest := NewFileSystem()

        /* act */
        actualChildFileInfos, _ := objectUnderTest.ListChildFileInfosOfDir(existentDirWithContent)

        /* assert */
        Expect(actualChildFileInfos).To(Equal(expectedChildFileInfos))

      })

      It("should return nill err", func() {

        /* arrange */
        existentDir, err := os.Getwd()
        if (nil != err) {
          panic(err)
        }

        objectUnderTest := NewFileSystem()

        /* act */
        _, err = objectUnderTest.ListChildFileInfosOfDir(existentDir)

        /* assert */
        Expect(err).To(BeNil())

      })

    })

  })

  Context("SaveFile", func() {

    Context("when passed path of non-existing file", func() {

      It("should create file with provided bytes", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        expectedBytes := []byte("dummyBytes")

        objectUnderTest := NewFileSystem()

        /* act */
        objectUnderTest.SaveFile(providedPath, expectedBytes)

        /* assert */
        actualBytes, err := ioutil.ReadFile(providedPath)
        if (nil != err) {
          panic(err)
        }

        Expect(actualBytes).To(Equal(expectedBytes))

      })

      It("should return nil error", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        objectUnderTest := NewFileSystem()

        /* act */
        err = objectUnderTest.SaveFile(providedPath, []byte("dummyBytes"))

        /* assert */
        Expect(err).To(BeNil())

      })

    })

    Context("when passed path of existing file", func() {

      It("should overwrite existing file with provided bytes", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        expectedBytes := []byte("dummyBytes")
        err = ioutil.WriteFile(providedPath, expectedBytes, 0644)
        if (nil != err) {
          panic(err)
        }

        objectUnderTest := NewFileSystem()

        /* act */
        objectUnderTest.SaveFile(providedPath, expectedBytes)

        /* assert */
        actualBytes, err := ioutil.ReadFile(providedPath)
        if (nil != err) {
          panic(err)
        }

        Expect(actualBytes).To(Equal(expectedBytes))

      })

      It("should return nil error", func() {

        /* arrange */
        uuid, err := uuid.NewV4()
        if (nil != err) {
          panic(err)
        }
        providedPath := path.Join(os.TempDir(), uuid.String())

        err = ioutil.WriteFile(providedPath, []byte("dummyBytes"), 0644)
        if (nil != err) {
          panic(err)
        }

        objectUnderTest := NewFileSystem()

        /* act */
        err = objectUnderTest.SaveFile(providedPath, []byte("dummyBytes"))

        /* assert */
        Expect(err).To(BeNil())

      })

    })

  })

})
