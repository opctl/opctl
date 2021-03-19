package call

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op"
)

var _ = Context("Interpret", func() {
	Context("callSpec.If not nil", func() {
		Context("predicates returns err", func() {
			It("should return expected result", func() {
				/* arrange */
				predicateSpec := []*model.PredicateSpec{{}}
				dataDir, err := ioutil.TempDir("", "")
				if err != nil {
					panic(err)
				}

				/* act */
				_, actualError := Interpret(
					context.Background(),
					map[string]*model.Value{},
					&model.CallSpec{
						If: &predicateSpec,
					},
					"providedID",
					"dummyOpPath",
					nil,
					"providedRootCallID",
					dataDir,
				)

				/* assert */
				Expect(actualError).To(MatchError("unable to interpret predicate: predicate was unexpected type &{Eq:<nil> Exists:<nil> Ne:<nil> NotExists:<nil>}"))
			})
		})
	})
	Context("callSpec.Container not nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScope := map[string]*model.Value{}
			providedID := "providedID"
			providedOpPath := "providedOpPath"
			providedParentIDValue := "providedParentID"
			providedParentID := &providedParentIDValue
			providedRootCallID := "providedRootCallID"
			providedDataDirPath, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			containerSpec := model.ContainerCallSpec{
				Image: &model.ContainerCallImageSpec{
					Ref: "ref",
				},
			}

			expectedContainer, err := container.Interpret(
				providedScope,
				&containerSpec,
				providedID,
				providedOpPath,
				providedDataDirPath,
			)
			if err != nil {
				panic(err)
			}

			expectedCall := &model.Call{
				Container: expectedContainer,
				ID:        providedID,
				ParentID:  providedParentID,
				RootID:    providedRootCallID,
			}

			/* act */
			actualCall, actualError := Interpret(
				context.Background(),
				providedScope,
				&model.CallSpec{
					Container: &containerSpec,
				},
				providedID,
				providedOpPath,
				providedParentID,
				providedRootCallID,
				providedDataDirPath,
			)

			/* assert */
			Expect(actualError).To(BeNil())
			Expect(actualCall).To(Equal(expectedCall))
		})
	})
	Context("callSpec.Op not nil", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScope := map[string]*model.Value{}
			providedID := "providedID"
			providedOpPath := "providedOpPath"
			providedParentIDValue := "providedParentID"
			providedParentID := &providedParentIDValue
			providedRootCallID := "providedRootCallID"
			providedDataDirPath, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			opRef := filepath.Join(wd, "testdata/testop")

			opSpec := model.OpCallSpec{
				Ref: opRef,
			}

			expectedOp, err := op.Interpret(
				context.Background(),
				providedScope,
				&opSpec,
				providedID,
				providedOpPath,
				providedDataDirPath,
			)
			if err != nil {
				panic(err)
			}

			expectedCall := &model.Call{
				Op:       expectedOp,
				ID:       providedID,
				ParentID: providedParentID,
				RootID:   providedRootCallID,
			}

			/* act */
			actualCall, actualError := Interpret(
				context.Background(),
				providedScope,
				&model.CallSpec{
					Op: &opSpec,
				},
				providedID,
				providedOpPath,
				providedParentID,
				providedRootCallID,
				providedDataDirPath,
			)

			/* assert */
			Expect(actualError).To(BeNil())
			// ignore Op.ChildCallID since it's a generated UUID
			actualCall.Op.ChildCallID = expectedCall.Op.ChildCallID
			Expect(*actualCall).To(Equal(*expectedCall))

		})
	})
	Context("callSpec.Parallel not empty", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScope := map[string]*model.Value{}
			providedID := "providedID"
			providedOpPath := "providedOpPath"
			providedParentIDValue := "providedParentID"
			providedParentID := &providedParentIDValue
			providedRootCallID := "providedRootCallID"
			providedDataDirPath, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			parallelSpec := []*model.CallSpec{}

			expectedCall := &model.Call{
				Parallel: parallelSpec,
				ID:       providedID,
				ParentID: providedParentID,
				RootID:   providedRootCallID,
			}

			/* act */
			actualCall, actualError := Interpret(
				context.Background(),
				providedScope,
				&model.CallSpec{
					Parallel: &parallelSpec,
				},
				providedID,
				providedOpPath,
				providedParentID,
				providedRootCallID,
				providedDataDirPath,
			)

			/* assert */
			Expect(actualError).To(BeNil())
			Expect(actualCall).To(Equal(expectedCall))

		})
	})
	Context("callSpec.Serial not empty", func() {
		It("should return expected result", func() {
			/* arrange */
			providedScope := map[string]*model.Value{}
			providedID := "providedID"
			providedOpPath := "providedOpPath"
			providedParentIDValue := "providedParentID"
			providedParentID := &providedParentIDValue
			providedRootCallID := "providedRootCallID"
			providedDataDirPath, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}

			serialSpec := []*model.CallSpec{}

			expectedCall := &model.Call{
				Serial:   serialSpec,
				ID:       providedID,
				ParentID: providedParentID,
				RootID:   providedRootCallID,
			}

			/* act */
			actualCall, actualError := Interpret(
				context.Background(),
				providedScope,
				&model.CallSpec{
					Serial: &serialSpec,
				},
				providedID,
				providedOpPath,
				providedParentID,
				providedRootCallID,
				providedDataDirPath,
			)

			/* assert */
			Expect(actualError).To(BeNil())
			Expect(*actualCall).To(Equal(*expectedCall))

		})
	})
})
