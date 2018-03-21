package op

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/op/dotyml"
)

var _ = Context("Lister", func() {
	Context("NewLister", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewLister()).Should(Not(BeNil()))
		})
	})
	Context("List", func() {
		It("should call dirHandle.ListContents w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()

			providedDirHandle := new(data.FakeHandle)
			// err to trigger immediate return
			providedDirHandle.ListContentsReturns(nil, errors.New("dummyError"))

			objectUnderTest := _lister{}

			/* act */
			objectUnderTest.List(
				providedCtx,
				providedDirHandle,
			)

			/* assert */
			actualCtx := providedDirHandle.ListContentsArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
		})
		Context("dirHandle.ListContents errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedDirHandle := new(data.FakeHandle)
				listContentsErr := errors.New("listContentsErr")
				providedDirHandle.ListContentsReturns(nil, listContentsErr)

				objectUnderTest := _lister{}

				/* act */
				_, actualErr := objectUnderTest.List(
					context.Background(),
					providedDirHandle,
				)

				/* assert */
				Expect(actualErr).To(Equal(listContentsErr))
			})
		})
		Context("dirHandle.ListContents doesn't err", func() {
			Context("dirHandle.ListContents returns items", func() {
				Context("item.Path ends w/ op.yml", func() {
					It("should call dirHandle.GetContent w/ expected args", func() {
						/* arrange */
						providedCtx := context.Background()

						providedDirHandle := new(data.FakeHandle)
						item := model.PkgContent{
							Path: dotyml.FileName,
						}
						providedDirHandle.ListContentsReturns([]*model.PkgContent{&item}, nil)

						// err to trigger immediate return
						providedDirHandle.GetContentReturns(nil, errors.New("dummyError"))

						objectUnderTest := _lister{}

						/* act */
						objectUnderTest.List(
							providedCtx,
							providedDirHandle,
						)

						/* assert */
						actualCtx,
							actualPath := providedDirHandle.GetContentArgsForCall(0)

						Expect(actualCtx).To(Equal(providedCtx))
						Expect(actualPath).To(Equal(item.Path))

					})
					Context("dirHandle.GetContent errs", func() {
						It("should return expected result", func() {
							/* arrange */
							providedDirHandle := new(data.FakeHandle)
							providedDirHandle.ListContentsReturns([]*model.PkgContent{{}}, nil)

							getContentErr := errors.New("getContentErr")
							providedDirHandle.GetContentReturns(nil, getContentErr)

							objectUnderTest := _lister{}

							/* act */
							actualOpYmls, actualErr := objectUnderTest.List(
								context.Background(),
								providedDirHandle,
							)

							/* assert */
							Expect(actualOpYmls).To(BeEmpty())
							Expect(actualErr).To(BeNil())
						})
					})
				})
			})
		})
	})
})
