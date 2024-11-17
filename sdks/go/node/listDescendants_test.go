package node

import (
	"context"
	"os"
	"path"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Context("core", func() {
	dbDir, err := os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	db, err := badger.Open(
		badger.DefaultOptions(dbDir).WithLogger(nil),
	)
	if err != nil {
		panic(err)
	}

	Context("ListDescendants", func() {
		Context("req.DataRef empty", func() {
			It("should return expected result", func() {
				/* arrange */
				objectUnderTest := core{}

				/* act */
				actualDescendants, actualErr := objectUnderTest.ListDescendants(
					context.Background(),
					model.ListDescendantsReq{},
				)

				/* assert */
				Expect(actualDescendants).To(BeEmpty())
				Expect(actualErr).To(MatchError(`"" not a valid data ref`))
			})
		})
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		Context("req.DataRef absolute path", func() {
			It("should return expected result", func() {
				/* arrange */
				providedDataRef := path.Join(wd, "testdata/listDescendants")

				objectUnderTest := core{
					stateStore: newStateStore(
						context.Background(),
						db,
						pubsub.New(db),
					),
				}

				expectedDescendants := []*model.DirEntry{
					{Path: "/empty.txt", Size: 0, Mode: 420},
				}

				/* act */
				actualDescendants, actualErr := objectUnderTest.ListDescendants(
					context.Background(),
					model.ListDescendantsReq{
						DataRef: providedDataRef,
					},
				)

				/* assert */
				Expect(actualDescendants).To(ConsistOf(expectedDescendants))
				Expect(actualErr).To(BeNil())
			})
		})
	})
})
