package docker

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/sdks/go/model"
	uuid "github.com/satori/go.uuid"

	"net"
	"os"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/glinton/ping"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const (
	UTUN = "utun"
)

var (
	defaultKey     wgtypes.Key
	listenPort     = 3333
	hostPrivateKey *wgtypes.Key
	vmPeerIp       = "10.33.33.2"
	vmPrivateKey   *wgtypes.Key
	hostPeerIp     = "10.33.33.1"
	hostPeerCIDR   = hostPeerIp + "/32"
)

// ensureNetworkAttached is concurrency safe and is intended to be run every time a container starts in order to self-heal
// in cases such as Docker For mac power saver mode activating. Power saver mode is enabled in docker by default and shuts
// the VM off after 5 min of no containers being run.
func ensureNetworkAttached(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
) error {

	// Routes to container IPs do not exist on docker for mac. This is inconsistent with linux docker
	// and is so important we implement this bandaid.
	//
	// note: This will likely be brittle because It's relying on undocumented docker for mac internals.
	if runtime.GOOS == "darwin" {
		networkInspect, networkInspectErr := dockerClient.NetworkInspect(
			ctx,
			networkName,
			network.InspectOptions{},
		)
		if networkInspectErr != nil {
			return fmt.Errorf("unable to inspect network: %w", networkInspectErr)
		}

		go func() {
			defer func() {
				if panic := recover(); panic != nil {
					// recover from panics
					fmt.Printf("recovered from panic: %s\n%s\n", panic, string(debug.Stack()))
				}
			}()

			err := wgUp(dockerClient, networkInspect)
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	return nil
}

func wgUp(
	dockerClient dockerClientPkg.CommonAPIClient,
	network network.Inspect,
) error {

	if hostPrivateKey == nil {
		pk, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return fmt.Errorf("Failed to generate host private key: %w", err)
		}
		hostPrivateKey = &pk
	}

	if vmPrivateKey == nil {
		pk, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return fmt.Errorf("Failed to generate VM private key: %w", err)
		}
		vmPrivateKey = &pk
	}

	go func() {
		defer func() {
			if panic := recover(); panic != nil {
				// recover from panics
				fmt.Printf("recovered from panic: %s\n%s\n", panic, string(debug.Stack()))
			}
		}()

		ctx := context.Background()
		pingCtx, _ := context.WithTimeout(ctx, time.Second)
		// test if VM pingable
		_, err := ping.IPv4(pingCtx, vmPeerIp)
		// if VM pingable, nothing to do
		if err == nil {
			return
		}

		if err := setupVm(
			ctx,
			dockerClient,
			listenPort,
			hostPeerIp,
			vmPeerIp,
			*hostPrivateKey,
			*vmPrivateKey,
		); err != nil {
			fmt.Println(err.Error())
		}
	}()

	tunDevice, err := createTunIfNotExists(context.Background())
	if err != nil {
		return err
	}

	if tunDevice == nil {
		// utunDevice alread exists; nothing to do...
		return nil
	}

	interfaceName, err := tunDevice.Name()
	if err != nil {
		return err
	}

	fileUAPI, err := ipc.UAPIOpen(interfaceName)
	if err != nil {
		return fmt.Errorf("UAPI listen error: %w", err)
	}

	errs := make(chan error)

	uapi, err := ipc.UAPIListen(interfaceName, fileUAPI)
	if err != nil {
		return fmt.Errorf("Failed to listen on UAPI socket: %w", err)
	}

	// Clean up
	defer uapi.Close()

	wgDevice := device.NewDevice(
		tunDevice,
		conn.NewDefaultBind(),
		device.NewLogger(
			1,
			fmt.Sprintf("(%s) ", interfaceName),
		),
	)

	go func() {
		defer func() {
			if panic := recover(); panic != nil {
				// recover from panics
				fmt.Printf("recovered from panic: %s\n%s\n", panic, string(debug.Stack()))
			}
		}()

		for {
			conn, err := uapi.Accept()
			if err != nil {
				errs <- err
				return
			}
			go wgDevice.IpcHandle(conn)
		}
	}()

	_, wildcardIpNet, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return fmt.Errorf("Failed to parse wildcard CIDR: %w", err)
	}

	_, vmIpNet, err := net.ParseCIDR(vmPeerIp + "/32")
	if err != nil {
		return fmt.Errorf("Failed to parse VM peer CIDR: %w", err)
	}

	wgClient, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("Failed to create new wgctrl client: %w", err)
	}

	defer wgClient.Close()

	peer := wgtypes.PeerConfig{
		PublicKey: vmPrivateKey.PublicKey(),
		AllowedIPs: []net.IPNet{
			*wildcardIpNet,
			*vmIpNet,
		},
	}

	err = wgClient.ConfigureDevice(
		interfaceName,
		wgtypes.Config{
			ListenPort: &listenPort,
			PrivateKey: hostPrivateKey,
			Peers:      []wgtypes.PeerConfig{peer},
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to configure Wireguard device: %w", err)
	}

	for _, config := range network.IPAM.Config {
		if network.Scope == "local" {
			cmd := exec.Command("route", "-q", "-n", "add", "-inet", config.Subnet, "-interface", interfaceName)

			outputBytes, err := cmd.CombinedOutput()

			if err != nil {
				return fmt.Errorf("Failed to add route: %w, %s", err, string(outputBytes))
			}

		}
	}

	// Wait for program to terminate
	<-wgDevice.Wait()

	return nil
}

