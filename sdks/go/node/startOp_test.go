package node

import (
	"context"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/internal/fakes"
	. "github.com/opctl/opctl/sdks/go/node/pubsub/fakes"
)

var _ = Context("core", func() {
	Context("StartOp", func() {
		Context("data.Resolve errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedCtx := context.Background()
				providedStartOpReq := model.StartOpReq{
					Op: model.StartOpReqOp{
						Ref: "dummyOpRef",
					},
				}

				objectUnderTest := core{}

				/* act */
				_, actualErr := objectUnderTest.StartOp(
					providedCtx,
					providedStartOpReq,
				)

				/* assert */
				Expect(actualErr).NotTo(BeNil())
			})
		})
		Context("data.Resolve doesn't err", func() {
			Context("opfile.Get doesn't err", func() {
				It("should call caller.Call w/ expected args", func() {
					/* arrange */
					providedCtx := context.Background()
					providedArg1String := "dummyArg1Value"
					providedArg2Dir := "/"
					providedArg3Dir := "/"
					providedArg4Dir := "/"

					// use local op
					wd, err := os.Getwd()
					if err != nil {
						panic(err)
					}
					providedOpPath := filepath.Join(wd, "testdata/startOp")
					providedReq := model.StartOpReq{
						Args: map[string]*ipld.Node{
							"dummyArg1Name": {String: &providedArg1String},
							"dummyArg2Name": {Dir: &providedArg2Dir},
							"dummyArg3Name": {Dir: &providedArg3Dir},
							"dummyArg4Name": {Dir: &providedArg4Dir},
						},
						Op: model.StartOpReqOp{
							Ref: providedOpPath,
						},
					}

					opFile := &model.OpSpec{
						Outputs: map[string]*model.ParamSpec{
							"dummyOutput1": {String: &model.StringParamSpec{}},
							"dummyOutput2": {String: &model.StringParamSpec{}},
						},
					}

					expectedOpCallSpec := &model.OpCallSpec{
						Ref:     providedOpPath,
						Inputs:  map[string]interface{}{},
						Outputs: map[string]string{},
					}
					for name := range providedReq.Args {
						expectedOpCallSpec.Inputs[name] = ""
					}
					for name := range opFile.Outputs {
						expectedOpCallSpec.Outputs[name] = ""
					}

					fakeCaller := new(FakeCaller)
					dataCachePath, err := os.MkdirTemp("", "")
					if err != nil {
						panic(err)
					}

					objectUnderTest := core{
						caller:        fakeCaller,
						dataCachePath: dataCachePath,
						pubSub:        new(FakePubSub),
					}

					/* act */
					objectUnderTest.StartOp(
						providedCtx,
						providedReq,
					)

					/* assert */
					// Call happens in go routine; wait 500ms to allow it to occur
					time.Sleep(time.Millisecond * 500)
					_,
						actualOpID,
						actualScope,
						actualCallSpec,
						actualOpPath,
						_,
						actualRootID := fakeCaller.CallArgsForCall(0)

					Expect(actualOpID).To(HaveLen(32))
					Expect(actualScope).To(Equal(providedReq.Args))
					Expect(*actualCallSpec).To(BeEquivalentTo(model.CallSpec{
						Op: expectedOpCallSpec,
					}))
					Expect(actualOpPath).To(Equal(providedOpPath))
					Expect(actualRootID).To(HaveLen(32))
				})
			})
		})
	})
})
