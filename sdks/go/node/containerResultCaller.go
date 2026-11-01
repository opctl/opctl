package node

import (
	"context"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"sort"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/sensitivity"
)

func (cc _containerCaller) callWithResults(ctx context.Context, call *model.ContainerCall, scope map[string]*model.Value, spec *model.ContainerCallSpec, root string, declarations map[string]*model.ParamSpec) (result map[string]*model.Value, resultErr error) {
	if err := opspec.ValidateContainerOutputParams(declarations, spec.Outputs); err != nil {
		return nil, err
	}
	if stdoutCaptureFrom(ctx) != nil {
		return nil, fmt.Errorf("public container output bindings are not allowed on a private bootstrap call")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if spec.Stdout != "" || spec.Stderr != "" {
		ctx = containerruntime.WithSeparateOutput(ctx)
	}
	runtimeCall := *call
	image := *call.Image
	runtimeCall.Image = &image
	if !image.Prepared && image.Ref != nil && image.PullCreds == nil {
		if auth := cc.stateStore.TryGetAuth(*image.Ref); auth != nil {
			image.PullCreds = &auth.Creds
		}
	}
	if len(spec.Outputs) > 0 {
		if call.OutputsFile == "" {
			return nil, fmt.Errorf("container outputs file path was not prepared by the interpreter")
		}
		dir := filepath.Dir(call.OutputsFile)
		if err := unsudo.CreateDir(dir); err != nil {
			return nil, fmt.Errorf("unable to prepare outputs directory")
		}
		defer func() {
			if err := os.RemoveAll(dir); err != nil && resultErr == nil {
				result = nil
				resultErr = fmt.Errorf("unable to remove outputs directory")
			}
		}()
		// Only the file crosses the mount boundary; root-owned siblings stay in the container.
		if err := unsudo.CreateFile(call.OutputsFile, nil); err != nil {
			return nil, fmt.Errorf("unable to prepare outputs file")
		}
		if err := os.Chmod(call.OutputsFile, 0666); err != nil {
			return nil, fmt.Errorf("unable to make outputs file writable")
		}
		runtimeCall.Files = maps.Clone(call.Files)
		if runtimeCall.Files == nil {
			runtimeCall.Files = map[string]string{}
		}

		runtimeCall.EnvVars = maps.Clone(call.EnvVars)
		if runtimeCall.EnvVars == nil {
			runtimeCall.EnvVars = map[string]string{}
		}
		runtimeCall.Files["/opctl/outputs/values.json"] = call.OutputsFile
		runtimeCall.EnvVars["OPCTL_OUTPUTS"] = "/opctl/outputs/values.json"
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
	for _, stream := range streams {
		stream.Lock()
		captureErr := stream.err
		stream.Unlock()
		if captureErr != nil {
			return nil, captureErr
		}
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
	values := map[string]*model.Value{}
	streamValues := map[string]*model.Value{}
	streamParams := map[string]*model.ParamSpec{}
	for _, binding := range []struct{ name, ref string }{{"stdout", spec.Stdout}, {"stderr", spec.Stderr}} {
		if binding.ref == "" {
			continue
		}
		stream := streams[binding.ref]
		if stream == nil {
			return nil, fmt.Errorf("container %s capture result is missing", binding.name)
		}
		value, captureErr := stream.value()
		if captureErr != nil {
			return nil, captureErr
		}
		value = sensitivity.Mark(value, secretDestination(ctx, scope, binding.ref))
		streamValues[binding.name] = value
		streamParams[binding.name] = &model.ParamSpec{String: &model.StringParamSpec{}}
		values[opspec.RefToName(binding.ref)] = value
	}
	if err := params.Validate(streamValues, streamParams, params.ValidationOptions{ExactTypes: true}); err != nil {
		return nil, fmt.Errorf("container output %w", err)
	}
	if len(spec.Outputs) > 0 {
		named, err := readContainerOutputs(call.OutputsFile)
		if err != nil {
			return nil, err
		}
		for name, value := range named {
			private := secretParam(declarations[name])
			if ref, ok := spec.Outputs[name]; ok {
				private = private || secretDestination(ctx, scope, ref)
			}
			named[name] = sensitivity.Mark(value, private)
		}
		if err := params.Validate(named, declarations, params.ValidationOptions{ExactTypes: true}); err != nil {
			return nil, fmt.Errorf("container output %w", err)
		}
		for _, name := range sortedOutputNames(spec.Outputs) {
			value, ok := named[name]
			if !ok {
				return nil, fmt.Errorf("outputs file is missing requested member %q", name)
			}
			values[opspec.RefToName(spec.Outputs[name])] = value
		}
	}
	outputs := cc.interpretOutputs(spec, call)
	for name, value := range values {
		outputs[name] = value
	}
	return outputs, nil
}

func sortedOutputNames(bindings map[string]string) []string {
	keys := make([]string, 0, len(bindings))
	for key := range bindings {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
