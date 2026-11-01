package node

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	runtimefakes "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Describe("YAML container output privacy", func() {
	It("activates beneath a YAML parent and protects renamed secret producers before execution", func() {
		dir, err := os.MkdirTemp("", "yaml-capture-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		child := filepath.Join(dir, "child")
		Expect(os.Mkdir(child, 0700)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(child, "op.yml"), []byte(`name: child
outputs:
  source: {string: {}}
run:
  container:
    image: {ref: alpine}
    cmd: [generate]
    stdout: $(source)
`), 0600)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(dir, "op.yml"), []byte(`name: parent
outputs:
  token: {string: {isSecret: true}}
run:
  serial:
    - op:
        ref: ./child
        outputs: {source: $(token)}
    - parallel:
        - container:
            image: {ref: alpine}
            cmd: [check]
            envVars: {TOKEN: $(token)}
        - container:
            image: {ref: alpine}
            cmd: [check]
            envVars: {TOKEN: $(token)}
`), 0600)).To(Succeed())
		db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLogger(nil))
		Expect(err).NotTo(HaveOccurred())
		defer db.Close()
		bus := pubsub.New(db)
		runtime := new(runtimefakes.FakeContainerRuntime)
		marker := "generated-private-capture-marker"
		var checks atomic.Int32
		runtime.RunContainerStub = func(_ context.Context, call *model.ContainerCall, _ string, _ pubsub.EventPublisher, out, stderr io.WriteCloser) (*int64, error) {
			defer out.Close()
			defer stderr.Close()
			if call.Cmd[0] == "generate" {
				out.Write([]byte(marker + "\n"))
			} else {
				if call.EnvVars["TOKEN"] != marker {
					return nil, fmt.Errorf("private value lost")
				}
				checks.Add(1)
			}
			return nil, nil
		}
		store := newStateStore(context.Background(), db, bus)
		clr := newCaller(newContainerCaller(runtime, bus, store), dir, bus, store)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = clr.Call(ctx, "root", nil, &model.CallSpec{Op: &model.OpCallSpec{Ref: dir}}, dir, nil, "root")
		Expect(err).NotTo(HaveOccurred())
		Expect(checks.Load()).To(Equal(int32(2)))
		since := time.Now().Add(-time.Minute)
		events, err := bus.Subscribe(ctx, model.EventFilter{Since: &since})
		Expect(err).NotTo(HaveOccurred())
		for event := range events {
			data, err := json.Marshal(event)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).NotTo(ContainSubstring(marker))
			Expect(string(data)).NotTo(ContainSubstring(base64.StdEncoding.EncodeToString([]byte(marker + "\n"))))
			if event.CallEnded != nil && event.CallEnded.Call.ID == "root" {
				break
			}
		}
	})
})
