package docker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"

	"io"
	"net"
	"os"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
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

func ensureNetworkAttached(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
) error {

	// Routes to container IPs do not exist on docker for mac. This is inconsistent with linux docker
	// and is so important we implement this bandaid.
	//
	// note: This will likely be brittle because It's relying on undocumented docker for mac internals.
	if runtime.GOOS == "darwin" {
		networkResource, networkInspectErr := dockerClient.NetworkInspect(
			ctx,
			networkName,
			types.NetworkInspectOptions{},
		)
		if networkInspectErr != nil {
			return fmt.Errorf("unable to inspect network: %w", networkInspectErr)
		}

		go func() {
			err := wgUp(dockerClient, networkResource)
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	return nil
}

func wgUp(
	dockerClient dockerClientPkg.CommonAPIClient,
	network types.NetworkResource,
) error {
	hostPeerIp := "10.33.33.1"
	hostPeerCIDR := hostPeerIp + "/32"
	vmPeerIp := "10.33.33.2"

	// always attempt to create to avoid races
	tunDevice, err := tun.CreateTUN("utun", device.DefaultMTU)
	if err != nil {
		return fmt.Errorf("Failed to create TUN device: %w", err)
	}

	interfaceName, err := tunDevice.Name()
	if err != nil {
		return fmt.Errorf("Failed to get TUN device name: %w", err)
	}

	cmd := exec.Command("ifconfig", interfaceName, "inet", hostPeerCIDR, vmPeerIp)

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to set interface address with ifconfig: %w, %s", err, string(outputBytes))
	}

	tunIndex, err := strconv.Atoi(strings.TrimPrefix(interfaceName, UTUN))
	if err != nil {
		return err
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("Failed to list interfaces: %w", err)
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
			return err
		}

		for _, a := range addrs {
			if a.String() == hostPeerCIDR && tunIndex > iTunIndex {
				// we got raced and lost...
				return tunDevice.Close()
			}
		}
	}

	fileUAPI, err := ipc.UAPIOpen(interfaceName)
	if err != nil {
		return fmt.Errorf("UAPI listen error: %w", err)
	}

	wgDevice := device.NewDevice(
		tunDevice, conn.NewDefaultBind(),
		device.NewLogger(
			2,
			fmt.Sprintf("(%s) ", interfaceName),
		),
	)

	errs := make(chan error)

	//"sudo networksetup -setdnsservers (Network Service) (DNS IP)"

	uapi, err := ipc.UAPIListen(interfaceName, fileUAPI)
	if err != nil {
		return fmt.Errorf("Failed to listen on UAPI socket: %w", err)
	}

	go func() {
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

	hostPrivateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("Failed to generate host private key: %w", err)
	}

	vmPrivateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("Failed to generate VM private key: %w", err)
	}

	// Wireguard configuration
	listenPort := 3333

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

	err = wgClient.ConfigureDevice(interfaceName, wgtypes.Config{
		ListenPort: &listenPort,
		PrivateKey: &hostPrivateKey,
		Peers:      []wgtypes.PeerConfig{peer},
	})
	if err != nil {
		return fmt.Errorf("Failed to configure Wireguard device: %w", err)
	}

	ctx := context.Background()

	err = setupVm(ctx, dockerClient, listenPort, hostPeerIp, vmPeerIp, hostPrivateKey, vmPrivateKey)
	if err != nil {
		return fmt.Errorf("Failed to setup VM: %w", err)
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

	// Clean up
	uapi.Close()
	wgDevice.Close()

	return nil
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
	imageName := "ghcr.io/chipmk/docker-mac-net-connect/setup:v0.1.3"

	_, _, err := dockerClient.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		fmt.Printf("Image doesn't exist locally. Pulling...\n")

		pullStream, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull setup image: %w", err)
		}

		io.Copy(os.Stdout, pullStream)
	}

	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env: []string{
			"SERVER_PORT=" + strconv.Itoa(serverPort),
			"HOST_PEER_IP=" + hostPeerIp,
			"VM_PEER_IP=" + vmPeerIp,
			"HOST_PUBLIC_KEY=" + hostPrivateKey.PublicKey().String(),
			"VM_PRIVATE_KEY=" + vmPrivateKey.String(),
		},
	}, &container.HostConfig{
		AutoRemove:  true,
		NetworkMode: "host",
		CapAdd:      []string{"NET_ADMIN"},
	}, nil, nil, "wireguard-setup")
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Run container to completion
	err = dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	func() error {
		reader, err := dockerClient.ContainerLogs(ctx, resp.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
		if err != nil {
			return fmt.Errorf("failed to get logs for container %s: %w", resp.ID, err)
		}

		defer reader.Close()

		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
		if err != nil {
			return err
		}

		return nil
	}()

	return nil
}
