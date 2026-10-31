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

var _ = Context("List", func() {
	It("should call dirHandle.ListDescendants w/ expected args", func() {
		/* arrange */
		providedCtx := context.Background()

		providedDirHandle := new(modelFakes.FakeDataHandle)
		// err to trigger immediate return
		providedDirHandle.ListDescendantsReturns(nil, errors.New("dummyError"))

		/* act */
		List(
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

			/* act */
			_, actualErr := List(
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

					/* act */
					List(
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

						/* act */
						actualOpYmls, actualErr := List(
							context.Background(),
							providedDirHandle,
						)

						/* assert */
						Expect(actualOpYmls).To(BeEmpty())
						Expect(actualErr).
							To(MatchError(fmt.Sprintf("error opening %s%s: %s", expectedRef, expectedPath, getContentErr)))
					})
				})
			})
		})
	})
})

var _ = Describe("List Go ops", func() {
	It("rejects both op declarations in one directory regardless of listing order", func() {
		for _, paths := range [][]string{{"nested/op.go", "nested/op.yml"}, {"nested/op.yml", "nested/op.go"}} {
			dir := new(modelFakes.FakeDataHandle)
			dir.ListDescendantsReturns([]*model.DirEntry{{Path: paths[0]}, {Path: paths[1]}}, nil)
			_, err := List(context.Background(), dir)
			Expect(err).To(MatchError("op dir contains both op.yml and op.go; an op must be a YAML op or a Go op, not both"))
		}
	})
	It("skips Go ops until describe is supported", func() {
		dir := new(modelFakes.FakeDataHandle)
		dir.ListDescendantsReturns([]*model.DirEntry{{Path: "nested/op.go"}}, nil)
		ops, err := List(context.Background(), dir)
		Expect(err).NotTo(HaveOccurred())
		Expect(ops).To(BeEmpty())
	})
})
