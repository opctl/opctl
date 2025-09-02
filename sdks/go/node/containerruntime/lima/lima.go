package lima

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

func New(
	ctx context.Context,
	dataDir string,
) (containerruntime.ContainerRuntime, error) {
	limaPath, err := ensureLimaBinary(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure lima binary: %w", err)
	}

	cr := _containerRuntime{
		limaPath: limaPath,
	}

	if _, err := cr.getDockerContainerRuntime(ctx); err != nil {
		return nil, err
	}

	return cr, nil
}

type _containerRuntime struct {
	limaPath string
}

func (cr _containerRuntime) Delete(
	ctx context.Context,
) error {
	// Stop and delete all lima instances
	cmd := exec.CommandContext(ctx, cr.limaPath, "stop", "--all")
	cmd.Env = cr.getEnv()
	if err := cmd.Run(); err != nil {
		// Ignore errors if no instances are running
		fmt.Printf("Warning: failed to stop lima instances: %v\n", err)
	}

	cmd = exec.CommandContext(ctx, cr.limaPath, "delete", "--all")
	return cmd.Run()
}

func (cr _containerRuntime) DeleteContainerIfExists(
	ctx context.Context,
	containerID string,
) error {
	if !cr.isVMRunning(ctx) {
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
	if !cr.isVMRunning(ctx) {
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
	cmd := exec.CommandContext(ctx, cr.limaPath, "stop", "default")
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
	if !cr.isVMRunning(ctx) {
		// cmd := exec.CommandContext(ctx, "limactl", "create", "--name=default", "template://docker")
		// cmd.Env = cr.getEnv()
		// if outBytes, err := cmd.CombinedOutput(); err != nil {
		// 	return nil, fmt.Errorf("failed to create lima instance: %w, %s", err, string(outBytes))
		// }

		// cmd = exec.CommandContext(ctx, "limactl", "start")
		// cmd.Env = cr.getEnv()
		// if outBytes, err := cmd.CombinedOutput(); err != nil {
		// 	return nil, fmt.Errorf("failed to start lima instance: %w, %s", err, string(outBytes))
		// }

		// Wait for docker to be ready
		// cmd = exec.CommandContext(ctx, "limactl", "shell", "default", "sudo", "service", "docker", "start")
		// cmd.Env = cr.getEnv()
		// if outBytes, err := cmd.CombinedOutput(); err != nil {
		// 	return nil, fmt.Errorf("failed to start docker service: %w, %s", err, string(outBytes))
		// }

		// Set proper permissions on docker socket
		// cmd = exec.CommandContext(ctx, "limactl", "shell", "default", "sudo", "chmod", "0666", "/var/run/docker.sock")
		// cmd.Env = cr.getEnv()
		// if outBytes, err := cmd.CombinedOutput(); err != nil {
		// 	return nil, fmt.Errorf("failed to set docker socket permissions: %w, %s", err, string(outBytes))
		// }
	}

	return docker.New(
		ctx,
		fmt.Sprintf("unix://%s",
			filepath.Join(
				filepath.Dir(cr.limaPath),
				"/docker/sock/docker.sock",
			),
		),
	)
}

// isVMRunning checks if the lima VM is running
func (cr _containerRuntime) isVMRunning(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "limactl", "list", "--json")
	cmd.Env = cr.getEnv()
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Parse the JSON output to check if default instance is running
	var instances []map[string]interface{}
	if err := json.Unmarshal(output, &instances); err != nil {
		return false
	}

	for _, instance := range instances {
		if name, ok := instance["name"].(string); ok && name == "default" {
			if status, ok := instance["status"].(string); ok && status == "Running" {
				return true
			}
		}
	}

	return false
}

// ensureLimaBinary downloads and caches the latest lima binary
func ensureLimaBinary(dataDir string) (string, error) {
	binaryPath := filepath.Join(dataDir, "vms", "limactl")

	// Check if binary already exists and is executable
	if info, err := os.Stat(binaryPath); err == nil && info.Mode().Perm()&0111 != 0 {
		return binaryPath, nil
	}

	// Create cache directory
	if err := os.MkdirAll(
		filepath.Dir(binaryPath),
		0755,
	); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Find the appropriate asset for current platform
	binaryURL, err := findBinaryURL()
	if err != nil {
		return "", err
	}

	// Download the binary
	if err := downloadFile(binaryURL, binaryPath); err != nil {
		return "", fmt.Errorf("failed to download lima binary: %w", err)
	}

	// Make binary executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return "", fmt.Errorf("failed to make binary executable: %w", err)
	}

	return binaryPath, nil
}

// findAssetForPlatform finds the appropriate asset for the current platform
func findBinaryURL() (string, error) {
	goarch := runtime.GOARCH

	// Map Go architecture to common names
	archMap := map[string]string{
		"amd64": "https://github.com/lima-vm/lima/releases/download/v1.2.1/lima-1.2.1-Darwin-x86_64.tar.gz",
		"arm64": "https://github.com/lima-vm/lima/releases/download/v1.2.1/lima-1.2.1-Darwin-arm64.tar.gz",
	}

	if mapped, ok := archMap[goarch]; ok {
		return mapped, nil
	}

	return "", fmt.Errorf("unsupported architecture: %s", goarch)
}

// downloadFile downloads a file from URL and extracts it if it's a tar.gz
func downloadFile(url, binaryPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Check if this is a tar.gz file by looking at the URL
	if strings.HasSuffix(url, ".tar.gz") || strings.Contains(url, ".tar.gz") {
		return extractTarGz(resp.Body, binaryPath)
	}

	// For direct binary downloads
	out, err := os.Create(binaryPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractTarGz extracts a tar.gz archive and finds the lima binary
func extractTarGz(reader io.Reader, binaryPath string) error {
	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Look for the lima binary (usually named "lima" or "limactl")
		if header.Typeflag == tar.TypeReg && (header.Name == "lima" || header.Name == "limactl" || strings.HasSuffix(header.Name, "/lima") || strings.HasSuffix(header.Name, "/limactl")) {
			out, err := os.Create(binaryPath)
			if err != nil {
				return fmt.Errorf("failed to create binary file: %w", err)
			}
			defer out.Close()

			if _, err := io.Copy(out, tr); err != nil {
				return fmt.Errorf("failed to extract binary: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("lima binary not found in archive")
}
