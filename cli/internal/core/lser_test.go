package core

import (
	"bytes"
	"context"
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clioutput "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	dataresolver "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/model/fakes"
)

var _ = Context("Lser", func() {
	Context("Ls", func() {
		It("should call dataResolver.Resolve w/ expected args", func() {
			/* arrange */
			providedDirRef := "dummyDirRef"

			fakeOpHandle := new(FakeDataHandle)
			fakeOpHandle.ListDescendantsReturns(nil, nil)

			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

			objectUnderTest := newLser(new(clioutput.FakeCliOutput), fakeDataResolver)

			/* act */
			err := objectUnderTest.Ls(
				context.Background(),
				providedDirRef,
			)

			/* assert */
			actualDirRef,
				actualPullCreds := fakeDataResolver.ResolveArgsForCall(0)
			Expect(err).To(BeNil())
			Expect(actualDirRef).To(Equal(providedDirRef))
			Expect(actualPullCreds).To(BeNil())
		})
		It("should return dataResolver.Resolve errors", func() {
			/* arrange */
			expectedError := errors.New("expected")
			fakeDataResolver := new(dataresolver.FakeDataResolver)
			fakeDataResolver.ResolveReturns(nil, expectedError)

			objectUnderTest := newLser(new(clioutput.FakeCliOutput), fakeDataResolver)

			/* act */
			err := objectUnderTest.Ls(context.Background(), "dummy")

			/* assert */
			Expect(err).To(MatchError(err))
		})
		Context("opLister.List", func() {
			It("works", func() {
				/* arrange */
				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.ListDescendantsReturns([]*model.DirEntry{{
					Path: "/path/op.yml",
				}}, nil)
				rs := bytes.NewReader([]byte(`{"name": "validater_valid"}`))
				fakeOpHandle.GetContentReturns(mockReadSeekCloser{rs}, nil)
				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _lsInvoker{
					dataResolver: fakeDataResolver,
					writer:       os.Stdout,
				}

				/* act */
				err := objectUnderTest.Ls(context.Background(), "dummyDirRef")

				/* assert */
				Expect(err).To(BeNil())
			})
			It("should return errors", func() {
				/* arrange */
				expectedError := errors.New("expected")
				fakeOpHandle := new(FakeDataHandle)
				fakeOpHandle.ListDescendantsReturns(nil, expectedError)
				fakeDataResolver := new(dataresolver.FakeDataResolver)
				fakeDataResolver.ResolveReturns(fakeOpHandle, nil)

				objectUnderTest := _lsInvoker{
					dataResolver: fakeDataResolver,
					writer:       os.Stdout,
				}

				/* act */
				err := objectUnderTest.Ls(context.Background(), "dummyDirRef")

				/* assert */
				Expect(err).To(MatchError(expectedError))
			})
		})
	})
})
