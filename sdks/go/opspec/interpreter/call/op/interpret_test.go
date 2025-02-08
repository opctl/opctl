package op

import (
	"context"
	"encoding/json"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

var _ = Context("Interpret", func() {
	It("should return expected result", func() {
		/* arrange */
		providedOpId := "opID"
		dataDir, err := os.MkdirTemp("", "")
		if err != nil {
			panic(err)
		}

		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		opRef := "testdata/interpret"

		expectedResult := model.OpCall{
			BaseCall: model.BaseCall{
				OpPath: path.Join(wd, opRef),
			},
			ChildCallCallSpec: &model.CallSpec{
				Container: &model.ContainerCallSpec{
					Image: &model.ContainerCallImageSpec{
						Ref: "ghcr.io/linuxcontainers/alpine",
					},
				},
			},
			Inputs: map[string]*model.Value{},
			OpID:   providedOpId,
		}

		/* act */
		actualResult, actualError := Interpret(
			context.Background(),
			map[string]*model.Value{},
			&model.OpCallSpec{
				Ref: opRef,
			},
			providedOpId,
			wd,
			dataDir,
		)

		/* assert */
		Expect(actualError).To((BeNil()))

		// ignore ChildCallID; it's dynamic
		expectedResult.ChildCallID = actualResult.ChildCallID

		// compare as JSON; otherwise we encounter pointer inequalities
		actualBytes, err := json.Marshal(actualResult)
		if err != nil {
			panic(err)
		}

		expectedBytes, err := json.Marshal(expectedResult)
		if err != nil {
			panic(err)
		}

		Expect(actualBytes).To(Equal(expectedBytes))
	})
})
