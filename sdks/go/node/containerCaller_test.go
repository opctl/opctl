package node

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	. "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Context("containerCaller", func() {
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

	Context("newContainerCaller", func() {
		It("should return containerCaller", func() {
			/* arrange/act/assert */
			pubSub := pubsub.New(db)

			Expect(newContainerCaller(
				new(FakeContainerRuntime),
				pubSub,
				newStateStore(
					context.Background(),
					db,
					pubSub,
				),
				dbDir,
			)).To(Not(BeNil()))
		})
	})
	Context("Call", func() {
		for _, testCase := range []struct {
			name        string
			exitCode    int64
			runtimeErr  error
			kill        bool
			expectedErr string
		}{
			{name: "success"},
			{name: "nonzero exit", exitCode: 2, expectedErr: "nonzero container exit code: 2"},
			{name: "runtime error", runtimeErr: errors.New("runtime failed"), expectedErr: "runtime failed"},
			{name: "kill", runtimeErr: context.Canceled, kill: true, expectedErr: context.Canceled.Error()},
		} {
			It("should restore only mounted cache dirs after "+testCase.name+" when privileged", func() {
				dataDirPath := GinkgoT().TempDir()
				cacheRootPath := filepath.Join(dataDirPath, "caches")
				cachePath := filepath.Join(cacheRootPath, "npm")
				secondCachePath := filepath.Join(cacheRootPath, "gomod")
				unmountedCachePath := filepath.Join(cacheRootPath, "unmounted")
				siblingPath := filepath.Join(dataDirPath, "caches-other", "dir")
				scratchCachePath := filepath.Join(dataDirPath, "nocache", "rootCallID", "npm")
				outsidePath := cacheRootPath + string(filepath.Separator) + "../outside"
				missingPath := filepath.Join(cacheRootPath, "missing")
				GinkgoT().Setenv("SUDO_UID", "1234")
				GinkgoT().Setenv("SUDO_GID", "2345")

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				fakeContainerRuntime := new(FakeContainerRuntime)
				fakeContainerRuntime.RunContainerStub = func(
					ctx context.Context,
					req *model.ContainerCall,
					rootCallID string,
					eventPublisher pubsub.EventPublisher,
					stdOut io.WriteCloser,
					stdErr io.WriteCloser,
				) (*int64, error) {
					defer stdOut.Close()
					defer stdErr.Close()
					for _, dirPath := range []string{cachePath, secondCachePath, unmountedCachePath, siblingPath, outsidePath, scratchCachePath} {
						Expect(os.MkdirAll(filepath.Join(dirPath, "nested"), 0700)).To(Succeed())
						Expect(os.WriteFile(filepath.Join(dirPath, "nested", "file"), nil, 0600)).To(Succeed())
					}
					if testCase.kill {
						cancel()
					}
					return &testCase.exitCode, testCase.runtimeErr
				}
				pubSub := pubsub.New(db)
				rootCallID := "cache-ownership-" + testCase.name
				objectUnderTest := newContainerCaller(fakeContainerRuntime, pubSub, nil, dataDirPath)

				_, actualErr := objectUnderTest.Call(
					ctx,
					&model.ContainerCall{
						BaseCall:    model.BaseCall{OpPath: "opPath"},
						ContainerID: "containerID",
						Image:       &model.ContainerCallImage{},
						Dirs: map[string]string{
							"/cache":      cachePath,
							"/second":     secondCachePath,
							"/cache-root": cacheRootPath,
							"/sibling":    siblingPath,
							"/outside":    outsidePath,
							"/missing":    missingPath,
							"/cold":       scratchCachePath,
						},
					},
					nil,
					&model.ContainerCallSpec{},
					rootCallID,
				)

				if testCase.expectedErr == "" {
					Expect(actualErr).To(BeNil())
				} else {
					Expect(actualErr).To(MatchError(testCase.expectedErr))
				}
				for _, dirPath := range []string{cachePath, secondCachePath, unmountedCachePath, siblingPath, outsidePath, scratchCachePath} {
					uid, gid := os.Geteuid(), os.Getegid()
					if uid == 0 && (dirPath == cachePath || dirPath == secondCachePath) {
						uid, gid = 1234, 2345
					}
					for _, entryPath := range []string{dirPath, filepath.Join(dirPath, "nested"), filepath.Join(dirPath, "nested", "file")} {
						info, err := os.Stat(entryPath)
						Expect(err).To(BeNil())
						stat := info.Sys().(*syscall.Stat_t)
						Expect(stat.Uid).To(Equal(uint32(uid)), entryPath)
						Expect(stat.Gid).To(Equal(uint32(gid)), entryPath)
					}
				}

				eventCtx, eventCancel := context.WithCancel(context.Background())
				defer eventCancel()
				events, err := pubSub.Subscribe(eventCtx, model.EventFilter{Roots: []string{rootCallID}})
				Expect(err).To(BeNil())
				if os.Geteuid() == 0 {
					var event model.Event
					Eventually(events).Should(Receive(&event))
					Expect(event.Timestamp).NotTo(BeZero())
					Expect(event.ContainerStdErrWrittenTo).NotTo(BeNil())
					Expect(string(event.ContainerStdErrWrittenTo.Data)).To(ContainSubstring("unable to restore cache dir ownership"))
					Expect(string(event.ContainerStdErrWrittenTo.Data)).To(ContainSubstring(missingPath))
					Expect(event.ContainerStdErrWrittenTo.ContainerID).To(Equal("containerID"))
					Expect(event.ContainerStdErrWrittenTo.OpRef).To(Equal("opPath"))
					Expect(event.ContainerStdErrWrittenTo.RootCallID).To(Equal(rootCallID))
				}
				Consistently(events).ShouldNot(Receive())
			})
		}

		It("should call containerRuntime.RunContainer w/ expected args", func() {
			/* arrange */
			providedCtx := context.Background()
			providedContainerCall := &model.ContainerCall{
				BaseCall: model.BaseCall{},
				Image:    &model.ContainerCallImage{},
			}
			providedRootCallID := "providedRootCallID"
			fakeContainerRuntime := new(FakeContainerRuntime)

			fakeContainerRuntime.RunContainerStub = func(
				ctx context.Context,
				req *model.ContainerCall,
				rootCallID string,
				eventPublisher pubsub.EventPublisher,
				stdOut io.WriteCloser,
				stdErr io.WriteCloser,
			) (*int64, error) {

				stdErr.Close()
				stdOut.Close()

				return nil, nil
			}

			pubSub := pubsub.New(db)

			objectUnderTest := _containerCaller{
				containerRuntime: fakeContainerRuntime,
				pubSub:           pubSub,
			}

			/* act */
			objectUnderTest.Call(
				providedCtx,
				providedContainerCall,
				map[string]*model.Value{},
				&model.ContainerCallSpec{},
				providedRootCallID,
			)

			/* assert */
			_,
				actualContainerCall,
				actualRootCallID,
				actualEventPublisher,
				_,
				_ := fakeContainerRuntime.RunContainerArgsForCall(0)
			Expect(actualContainerCall).To(Equal(providedContainerCall))
			Expect(actualRootCallID).To(Equal(providedRootCallID))
			Expect(actualEventPublisher).To(Equal(pubSub))
		})
		Context("containerRuntime.RunContainer errors", func() {
			It("should publish expected ContainerExited", func() {
				/* arrange */
				expectedErrorMessage := "expectedErrorMessage"

				fakeContainerRuntime := new(FakeContainerRuntime)

				fakeContainerRuntime.RunContainerStub = func(
					ctx context.Context,
					req *model.ContainerCall,
					rootCallID string,
					eventPublisher pubsub.EventPublisher,
					stdOut io.WriteCloser,
					stdErr io.WriteCloser,
				) (*int64, error) {

					stdErr.Close()
					stdOut.Close()

					return nil, errors.New(expectedErrorMessage)
				}

				objectUnderTest := _containerCaller{
					containerRuntime: fakeContainerRuntime,
					pubSub:           pubsub.New(db),
				}

				/* act */
				actualOutputs, actualErr := objectUnderTest.Call(
					context.Background(),
					&model.ContainerCall{
						BaseCall: model.BaseCall{},
						Image:    &model.ContainerCallImage{},
					},
					map[string]*model.Value{},
					&model.ContainerCallSpec{},
					"rootCallID",
				)

				/* assert */
				Expect(actualOutputs).To(Equal(map[string]*model.Value{}))
				Expect(actualErr).To(MatchError(expectedErrorMessage))
			})
		})
	})

	It("should return expected results", func() {
		/* arrange */
		providedOpPath := "providedOpPath"
		providedContainerCall := &model.ContainerCall{
			BaseCall: model.BaseCall{
				OpPath: providedOpPath,
			},
			ContainerID: "providedContainerID",
			Image:       &model.ContainerCallImage{},
		}
		providedInboundScope := map[string]*model.Value{}
		providedContainerCallSpec := &model.ContainerCallSpec{}

		fakeContainerRuntime := new(FakeContainerRuntime)

		expectedErr := errors.New("io: read/write on closed pipe")
		fakeContainerRuntime.RunContainerStub = func(
			ctx context.Context,
			req *model.ContainerCall,
			rootCallID string,
			eventPublisher pubsub.EventPublisher,
			stdOut io.WriteCloser,
			stdErr io.WriteCloser,
		) (*int64, error) {

			stdErr.Close()
			stdOut.Close()

			return nil, expectedErr
		}

		objectUnderTest := _containerCaller{
			containerRuntime: fakeContainerRuntime,
			pubSub:           pubsub.New(db),
		}

		/* act */
		actualOutputs, actualErr := objectUnderTest.Call(
			context.Background(),
			providedContainerCall,
			providedInboundScope,
			providedContainerCallSpec,
			"rootCallID",
		)

		/* assert */
		Expect(actualOutputs).To(Equal(map[string]*model.Value{}))
		Expect(actualErr).To(Equal(expectedErr))
	})
})
