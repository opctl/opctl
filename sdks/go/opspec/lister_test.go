package opspec

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	modelFakes "github.com/opctl/opctl/sdks/go/model/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

var _ = Context("Lister", func() {
	Context("NewLister", func() {
		It("should not return nil", func() {
			/* arrange/act/assert */
			Expect(NewLister()).Should(Not(BeNil()))
		})
	})
	Context("List", func() {
		It("should call dirHandle.ListDescendants w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()

			providedDirHandle := new(modelFakes.FakeDataHandle)
			// err to trigger immediate return
			providedDirHandle.ListDescendantsReturns(nil, errors.New("dummyError"))

			objectUnderTest := _lister{}

			/* act */
			objectUnderTest.List(
				providedCtx,
				providedDirHandle,
			)

			/* assert */
			actualCtx := providedDirHandle.ListDescendantsArgsForCall(0)

			Expect(actualCtx).To(Equal(providedCtx))
		})
		Context("dirHandle.ListDescendants errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedDirHandle := new(modelFakes.FakeDataHandle)
				listDescendantsErr := errors.New("listDescendantsErr")
				providedDirHandle.ListDescendantsReturns(nil, listDescendantsErr)

				objectUnderTest := _lister{}

				/* act */
				_, actualErr := objectUnderTest.List(
					context.Background(),
					providedDirHandle,
				)

				/* assert */
				Expect(actualErr).To(Equal(listDescendantsErr))
			})
		})
		Context("dirHandle.ListDescendants doesn't err", func() {
			Context("dirHandle.ListDescendants returns items", func() {
				Context("item.Path ends w/ op.yml", func() {
					It("should call dirHandle.GetContent w/ expected args", func() {
						/* arrange */
						providedCtx := context.Background()

						providedDirHandle := new(modelFakes.FakeDataHandle)
						item := model.DirEntry{
							Path: opfile.FileName,
						}
						providedDirHandle.ListDescendantsReturns([]*model.DirEntry{&item}, nil)

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
							providedDirHandle := new(modelFakes.FakeDataHandle)

							expectedRef := "expectedRef"
							providedDirHandle.RefReturns(expectedRef)

							expectedPath := "op.yml"
							providedDirHandle.ListDescendantsReturns([]*model.DirEntry{{Path: expectedPath}}, nil)

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
							Expect(actualErr).
								To(Equal(fmt.Errorf("error opening %s%s; %s", expectedRef, expectedPath, getContentErr)))
						})
					})
				})
			})
		})
	})
})
