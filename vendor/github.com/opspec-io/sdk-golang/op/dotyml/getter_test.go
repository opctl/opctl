package dotyml

import (
	"context"
	"errors"
	"github.com/golang-interfaces/iioutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
)

var _ = Context("pkg", func() {

	Context("GetManifest", func() {

		It("should call opDirHandle.GetContent w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedOpDirHandle := new(data.FakeHandle)
			// err to trigger immediate return
			providedOpDirHandle.GetContentReturns(nil, errors.New("dummyError"))

			objectUnderTest := _getter{}

			/* act */
			objectUnderTest.Get(
				providedCtx,
				providedOpDirHandle,
			)

			/* assert */
			actualCtx,
				actualPath := providedOpDirHandle.GetContentArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
			Expect(actualPath).To(Equal(FileName))
		})
		Context("opDirHandle.GetContent errs", func() {
			It("should return error", func() {
				/* arrange */
				getContentErr := errors.New("dummyError")

				providedOpDirHandle := new(data.FakeHandle)
				// err to trigger immediate return
				providedOpDirHandle.GetContentReturns(nil, getContentErr)

				objectUnderTest := _getter{}

				/* act */
				_, actualErr := objectUnderTest.Get(
					context.Background(),
					providedOpDirHandle,
				)

				/* assert */
				Expect(actualErr).To(Equal(getContentErr))
			})
		})
		Context("opDirHandle.GetContent doesn't err", func() {
			It("should call ioUtil.ReadAll w/ expected args", func() {
				/* arrange */
				file, err := ioutil.TempFile("", "")
				if nil != err {
					panic(err)
				}

				providedOpDirHandle := new(data.FakeHandle)
				providedOpDirHandle.GetContentReturns(file, nil)

				fakeIOUtil := new(iioutil.Fake)
				// err to trigger immediate return
				fakeIOUtil.ReadAllReturns(nil, errors.New("dummyError"))

				objectUnderTest := _getter{
					ioUtil: fakeIOUtil,
				}

				/* act */
				objectUnderTest.Get(
					context.Background(),
					providedOpDirHandle,
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

					providedOpDirHandle := new(data.FakeHandle)
					providedOpDirHandle.GetContentReturns(file, nil)

					readAllErr := errors.New("dummyError")
					fakeIOUtil := new(iioutil.Fake)
					fakeIOUtil.ReadAllReturns(nil, readAllErr)

					objectUnderTest := _getter{
						ioUtil: fakeIOUtil,
					}

					/* act */
					_, actualErr := objectUnderTest.Get(
						context.Background(),
						providedOpDirHandle,
					)

					/* assert */
					Expect(actualErr).To(Equal(readAllErr))
				})
			})
			Context("ioUtil.ReadAll doesn't err", func() {
				It("should call manifest.Unmarshal w/ expected args & return result", func() {
					/* arrange */
					file, err := ioutil.TempFile("", "")
					if nil != err {
						panic(err)
					}

					providedOpDirHandle := new(data.FakeHandle)
					providedOpDirHandle.GetContentReturns(file, nil)

					bytesFromReadAll := []byte{2, 3}
					fakeIOUtil := new(iioutil.Fake)
					fakeIOUtil.ReadAllReturns(bytesFromReadAll, nil)

					expectedPkgManifest := &model.PkgManifest{
						Name: "dummyName",
					}
					expectedErr := errors.New("dummyError")
					FakeUnmarshaller := new(FakeUnmarshaller)

					FakeUnmarshaller.UnmarshalReturns(expectedPkgManifest, expectedErr)

					objectUnderTest := _getter{
						ioUtil:       fakeIOUtil,
						unmarshaller: FakeUnmarshaller,
					}

					/* act */
					actualPkgManifest, actualErr := objectUnderTest.Get(
						context.Background(),
						providedOpDirHandle,
					)

					/* assert */
					Expect(FakeUnmarshaller.UnmarshalArgsForCall(0)).To(Equal(bytesFromReadAll))
					Expect(actualPkgManifest).To(Equal(expectedPkgManifest))
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
		})
	})
})
