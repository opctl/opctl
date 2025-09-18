package containerd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"syscall"

	containerd "github.com/containerd/containerd/v2/client"
	dkr "github.com/containerd/containerd/v2/core/remotes/docker"
	"github.com/containerd/containerd/v2/pkg/cio"
	"github.com/containerd/containerd/v2/pkg/oci"
	"github.com/containerd/containerd/v2/plugins"
	"github.com/nxadm/tail"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type runContainer interface {
	RunContainer(
		ctx context.Context,
		req *model.ContainerCall,
		rootCallID string,
		eventPublisher pubsub.EventPublisher,
		stdout io.WriteCloser,
		stderr io.WriteCloser,
	) (*int64, error)
}

type _runContainer struct {
	client *containerd.Client
}

func newRunContainer(
	client *containerd.Client,
) (runContainer, error) {
	return _runContainer{
		client: client,
	}, nil
}

func (rc _runContainer) RunContainer(
	ctx context.Context,
	req *model.ContainerCall,
	rootCallID string,
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	defer stdout.Close()
	defer stderr.Close()

	if req.Image == nil || (req.Image.Ref == nil && req.Image.Src == nil) {
		return nil, fmt.Errorf("containerd: image not specified")
	}

	// Resolve image
	var (
		img containerd.Image
		err error
	)

	if req.Image.Src != nil {
		// @TODO: support importing from OCI layout dir similar to docker pushImage
		// For now, require Ref be provided or image present locally.
		return nil, fmt.Errorf("containerd: Image.Src import not implemented")
	}

	// Pull by ref (or get if already present)
	if req.Image.Ref != nil {
		resolverOpts := dkr.ResolverOptions{}
		if req.Image.PullCreds != nil && req.Image.PullCreds.Username != "" {
			resolverOpts.Hosts = dkr.ConfigureDefaultRegistries(
				dkr.WithAuthorizer(
					dkr.NewDockerAuthorizer(
						dkr.WithAuthCreds(
							func(s string) (string, string, error) {
								return req.Image.PullCreds.Username, req.Image.PullCreds.Password, nil
							},
						),
					),
				),
			)
		}

		pullOpts := []containerd.RemoteOpt{containerd.WithPullUnpack}
		if resolverOpts.Hosts != nil {
			pullOpts = append(pullOpts, containerd.WithResolver(dkr.NewResolver(resolverOpts)))
		}
		// platform arch if provided (assume linux)
		if req.Image.Platform != nil && req.Image.Platform.Arch != nil && *req.Image.Platform.Arch != "" {
			pullOpts = append(pullOpts, containerd.WithPlatform(fmt.Sprintf("linux/%s", *req.Image.Platform.Arch)))
		}

		img, err = rc.client.Pull(ctx, *req.Image.Ref, pullOpts...)
		if err != nil {
			// allow cached/offline use
			img, err = rc.client.GetImage(ctx, *req.Image.Ref)
			if err != nil {
				return nil, err
			}
		}

	} else {
		return nil, fmt.Errorf("containerd: image ref required")
	}

	name := getContainerName(req.ContainerID)

	// Build env
	env := make([]string, 0, len(req.EnvVars))
	for k, v := range req.EnvVars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(env)

	logFilePath := path.Join(
		"/Users/chrisd/Library/Application Support/opctl/vms/default/",
		req.ContainerID+".txt",
	)

	f, err := os.Create(logFilePath)
	if err != nil {
		return nil, err
	}

	f.Close()

	// Mounts
	mounts := []specs.Mount{
		// logfile has to be mounted
		{
			Type:        "bind",
			Source:      logFilePath,
			Destination: logFilePath,
			Options:     []string{"rbind", "rw"},
		},
	}
	for containerPath, hostPath := range req.Files {
		mounts = append(mounts, specs.Mount{
			Type:        "bind",
			Source:      hostPath,
			Destination: containerPath,
			Options:     []string{"rbind", "rw"},
		})
	}
	for containerPath, hostPath := range req.Dirs {
		mounts = append(mounts, specs.Mount{
			Type:        "bind",
			Source:      hostPath,
			Destination: containerPath,
			Options:     []string{"rbind", "rw"},
		})
	}
	for containerSocketPath, hostSocketAddr := range req.Sockets {
		// naive detection of unix sockets (same as docker impl approach)
		if strings.Contains(hostSocketAddr, "/") || strings.Contains(hostSocketAddr, `\`) {
			mounts = append(mounts, specs.Mount{
				Type:        "bind",
				Source:      hostSocketAddr,
				Destination: containerSocketPath,
				Options:     []string{"rbind"},
			})
		}
	}

	// Ports are not managed directly by containerd (requires CNI). Emit a warning.
	if len(req.Ports) > 0 {
		_, _ = stdout.Write([]byte("warning: containerd runtime does not publish ports; req.Ports will be ignored\n"))
	}

	// Spec options
	specOpts := []oci.SpecOpts{
		oci.WithDefaultSpecForPlatform("linux"),
		oci.WithImageConfig(img),
		oci.WithEnv(env),
		// Disable cgroup hierarchy and make cgroup filesystem writable to avoid permission issues
		// oci.WithCgroup(""),
		oci.WithWriteableCgroupfs,
		oci.WithMounts(mounts),
		oci.WithHostname(name),
		oci.WithTTY, // docker sets TTY=true
	}
	if len(req.Cmd) > 0 {
		specOpts = append(specOpts, oci.WithProcessArgs(req.Cmd...))
	}
	if req.WorkDir != "" {
		specOpts = append(specOpts, oci.WithProcessCwd(req.WorkDir))
	}

	ctr, err := rc.client.NewContainer(
		ctx,
		name,
		containerd.WithImage(img),
		containerd.WithNewSnapshot(name+"-snapshot", img),
		containerd.WithRuntime(plugins.RuntimeRuncV2, nil),
		containerd.WithNewSpec(specOpts...),
	)
	if err != nil {
		select {
		case <-ctx.Done():
			// we got killed
			return nil, nil
		default:
			return nil, err
		}
	}

	opts := []containerd.NewTaskOpts{}
	spec, err := ctr.Spec(ctx)
	if err != nil {
		return nil, err
	}
	if spec.Linux != nil {
		if len(spec.Linux.UIDMappings) != 0 {
			opts = append(opts, containerd.WithUIDOwner(spec.Linux.UIDMappings[0].HostID))
		}
		if len(spec.Linux.GIDMappings) != 0 {
			opts = append(opts, containerd.WithGIDOwner(spec.Linux.GIDMappings[0].HostID))
		}
	}

	// Ensure cleanup
	defer func() {
		newCtx := context.Background()
		if t, err := ctr.Task(newCtx, nil); err == nil {
			_ = t.Kill(newCtx, syscall.SIGTERM)
			_ = t.Kill(newCtx, syscall.SIGKILL)
			_, _ = t.Delete(newCtx, containerd.WithProcessKill)
		}
		_ = ctr.Delete(newCtx, containerd.WithSnapshotCleanup)
	}()

	task, err := ctr.NewTask(
		ctx,
		// has to be mounted
		cio.LogFile(logFilePath),
		// permission denied: unknown
		//cio.NewCreator(cio.WithStreams(nil, stdout, stderr), cio.WithFIFODir(fifoPath)),
		opts...,
	)
	if err != nil {
		return nil, err
	}
	waitC, err := task.Wait(ctx)
	if err != nil {
		return nil, err
	}
	if err := task.Start(ctx); err != nil {
		return nil, err
	}

	t, err := tail.TailFile(logFilePath, tail.Config{
		Follow: true,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		for line := range t.Lines {
			io.Copy(stdout, strings.NewReader(line.Text))
		}

		if t.Err() != nil {
			return
		}
	}()

	status := <-waitC

	exit := int64(status.ExitCode())
	return &exit, status.Error()
}
