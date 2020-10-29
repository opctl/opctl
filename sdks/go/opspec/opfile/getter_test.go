package opfile

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("opfile", func() {

	Context("Get", func() {
		It("should call ioUtil.ReadFile w/ expected args", func() {
			/* arrange */
			providedOpPath := "providedOpPath"

			fakeIOUtil := new(iioutil.Fake)
			// err to trigger immediate return
			fakeIOUtil.ReadFileReturns(nil, errors.New("dummyError"))

			objectUnderTest := _getter{
				ioUtil: fakeIOUtil,
			}

			/* act */
			objectUnderTest.Get(
				context.Background(),
				providedOpPath,
			)

			/* assert */
			actualPath := fakeIOUtil.ReadFileArgsForCall(0)

			Expect(actualPath).To(Equal(filepath.Join(providedOpPath, FileName)))
		})
		Context("ioUtil.ReadFile errs", func() {
			It("should return error", func() {
				/* arrange */
				readAllErr := errors.New("dummyError")
				fakeIOUtil := new(iioutil.Fake)
				fakeIOUtil.ReadFileReturns(nil, readAllErr)

				objectUnderTest := _getter{
					ioUtil: fakeIOUtil,
				}

				/* act */
				_, actualErr := objectUnderTest.Get(
					context.Background(),
					"dummyOpPath",
				)

				/* assert */
				Expect(actualErr).To(Equal(readAllErr))
			})
		})
		// @TODO: re-enable once unmarshaller is internal
		// Context("ioUtil.ReadAll doesn't err", func() {
		// 	It("should call opFile.Unmarshal w/ expected args & return result", func() {
		// 		/* arrange */
		// 		file, err := ioutil.TempFile("", "")
		// 		if nil != err {
		// 			panic(err)
		// 		}

		// 		providedOpHandle := new(modelFakes.FakeDataHandle)
		// 		providedOpHandle.GetContentReturns(file, nil)

		// 		bytesFromReadAll := []byte{2, 3}
		// 		fakeIOUtil := new(iioutil.Fake)
		// 		fakeIOUtil.ReadAllReturns(bytesFromReadAll, nil)

		// 		expectedOpFile := &model.OpSpec{
		// 			Name: "dummyName",
		// 		}
		// 		expectedErr := errors.New("dummyError")
		// 		FakeUnmarshaller := new(FakeUnmarshaller)

		// 		FakeUnmarshaller.UnmarshalReturns(expectedOpFile, expectedErr)

		// 		objectUnderTest := _getter{
		// 			ioUtil:       fakeIOUtil,
		// 			unmarshaller: FakeUnmarshaller,
		// 		}

		// 		/* act */
		// 		actualOpFile, actualErr := objectUnderTest.Get(
		// 			context.Background(),
		// 			providedOpHandle,
		// 		)

		// 		/* assert */
		// 		Expect(FakeUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(bytesFromReadAll))
		// 		Expect(actualOpFile).To(Equal(expectedOpFile))
		// 		Expect(actualErr).To(Equal(expectedErr))
		// 	})
		// })
	})
})
