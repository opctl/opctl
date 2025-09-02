package embedded

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	"golang.org/x/sync/singleflight"
)

//go:embed template.yaml
var embeddedTemplateYAML []byte

// singleFlightGroup is used to ensure resolves don't race across invocations
var resolveSingleFlightGroup singleflight.Group

func New(
	ctx context.Context,
	dataDir string,
) (containerruntime.ContainerRuntime, error) {

	return _containerRuntime{
		limaPath: limaBinPath(dataDir),
	}, nil
}

type _containerRuntime struct {
	limaPath string
}

func (cr _containerRuntime) Delete(
	ctx context.Context,
) error {
	cr.Kill(ctx)

	cmd := exec.CommandContext(ctx, cr.limaPath, "delete", "default", "--tty=false")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: 502, Gid: 20},
	}
	cmd.Env = cr.getEnv()
	return cmd.Run()
}

func (cr _containerRuntime) DeleteContainerIfExists(
	ctx context.Context,
	containerID string,
) error {
	if _, isRunning := cr.getInstanceStatus(ctx); !isRunning {
		return nil
	}

	dockerCR, err := cr.getDockerContainerRuntime(ctx)
	if err != nil {
		return err
	}

	return dockerCR.DeleteContainerIfExists(ctx, containerID)
}

func (cr _containerRuntime) Kill(
	ctx context.Context,
) error {
	if _, isRunning := cr.getInstanceStatus(ctx); !isRunning {
		return nil
	}

	dockerCR, err := cr.getDockerContainerRuntime(ctx)
	if err != nil {
		return err
	}

	if err := dockerCR.Kill(ctx); err != nil {
		return err
	}

	// Stop the lima instance
	cmd := exec.CommandContext(ctx, cr.limaPath, "stop", "default", "-f", "--tty=false")
	// limactl doesn't allow running stop as root
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: 502, Gid: 20},
	}
	cmd.Env = cr.getEnv()
	return cmd.Run()
}

func (cr _containerRuntime) getEnv() []string {

	n := append(
		os.Environ(),
		"LIMA_HOME="+filepath.Dir(cr.limaPath),
		"PATH="+filepath.Dir(cr.limaPath)+":"+os.Getenv("PATH"),
	)

	return n
}

// RunContainer creates, starts, and waits on a container. ExitCode &/Or an error will be returned
func (cr _containerRuntime) RunContainer(
	ctx context.Context,
	req *model.ContainerCall,
	// @TODO: get rid of in combination with eventPublisher
	rootCallID string,
	// @TODO: get rid of this; just use stdout/stderr
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	dockerCR, err := cr.getDockerContainerRuntime(ctx)
	if err != nil {
		return nil, err
	}

	return dockerCR.RunContainer(ctx, req, rootCallID, eventPublisher, stdout, stderr)
}

func (cr _containerRuntime) getDockerContainerRuntime(
	ctx context.Context,
) (containerruntime.ContainerRuntime, error) {

	// attempt to resolve within singleFlight.Group to ensure concurrent resolves don't race
	crt, err, _ := resolveSingleFlightGroup.Do(
		"",
		func() (interface{}, error) {
			err := inflateLimaBinaries(cr.limaPath)
			if err != nil {
				return nil, fmt.Errorf("failed to ensure lima binaries: %w", err)
			}

			// Ensure the embedded template is available at the expected location
			templatePath := filepath.Join(filepath.Dir(cr.limaPath), "default.yaml")
			if err := cr.ensureTemplate(templatePath); err != nil {
				return nil, fmt.Errorf("failed to ensure template: %w", err)
			}

			// Check if the instance already exists and is running
			exists, _ := cr.getInstanceStatus(ctx)

			args := []string{
				"start",
				"--tty=false",
				"--mount-inotify",
				"--mount-type=virtiofs",
				"--mount-writable",
				"--rosetta",
			}

			if !exists {
				args = append(args, templatePath)
			}

			// Start the VM (this will create it if it doesn't exist, or start it if it exists but is stopped)
			cmd := exec.CommandContext(ctx, cr.limaPath, args...)
			// limactl doesn't allow running start as root
			cmd.SysProcAttr = &syscall.SysProcAttr{
				Credential: &syscall.Credential{Uid: 502, Gid: 20},
			}
			cmd.Env = cr.getEnv()

			if outBytes, err := cmd.CombinedOutput(); err != nil {
				return nil, fmt.Errorf("failed to start lima instance: %w, %s", err, string(outBytes))
			}

			return docker.New(
				ctx,
				fmt.Sprintf("unix://%s",
					filepath.Join(
						filepath.Dir(cr.limaPath),
						"/default/sock/docker.sock",
					),
				),
			)

		},
	)
	if err != nil {
		return nil, err
	}

	return crt.(containerruntime.ContainerRuntime), nil
}

// getInstanceStatus checks if the lima instance exists and if it's running
func (cr _containerRuntime) getInstanceStatus(ctx context.Context) (exists bool, running bool) {
	cmd := exec.CommandContext(ctx, cr.limaPath, "list", "--json", "default")
	// limactl doesn't allow running list as root
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: 502, Gid: 20},
	}
	cmd.Env = cr.getEnv()
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Warning: failed to get instance status: %v\n", err)
		return false, false
	}

	// Parse the JSON output to check if default instance exists and is running
	var instance map[string]interface{}
	if err := json.Unmarshal(output, &instance); err != nil {
		return false, false
	}

	if name, ok := instance["name"].(string); ok && name == "default" {
		exists = true
		if status, ok := instance["status"].(string); ok && status == "Running" {
			running = true
		}
	}

	return exists, running
}

// ensureTemplate ensures the embedded template is available at the specified path
func (cr _containerRuntime) ensureTemplate(templatePath string) error {
	// Check if template already exists and matches the embedded version
	if existingData, err := os.ReadFile(templatePath); err == nil {
		// If the existing file matches the embedded template, no need to write
		if string(existingData) == string(embeddedTemplateYAML) {
			return nil
		}
	}

	if err := unsudo.CreateFile(templatePath, embeddedTemplateYAML); err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	return nil
}
