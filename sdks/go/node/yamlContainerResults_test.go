package node

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	for _, mode := range []string{"stream", "outputs file", "declared outputs file"} {
		mode := mode
		It("protects renamed secret "+mode+" producers beneath a YAML parent", func() {
			dir, err := os.MkdirTemp("", "yaml-capture-")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(dir)
			child := filepath.Join(dir, "child")
			Expect(os.Mkdir(child, 0700)).To(Succeed())
			declaration := `name: child
outputs:
  source: {string: {}}
run:
  container:
    image: {ref: alpine}
    cmd: [generate]
    stdout: $(source)
`
			if mode != "stream" {
				declaration = strings.ReplaceAll(declaration, "stdout: $(source)", "outputs: {source: $(source)}")
			}
			if mode == "declared outputs file" {
				declaration += "    outputParams: {source: {string: {isSecret: true}}}\n"
			}
			Expect(os.WriteFile(filepath.Join(child, "op.yml"), []byte(declaration), 0600)).To(Succeed())
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
			if mode == "declared outputs file" {
				parentPath := filepath.Join(dir, "op.yml")
				data, err := os.ReadFile(parentPath)
				Expect(err).NotTo(HaveOccurred())
				Expect(os.WriteFile(parentPath, []byte(strings.ReplaceAll(string(data), "{string: {isSecret: true}}", "{string: {}}")), 0600)).To(Succeed())
			}
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
					if mode != "stream" {
						data, _ := json.Marshal(map[string]string{"source": marker})
						if err := os.WriteFile(call.OutputsFile, data, 0600); err != nil {
							return nil, err
						}
					} else {
						out.Write([]byte(marker + "\n"))
					}
				} else {
					if call.EnvVars["TOKEN"] != marker {
						return nil, fmt.Errorf("private value lost")
					}
					checks.Add(1)
				}
				return nil, nil
			}
			store := newStateStore(context.Background(), db, bus)
			clr := newCaller(newContainerCaller(runtime, bus, store), dir, bus, store, &goBootstrapCache{})
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
	}
})
