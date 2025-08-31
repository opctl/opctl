//go:build darwin
// +build darwin

package applevf

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/lima-vm/lima/v2/pkg/imgutil/proxyimgutil"
	"github.com/Code-Hex/vz/v3"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

func New(
	cacheDir string,
) (containerruntime.ContainerRuntime, error) {

	// Prepare VM configuration
	cr := containerRuntime{cacheDir: cacheDir}
	vm, err := cr.prepareVMConfig(
		cacheDir,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare VM config: %w", err)
	}

	// Start a goroutine to monitor VM state changes (silently)
	go func() {
		for newState := range vm.StateChangedNotify() {
			// Silent monitoring - only log if needed for debugging
			_ = newState
		}
	}()

	cr.vm = vm

	return &cr, nil
}

type containerRuntime struct {
	vm       *vz.VirtualMachine
	cacheDir string
}

func (cr *containerRuntime) prepareVMConfig(
	cacheDir string,
) (*vz.VirtualMachine, error) {

	kernelCommandLineArguments := []string{
		// Use the first virtio console device as system console.
		"console=hvc0",
		// Ubuntu cloud images typically have the root filesystem on a partition
		"root=/dev/vda1",
		// Specify the root filesystem type
		"rootfstype=ext4",
		// Mount options for better compatibility
		"rootflags=rw,errors=remount-ro",
		// Additional kernel parameters for cloud images
		"ro",
		"init=/sbin/init",
	}
	
	// Download or use cached kernel and initrd
	kernelPath := filepath.Join(cacheDir, "vmlinuz")
	initrdPath := filepath.Join(cacheDir, "initrd")

	err := cr.downloadFile(
		"https://cloud-images.ubuntu.com/releases/plucky/release-20250701/unpacked/ubuntu-25.04-server-cloudimg-arm64-vmlinuz-generic",
		kernelPath,
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to download and decompress kernel: %w", err)
	}

	if err := cr.downloadFile(
		"https://cloud-images.ubuntu.com/releases/plucky/release-20250701/unpacked/ubuntu-25.04-server-cloudimg-arm64-initrd-generic",
		initrdPath,
		false,
	); err != nil {
		return nil, fmt.Errorf("failed to download initrd: %w", err)
	}

	bootLoader, err := vz.NewLinuxBootLoader(
		kernelPath,
		vz.WithCommandLine(strings.Join(kernelCommandLineArguments, " ")),
		vz.WithInitrd(initrdPath),
	)
	if err != nil {
		return nil, fmt.Errorf("bootloader creation failed: %s", err)
	}

	config, err := vz.NewVirtualMachineConfiguration(
		bootLoader,
		1,
		2*1024*1024*1024,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual machine configuration: %s", err)
	}

	// console
	serialPortAttachment, err := vz.NewFileHandleSerialPortAttachment(os.Stdin, os.Stdout)
	if err != nil {
		return nil, fmt.Errorf("serial port attachment creation failed: %s", err)
	}
	consoleConfig, err := vz.NewVirtioConsoleDeviceSerialPortConfiguration(serialPortAttachment)
	if err != nil {
		return nil, fmt.Errorf("failed to create serial configuration: %s", err)
	}
	config.SetSerialPortsVirtualMachineConfiguration([]*vz.VirtioConsoleDeviceSerialPortConfiguration{
		consoleConfig,
	})

	// network
	natAttachment, err := vz.NewNATNetworkDeviceAttachment()
	if err != nil {
		return nil, fmt.Errorf("nat network device creation failed: %s", err)
	}
	networkConfig, err := vz.NewVirtioNetworkDeviceConfiguration(natAttachment)
	if err != nil {
		return nil, fmt.Errorf("creation of the networking configuration failed: %s", err)
	}
	config.SetNetworkDevicesVirtualMachineConfiguration([]*vz.VirtioNetworkDeviceConfiguration{
		networkConfig,
	})
	mac, err := vz.NewRandomLocallyAdministeredMACAddress()
	if err != nil {
		return nil, fmt.Errorf("random MAC address creation failed: %s", err)
	}
	networkConfig.SetMACAddress(mac)

	// entropy
	entropyConfig, err := vz.NewVirtioEntropyDeviceConfiguration()
	if err != nil {
		return nil, fmt.Errorf("entropy device creation failed: %s", err)
	}
	config.SetEntropyDevicesVirtualMachineConfiguration([]*vz.VirtioEntropyDeviceConfiguration{
		entropyConfig,
	})

	diskImageQcowPath := filepath.Join(cacheDir, "disk.img")
	diskImageRawPath := filepath.Join(cacheDir, "disk.raw")

	err = cr.downloadFile(
		"https://cloud-images.ubuntu.com/releases/plucky/release-20250701/ubuntu-25.04-server-cloudimg-arm64.img",
		diskImageQcowPath,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to download disk image: %w", err)
	}
	diskUtil := proxyimgutil.NewDiskUtil(
		context.Background(),
	)

	diskSize, _ := units.RAMInBytes("100GiB")
	if diskSize == 0 {
		return nil, fmt.Errorf("invalid disk size")
	}

	err = diskUtil.ConvertToRaw(
		context.Background(),
		diskImageQcowPath,
		diskImageRawPath,
		&diskSize,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create disk image: %s", err)
	}

	diskImageAttachment, err := vz.NewDiskImageStorageDeviceAttachment(
		diskImageRawPath,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create disk image attachment: %w", err)
	}

	storageDeviceConfig, err := vz.NewVirtioBlockDeviceConfiguration(diskImageAttachment)
	if err != nil {
		return nil, fmt.Errorf("block device creation failed: %s", err)
	}

	config.SetStorageDevicesVirtualMachineConfiguration([]vz.StorageDeviceConfiguration{
		storageDeviceConfig,
	})

	// traditional memory balloon device which allows for managing guest memory. (optional)
	memoryBalloonDevice, err := vz.NewVirtioTraditionalMemoryBalloonDeviceConfiguration()
	if err != nil {
		return nil, fmt.Errorf("balloon device creation failed: %s", err)
	}
	config.SetMemoryBalloonDevicesVirtualMachineConfiguration([]vz.MemoryBalloonDeviceConfiguration{
		memoryBalloonDevice,
	})

	// socket device (optional)
	vsockDevice, err := vz.NewVirtioSocketDeviceConfiguration()
	if err != nil {
		return nil, fmt.Errorf("virtio-vsock device creation failed: %s", err)
	}
	config.SetSocketDevicesVirtualMachineConfiguration([]vz.SocketDeviceConfiguration{
		vsockDevice,
	})
	fmt.Printf("Validating VM configuration...\n")
	validated, err := config.Validate()
	if !validated || err != nil {
		fmt.Printf("VM configuration validation failed: %v\n", err)
		return nil, fmt.Errorf("virtual machine configuration validation failed: %s", err)
	}
	fmt.Printf("VM configuration validated successfully\n")

	fmt.Printf("Creating virtual machine...\n")
	vm, err := vz.NewVirtualMachine(config)
	if err != nil {
		fmt.Printf("VM creation failed: %v\n", err)
		return nil, fmt.Errorf("failed to create virtual machine: %w", err)
	}
	fmt.Printf("Virtual machine created successfully\n")

	return vm, nil
}

