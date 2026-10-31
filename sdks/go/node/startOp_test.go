package node

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/internal/runcontext"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/internal/fakes"
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

	Context("StartOp", func() {
		It("should propagate no-cache through nested ops, serial, parallel, and loops without affecting concurrent runs", func() {
			dataDirPath := GinkgoT().TempDir()
			opPath := GinkgoT().TempDir()
			nestedOpPath := filepath.Join(opPath, "nested")
			Expect(os.Mkdir(nestedOpPath, 0700)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(nestedOpPath, "op.yml"), []byte(`name: nested
run:
  parallel:
    - container:
        image: { ref: node:22 }
        caches: { /cache: shared }
    - container:
        image: { ref: node:22 }
        caches: { /cache: shared }
`), 0600)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(opPath, "op.yml"), []byte(fmt.Sprintf(`name: cache-run
run:
  serial:
    - container:
        image: { ref: node:22 }
        caches: { /cache: shared }
    - op: { ref: %q }
    - serialLoop:
        range: [0, 1]
        run:
          container:
            image: { ref: node:22 }
            caches: { /cache: shared }
    - parallelLoop:
        range: [0, 1]
        run:
          container:
            image: { ref: node:22 }
            caches: { /cache: shared }
`, nestedOpPath)), 0600)).To(Succeed())

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			pubSub := pubsub.New(db)
			stateStore := newStateStore(ctx, db, pubSub)
			coldCalls := make(chan *model.ContainerCall, 7)
			warmCalls := make(chan *model.ContainerCall, 7)
			releaseCold := make(chan struct{})
			fakeContainerCaller := new(FakeContainerCaller)
			fakeContainerCaller.CallStub = func(ctx context.Context, call *model.ContainerCall, scope map[string]*model.Value, spec *model.ContainerCallSpec, rootID string) (map[string]*model.Value, error) {
				if runcontext.NoCacheRoot(ctx) != "" {
					coldCalls <- call
					select {
					case <-releaseCold:
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				} else {
					warmCalls <- call
				}
				return map[string]*model.Value{}, nil
			}
			objectUnderTest := core{
				caller:        newCaller(fakeContainerCaller, dataDirPath, pubSub, stateStore),
				dataCachePath: filepath.Join(dataDirPath, "ops"),
				dataDirPath:   dataDirPath,
				pubSub:        pubSub,
				stateStore:    stateStore,
			}
			waitForEnd := func(rootID string) {
				events, err := pubSub.Subscribe(ctx, model.EventFilter{Roots: []string{rootID}})
				Expect(err).To(BeNil())
				for {
					var event model.Event
					Eventually(events, 5*time.Second).Should(Receive(&event))
					if event.CallEnded != nil && event.CallEnded.Call.ID == rootID {
						Expect(event.CallEnded.Outcome).To(Equal(model.OpOutcomeSucceeded), "%+v", event.CallEnded.Error)
						return
					}
				}
			}

			coldID, err := objectUnderTest.StartOp(ctx, model.StartOpReq{
				NoCache: true,
				Op:      model.StartOpReqOp{Ref: opPath},
			})
			Expect(err).To(BeNil())
			var firstColdCall *model.ContainerCall
			Eventually(coldCalls).Should(Receive(&firstColdCall))
			Expect(filepath.Join(dataDirPath, "caches")).NotTo(BeAnExistingFile())

			// Complete a normal run while the cold run is still in its first container.
			warmID, err := objectUnderTest.StartOp(ctx, model.StartOpReq{Op: model.StartOpReqOp{Ref: opPath}})
			Expect(err).To(BeNil())
			waitForEnd(warmID)
			cachePath := filepath.Join(dataDirPath, "caches", "shared")
			for i := 0; i < 7; i++ {
				var call *model.ContainerCall
				Eventually(warmCalls).Should(Receive(&call))
				Expect(call.Dirs["/cache"]).To(Equal(cachePath))
			}
			cacheFilePath := filepath.Join(cachePath, "contents")
			Expect(os.WriteFile(cacheFilePath, []byte("warm"), 0600)).To(Succeed())

			coldRootPath := filepath.Join(dataDirPath, "nocache", coldID)
			entries, err := os.ReadDir(firstColdCall.Dirs["/cache"])
			Expect(err).To(BeNil())
			Expect(entries).To(BeEmpty())
			close(releaseCold)
			waitForEnd(coldID)
			for i := 0; i < 7; i++ {
				call := firstColdCall
				if i > 0 {
					Eventually(coldCalls).Should(Receive(&call))
				}
				// a cold run shares one empty dir per cache name, like the cache root
				Expect(call.Dirs["/cache"]).To(Equal(filepath.Join(coldRootPath, "shared")))
			}
			Eventually(coldRootPath).ShouldNot(BeAnExistingFile())
			contents, err := os.ReadFile(cacheFilePath)
			Expect(err).To(BeNil())
			Expect(string(contents)).To(Equal("warm"))
		})

		Context("data.Resolve errs", func() {
			It("should return expected result", func() {
				/* arrange */
				providedCtx := context.Background()
				providedStartOpReq := model.StartOpReq{
					Op: model.StartOpReqOp{
						Ref: "dummyOpRef",
					},
				}

				objectUnderTest := core{
					stateStore: newStateStore(
						context.Background(),
						db,
						pubsub.New(db),
					),
				}

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
						Args: map[string]*model.Value{
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

					pubSub := pubsub.New(db)

					objectUnderTest := core{
						caller:        fakeCaller,
						dataCachePath: dataCachePath,
						pubSub:        pubSub,
						stateStore: newStateStore(
							context.Background(),
							db,
							pubSub,
						),
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
