package opfile

import (
	"context"
	"errors"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	"io/ioutil"
)

var _ = Context("pkg", func() {

	Context("GetManifest", func() {

		It("should call opHandle.GetContent w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedOpHandle := new(modelFakes.FakeDataHandle)
			// err to trigger immediate return
			providedOpHandle.GetContentReturns(nil, errors.New("dummyError"))

			objectUnderTest := _getter{}

			/* act */
			objectUnderTest.Get(
				providedCtx,
				providedOpHandle,
			)

			/* assert */
			actualCtx,
				actualPath := providedOpHandle.GetContentArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualPath).To(Equal(FileName))
		})
		Context("opHandle.GetContent errs", func() {
			It("should return error", func() {
				/* arrange */
				getContentErr := errors.New("dummyError")

				providedOpHandle := new(modelFakes.FakeDataHandle)
				// err to trigger immediate return
				providedOpHandle.GetContentReturns(nil, getContentErr)

				objectUnderTest := _getter{}

				/* act */
				_, actualErr := objectUnderTest.Get(
					context.Background(),
					providedOpHandle,
				)

				/* assert */
				Expect(actualErr).To(Equal(getContentErr))
			})
		})
		Context("opHandle.GetContent doesn't err", func() {
			It("should call ioUtil.ReadAll w/ expected args", func() {
				/* arrange */
				file, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}

				providedOpHandle := new(modelFakes.FakeDataHandle)
				providedOpHandle.GetContentReturns(file, nil)

				fakeIOUtil := new(iioutil.Fake)
				// err to trigger immediate return
				fakeIOUtil.ReadAllReturns(nil, errors.New("dummyError"))

				objectUnderTest := _getter{
					ioUtil: fakeIOUtil,
				}

				/* act */
				objectUnderTest.Get(
					context.Background(),
					providedOpHandle,
				)

				/* assert */
				actualReader := fakeIOUtil.ReadAllArgsForCall(0)

				Expect(actualReader).To(Equal(file))
			})
			Context("ioUtil.ReadAll errs", func() {
				It("should return error", func() {
					/* arrange */
					file, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}

					providedOpHandle := new(modelFakes.FakeDataHandle)
					providedOpHandle.GetContentReturns(file, nil)

					readAllErr := errors.New("dummyError")
					fakeIOUtil := new(iioutil.Fake)
					fakeIOUtil.ReadAllReturns(nil, readAllErr)

					objectUnderTest := _getter{
						ioUtil: fakeIOUtil,
					}

					/* act */
					_, actualErr := objectUnderTest.Get(
						context.Background(),
						providedOpHandle,
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

			// 		expectedOpFile := &model.OpFile{
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
})