func (cr *containerRuntime) downloadFile(url, destPath string, decompressGzip bool) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", destPath, err)
	}

	var err error
	if _, err = os.Stat(destPath); os.IsNotExist(err) {
		// Download file using Go's native HTTP client
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP %d", resp.StatusCode)
		}

		// Create the output file
		out, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer out.Close()

		var reader io.Reader = resp.Body
		if decompressGzip {
			gzipReader, err := gzip.NewReader(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to create gzip reader: %w", err)
			}
			defer gzipReader.Close()
			reader = gzipReader
		}

		// Copy the response body to the file
		_, err = io.Copy(out, reader)
		return err
	}
	return err
}

func (cr *containerRuntime) Delete(ctx context.Context) error {
	if cr.vm.State() == vz.VirtualMachineStateRunning {
		return cr.vm.Stop()
	}
	return nil
}

func (cr *containerRuntime) DeleteContainerIfExists(ctx context.Context, containerID string) error {
	// vfkit doesn't have persistent containers like Docker
	// Each run is a new VM instance
	return nil
}

func (cr *containerRuntime) Kill(ctx context.Context) error {
	// For vfkit, killing is handled by the VM lifecycle
	return nil
}

func (cr *containerRuntime) RunContainer(
	ctx context.Context,
	req *model.ContainerCall,
	rootCallID string,
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {

	if cr.vm.State() != vz.VirtualMachineStateRunning {
		if err := cr.vm.Start(); err != nil {
			return nil, fmt.Errorf("start virtual machine is failed: %s", err)
		}

		// Give the guest VM time to boot and produce output
		time.Sleep(5 * time.Second)
	}

	exitCode := int64(0)
	return &exitCode, nil
}
