package node

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

type cacheLockTestPubSub struct {
	pubsub.PubSub
	publish func(model.Event)
}

func (ps cacheLockTestPubSub) Publish(event model.Event) { ps.publish(event) }

type cacheLockTestRun struct {
	id     string
	finish chan error
}

var _ = Context("locked cache dirs", func() {
	var caller *_containerCaller
	var started chan cacheLockTestRun
	var events chan model.Event
	var runtime *fakes.FakeContainerRuntime
	var cacheA, cacheB string
	var ctx context.Context
	var cancel context.CancelFunc
	var startCall func(context.Context, string, []string, map[string]string) <-chan error

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		dataDir := GinkgoT().TempDir()
		cacheA, cacheB = filepath.Join(dataDir, "caches", "a"), filepath.Join(dataDir, "caches", "b")
		Expect(os.MkdirAll(cacheA, 0777)).To(Succeed())
		Expect(os.MkdirAll(cacheB, 0777)).To(Succeed())
		started = make(chan cacheLockTestRun, 10)
		events = make(chan model.Event, 10)
		runtime = new(fakes.FakeContainerRuntime)
		runtime.RunContainerStub = func(ctx context.Context, call *model.ContainerCall, _ string, _ pubsub.EventPublisher, stdout, stderr io.WriteCloser) (*int64, error) {
			defer stdout.Close()
			defer stderr.Close()
			run := cacheLockTestRun{id: call.ContainerID, finish: make(chan error, 1)}
			started <- run
			select {
			case err := <-run.finish:
				return nil, err
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
		ps := cacheLockTestPubSub{publish: func(event model.Event) { events <- event }}
		caller = newContainerCaller(runtime, ps, nil, dataDir).(*_containerCaller)
		startCall = func(callCtx context.Context, id string, paths []string, dirs map[string]string) <-chan error {
			result := make(chan error, 1)
			go func() {
				_, err := caller.Call(callCtx, &model.ContainerCall{
					BaseCall:    model.BaseCall{OpPath: "opPath"},
					ContainerID: id, Image: &model.ContainerCallImage{},
					LockedCacheDirs: paths, Dirs: dirs,
				}, nil, &model.ContainerCallSpec{}, "root-"+id)
				result <- err
			}()
			return result
		}
	})

	AfterEach(func() { cancel() })

	waitForRun := func() cacheLockTestRun {
		var run cacheLockTestRun
		Eventually(started).Should(Receive(&run))
		return run
	}
	waitForWaiter := func(id, name string) {
		var event model.Event
		Eventually(events).Should(Receive(&event))
		Expect(event.Timestamp).NotTo(BeZero())
		Expect(event.ContainerStdErrWrittenTo).NotTo(BeNil())
		Expect(string(event.ContainerStdErrWrittenTo.Data)).To(Equal("waiting for locked cache dir " + name + "\n"))
		Expect(event.ContainerStdErrWrittenTo.ContainerID).To(Equal(id))
		Expect(event.ContainerStdErrWrittenTo.OpRef).To(Equal("opPath"))
		Expect(event.ContainerStdErrWrittenTo.RootCallID).To(Equal("root-" + id))
		select {
		case run := <-started:
			Fail("waiting call reached the runtime: " + run.id)
		default:
		}
	}
	finishRun := func(run cacheLockTestRun, result <-chan error, runtimeErr error) {
		run.finish <- runtimeErr
		var err error
		Eventually(result).Should(Receive(&err))
		if runtimeErr == nil {
			Expect(err).NotTo(HaveOccurred())
		} else {
			Expect(err).To(Equal(runtimeErr))
		}
	}
	checkEmpty := func() {
		caller.cacheDirLocks.mu.Lock()
		defer caller.cacheDirLocks.mu.Unlock()
		Expect(caller.cacheDirLocks.byHostPath).To(BeEmpty())
	}

	It("serializes locked calls and reports one waiting line", func() {
		first := startCall(ctx, "first", []string{cacheA}, nil)
		firstRun := waitForRun()
		second := startCall(ctx, "second", []string{cacheA}, nil)
		waitForWaiter("second", "a")
		finishRun(firstRun, first, nil)
		secondRun := waitForRun()
		Expect(secondRun.id).To(Equal("second"))
		finishRun(secondRun, second, nil)
		Expect(events).To(BeEmpty())
		checkEmpty()
	})

	It("keeps the lock through the post-call permission walk", func() {
		cleanupReached, allowCleanup := make(chan struct{}), make(chan struct{})
		defer close(allowCleanup)
		caller.pubSub = cacheLockTestPubSub{publish: func(event model.Event) {
			if event.ContainerStdErrWrittenTo.ContainerID == "first" && strings.HasPrefix(string(event.ContainerStdErrWrittenTo.Data), "warning:") {
				close(cleanupReached)
				<-allowCleanup
				return
			}
			events <- event
		}}
		first := startCall(ctx, "first", []string{cacheA}, map[string]string{"/cache": filepath.Join(cacheA, "missing")})
		firstRun := waitForRun()
		firstRun.finish <- nil
		Eventually(cleanupReached).Should(BeClosed())
		second := startCall(ctx, "second", []string{cacheA}, nil)
		waitForWaiter("second", "a")
		allowCleanup <- struct{}{}
		Eventually(first).Should(Receive(BeNil()))
		finishRun(waitForRun(), second, nil)
		checkEmpty()
	})

	It("allows shared users alongside a locked holder and does not lock for shared holders", func() {
		locked := startCall(ctx, "locked", []string{cacheA}, nil)
		lockedRun := waitForRun()
		shared := startCall(ctx, "shared", nil, map[string]string{"/cache": cacheA})
		sharedRun := waitForRun()
		Expect(sharedRun.id).To(Equal("shared"))
		finishRun(lockedRun, locked, nil)
		next := startCall(ctx, "next", []string{cacheA}, nil)
		nextRun := waitForRun()
		Expect(nextRun.id).To(Equal("next"))
		finishRun(nextRun, next, nil)
		finishRun(sharedRun, shared, nil)
		Expect(events).To(BeEmpty())
		checkEmpty()
	})

	It("does not contend for different resolved paths including cold run dirs", func() {
		coldPath := filepath.Join(caller.dataDirPath, "nocache", "run", "a")
		first := startCall(ctx, "first", []string{cacheA}, nil)
		firstRun := waitForRun()
		second := startCall(ctx, "second", []string{cacheB}, nil)
		secondRun := waitForRun()
		cold := startCall(ctx, "cold", []string{coldPath}, nil)
		coldRun := waitForRun()
		finishRun(coldRun, cold, nil)
		finishRun(secondRun, second, nil)
		finishRun(firstRun, first, nil)
		Expect(events).To(BeEmpty())
		checkEmpty()
	})

	It("sorts multiple locks and deduplicates the same cache mounted twice", func() {
		holder := startCall(ctx, "holder", []string{cacheA}, nil)
		holderRun := waitForRun()
		multi := startCall(ctx, "multi", []string{cacheB, cacheA, cacheA}, map[string]string{"/one": cacheA, "/two": cacheA})
		waitForWaiter("multi", "a")
		// If multi took b first, this call would be blocked while a is held.
		other := startCall(ctx, "other", []string{cacheB}, nil)
		otherRun := waitForRun()
		Expect(otherRun.id).To(Equal("other"))
		finishRun(holderRun, holder, nil)
		waitForWaiter("multi", "b")
		finishRun(otherRun, other, nil)
		multiRun := waitForRun()
		Expect(multiRun.id).To(Equal("multi"))
		finishRun(multiRun, multi, nil)
		checkEmpty()
	})

	It("kills a waiting call promptly and releases any earlier locks", func() {
		holder := startCall(ctx, "holder", []string{cacheB}, nil)
		holderRun := waitForRun()
		waitCtx, kill := context.WithCancel(ctx)
		defer kill()
		waiting := startCall(waitCtx, "waiting", []string{cacheB, cacheA}, nil)
		waitForWaiter("waiting", "b")
		kill()
		Eventually(waiting).Should(Receive(Equal(context.Canceled)))
		Expect(runtime.RunContainerCallCount()).To(Equal(1))
		other := startCall(ctx, "other", []string{cacheA}, nil)
		finishRun(waitForRun(), other, nil)
		finishRun(holderRun, holder, nil)
		checkEmpty()
	})

	It("releases locks after a runtime error", func() {
		first := startCall(ctx, "first", []string{cacheA}, nil)
		firstRun := waitForRun()
		second := startCall(ctx, "second", []string{cacheA}, nil)
		waitForWaiter("second", "a")
		finishRun(firstRun, first, errors.New("runtime failed"))
		finishRun(waitForRun(), second, nil)
		checkEmpty()
	})

	It("releases locks when the running call is killed", func() {
		holderCtx, kill := context.WithCancel(ctx)
		defer kill()
		first := startCall(holderCtx, "first", []string{cacheA}, nil)
		waitForRun()
		second := startCall(ctx, "second", []string{cacheA}, nil)
		waitForWaiter("second", "a")
		kill()
		Eventually(first).Should(Receive(Equal(context.Canceled)))
		finishRun(waitForRun(), second, nil)
		checkEmpty()
	})
})
