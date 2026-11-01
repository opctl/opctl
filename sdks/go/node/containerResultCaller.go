package node

import (
	"context"
	"fmt"
	"io"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/sensitivity"
)

func (cc _containerCaller) callWithResults(ctx context.Context, call *model.ContainerCall, scope map[string]*model.Value, spec *model.ContainerCallSpec, root string) (map[string]*model.Value, error) {
	if stdoutCaptureFrom(ctx) != nil {
		return nil, fmt.Errorf("public container output bindings are not allowed on a private bootstrap call")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = containerruntime.WithSeparateOutput(ctx)
	runtimeCall := *call
	image := *call.Image
	runtimeCall.Image = &image
	if !image.Prepared && image.Ref != nil && image.PullCreds == nil {
		if auth := cc.stateStore.TryGetAuth(*image.Ref); auth != nil {
			image.PullCreds = &auth.Creds
		}
	}
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()
	logs := make(chan error, 1)
	go func() { logs <- cc.interpretLogs(outReader, errReader, call, root) }()
	streams := map[string]*resultStream{}
	writers := []io.WriteCloser{outWriter, errWriter}
	for i, binding := range []struct{ name, ref string }{{"stdout", spec.Stdout}, {"stderr", spec.Stderr}} {
		if binding.ref == "" {
			continue
		}
		buffer := newResultStream(binding.name, cancel)
		streams[binding.ref] = buffer
		var writer io.Writer = io.MultiWriter(buffer, writers[i])
		if secretDestination(ctx, scope, binding.ref) {
			writer = buffer
		}
		writers[i] = teeCloser{Writer: writer, Closer: writers[i]}
	}
	code, err := cc.containerRuntime.RunContainer(ctx, &runtimeCall, root, cc.pubSub, writers[0], writers[1])
	// Runtimes normally close these; also close on failed starts or custom runtimes.
	writers[0].Close()
	writers[1].Close()
	logErr := <-logs
	values := map[string]*model.Value{}
	for _, ref := range []string{spec.Stdout, spec.Stderr} {
		if ref == "" {
			continue
		}
		value, captureErr := streams[ref].value()
		if captureErr != nil {
			return nil, captureErr
		}
		values[opspec.RefToName(ref)] = sensitivity.Mark(value, secretDestination(ctx, scope, ref))
	}
	if err != nil {
		return nil, err
	}
	if code != nil && *code != 0 {
		return nil, fmt.Errorf("nonzero container exit code: %d", *code)
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if logErr != nil {
		return nil, logErr
	}
	outputs := cc.interpretOutputs(spec, call)
	for name, value := range values {
		outputs[name] = value
	}
	return outputs, nil
}
