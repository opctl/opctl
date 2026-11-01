package node

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	runtimefakes "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

var _ = Describe("container outputs file", func() {
	It("copies runtime maps, supports atomic rename, binds only requested keys, and cleans its scratch directory", func() {
		dir, err := os.MkdirTemp("", "outputs-caller-")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)
		runtime := new(runtimefakes.FakeContainerRuntime)
		call := &model.ContainerCall{Image: &model.ContainerCallImage{}, OutputsFile: filepath.Join(dir, "outputs", "values.json"), EnvVars: map[string]string{"EXISTING": "yes"}, Dirs: map[string]string{"/ordinary": dir}}
		runtime.RunContainerStub = func(ctx context.Context, c *model.ContainerCall, _ string, _ pubsub.EventPublisher, out, stderr io.WriteCloser) (*int64, error) {
			defer out.Close()
			defer stderr.Close()
			Expect(containerruntime.SeparateOutputFrom(ctx)).To(BeFalse())
			Expect(c.EnvVars["OPCTL_OUTPUTS"]).To(Equal("/opctl/outputs/values.json"))
			Expect(c.Dirs["/opctl/outputs"]).To(Equal(filepath.Join(dir, "outputs")))
			stat, err := os.Stat(c.Dirs["/opctl/outputs"])
			Expect(err).NotTo(HaveOccurred())
			Expect(stat.Mode() & (os.ModeSticky | 0777)).To(Equal(os.ModeSticky | os.FileMode(0777)))
			Expect(os.WriteFile(c.OutputsFile+".tmp", []byte(`{"version":"$(literal)\n","count":3,"ignored":true}`), 0600)).To(Succeed())
			Expect(os.Rename(c.OutputsFile+".tmp", c.OutputsFile)).To(Succeed())
			return nil, nil
		}
		cc := newContainerCaller(runtime, new(fakes.FakePubSub), &_stateStore{})
		result, err := cc.Call(context.Background(), call, nil, &model.ContainerCallSpec{Outputs: map[string]string{"version": "$(renamed)", "count": "$(count)"}}, "root")
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(HaveLen(2))
		Expect(*result["renamed"].String).To(Equal("$(literal)\n"))
		Expect(*result["count"].Number).To(Equal(float64(3)))
		Expect(call.EnvVars).To(Equal(map[string]string{"EXISTING": "yes"}))
		Expect(call.Dirs).To(HaveLen(1))
		_, err = os.Stat(filepath.Dir(call.OutputsFile))
		Expect(os.IsNotExist(err)).To(BeTrue())
		_, err = os.Stat(dir)
		Expect(err).NotTo(HaveOccurred())
	})
	It("collects atomically and preserves runtime errors without exposing malformed values", func() {
		for _, data := range []string{"", `{"other":1}`, `{"secret-marker":}`, `{"v":null}`} {
			for _, failed := range []bool{false, true} {
				dir, err := os.MkdirTemp("", "outputs-failure-")
				Expect(err).NotTo(HaveOccurred())
				runtime := new(runtimefakes.FakeContainerRuntime)
				call := &model.ContainerCall{Image: &model.ContainerCallImage{}, OutputsFile: filepath.Join(dir, "outputs", "values.json")}
				runtime.RunContainerStub = func(_ context.Context, c *model.ContainerCall, _ string, _ pubsub.EventPublisher, out, stderr io.WriteCloser) (*int64, error) {
					defer out.Close()
					defer stderr.Close()
					if data != "" {
						os.WriteFile(c.OutputsFile, []byte(data), 0600)
					}
					if failed {
						return nil, errors.New("primary runtime error")
					}
					return nil, nil
				}
				result, err := newContainerCaller(runtime, new(fakes.FakePubSub), &_stateStore{}).Call(context.Background(), call, nil, &model.ContainerCallSpec{Outputs: map[string]string{"v": "$(v)"}}, "root")
				Expect(result).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).NotTo(ContainSubstring("secret-marker"))
				if failed {
					Expect(err).To(MatchError("primary runtime error"))
				}
				_, err = os.Stat(filepath.Dir(call.OutputsFile))
				Expect(os.IsNotExist(err)).To(BeTrue())
				os.RemoveAll(dir)
			}
		}
	})
})