// createTunIfNotExists returns nil if a tun device already exists; otherwise returns the created device
func createTunIfNotExists(
	ctx context.Context,
) (tun.Device, error) {
	// always attempt to create to avoid races
	tunDevice, err := tun.CreateTUN("utun", device.DefaultMTU)
	if err != nil {
		return nil, fmt.Errorf("Failed to create TUN device: %w", err)
	}

	interfaceName, err := tunDevice.Name()
	if err != nil {
		return nil, fmt.Errorf("Failed to get TUN device name: %w", err)
	}

	cmd := exec.Command("ifconfig", interfaceName, "inet", hostPeerCIDR, vmPeerIp)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Failed to set interface address with ifconfig: %w, %s", err, string(outputBytes))
	}

	tunIndex, err := strconv.Atoi(strings.TrimPrefix(interfaceName, UTUN))
	if err != nil {
		return nil, err
	}

	lowestTunIndex, err := getLowestTunIndex(
		context.Background(),
	)
	if nil != err {
		return nil, err
	}

	if lowestTunIndex < tunIndex {
		// we got raced and lost..
		tunDevice.Close()
		return nil, nil
	}

	return tunDevice, nil
}

// getLowestTunIndex retrieves the lowest index of any existing utun interface on the system
func getLowestTunIndex(
	ctx context.Context,
) (int, error) {

	lowestTunIndex := 1000000

	interfaces, err := net.Interfaces()
	if err != nil {
		return lowestTunIndex, fmt.Errorf("Failed to list interfaces: %w", err)
	}

	for _, i := range interfaces {
		if !strings.HasPrefix(i.Name, UTUN) {
			continue
		}

		iTunIndex, err := strconv.Atoi(strings.TrimPrefix(i.Name, UTUN))
		if err != nil {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			return lowestTunIndex, err
		}

		for _, a := range addrs {
			if a.String() == hostPeerCIDR {
				// we got raced and lost...
				if lowestTunIndex > iTunIndex {
					lowestTunIndex = iTunIndex
				}
			}
		}
	}

	return lowestTunIndex, nil
}

func setupVm(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
	serverPort int,
	hostPeerIp string,
	vmPeerIp string,
	hostPrivateKey wgtypes.Key,
	vmPrivateKey wgtypes.Key,
) error {
	imageRef := "ghcr.io/chipmk/docker-mac-net-connect/setup:v0.1.3"

	err := pullImage(
		ctx,
		&model.ContainerCall{
			Image: &model.ContainerCallImage{
				Ref: &imageRef,
			},
		},
		dockerClient,
		"",
		noOpEventPublisher{},
	)
	if err != nil {
		return err
	}

	containerName := getContainerName(fmt.Sprintf("wireguard-setup-%s", uuid.NewV4().String()))

	defer dockerClient.ContainerRemove(
		context.Background(),
		containerName,
		container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		},
	)

	resp, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageRef,
			Env: []string{
				"SERVER_PORT=" + strconv.Itoa(serverPort),
				"HOST_PEER_IP=" + hostPeerIp,
				"VM_PEER_IP=" + vmPeerIp,
				"HOST_PUBLIC_KEY=" + hostPrivateKey.PublicKey().String(),
				"VM_PRIVATE_KEY=" + vmPrivateKey.String(),
			},
		},
		&container.HostConfig{
			AutoRemove:  true,
			NetworkMode: "host",
			CapAdd: []string{
				"NET_ADMIN",
			},
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Run container to completion
	err = dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	reader, err := dockerClient.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return fmt.Errorf("failed to get logs for container %s: %w", resp.ID, err)
	}

	defer reader.Close()

	_, err = stdcopy.StdCopy(io.Discard, os.Stderr, reader)
	if err != nil {
		return err
	}

	return nil
}
