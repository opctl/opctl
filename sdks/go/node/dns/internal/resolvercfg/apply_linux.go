package resolvercfg

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/iguanesolutions/go-systemd/v5/resolved"
)

// Apply applies the split DNS configuration for a single domain and nameserver.
func Apply(ctx context.Context, domain, nsIPAddress, nsPort string) error {
	fmt.Println("Attempting to configure split DNS via systemd-resolved...")
	conn, err := resolved.NewConn()
	if err != nil {
		fmt.Println("Failed to create resolved connection: %w\n", err)
		fmt.Println("Falling back to /etc/resolv.conf configuration...")
		return configureFallbackDNS(domain, nsIPAddress, nsPort)
	}
	defer conn.Close()

	if err := configureSplitDNS(ctx, conn, domain, nsIPAddress, nsPort); err != nil {
		fmt.Println("Failed to configure split DNS via systemd-resolved: %w\n", err)
		fmt.Println("Falling back to /etc/resolv.conf configuration...")
		return configureFallbackDNS(domain, nsIPAddress, nsPort)
	}

	fmt.Println("Split DNS configuration applied successfully via systemd-resolved.")
	return nil
}

// getInterfaces retrieves all network interfaces.
func getInterfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to list network interfaces: %v", err)
	}
	return ifaces, nil
}

// configureSplitDNS configures systemd-resolved with split DNS settings for all interfaces.
func configureSplitDNS(ctx context.Context, conn *resolved.Conn, domain, nsIPAddress, nsPort string) error {
	// Parse the nameserver IP address.
	nsIP := net.ParseIP(nsIPAddress)
	if nsIP == nil {
		return fmt.Errorf("invalid nameserver IP address: %s", nsIPAddress)
	}

	// Determine the address family (IPv4 or IPv6).
	family := 2 // IPv4
	if nsIP.To4() == nil {
		family = 10 // IPv6
	}

	// Create the DNS configuration.
	dns := []resolved.LinkDNS{
		{
			Family:  family,
			Address: nsIP,
		},
	}

	// Create the domain configuration.
	domains := []resolved.LinkDomain{
		{
			Domain:        domain,
			RoutingDomain: strings.HasPrefix(domain, "~"),
		},
	}

	// Get all network interfaces.
	interfaces, err := getInterfaces()
	if err != nil {
		return err
	}

	// Apply the configuration to all interfaces.
	for _, iface := range interfaces {
		ifindex := iface.Index

		// Set DNS servers for the interface.
		if err := conn.SetLinkDNS(ctx, ifindex, dns); err != nil {
			return fmt.Errorf("failed to set DNS for interface %s (ifindex %d): %v", iface.Name, ifindex, err)
		}
		fmt.Printf("Set DNS for interface %s (ifindex %d)\n", iface.Name, ifindex)

		// Set domains for the interface.
		if err := conn.SetLinkDomains(ctx, ifindex, domains); err != nil {
			return fmt.Errorf("failed to set domains for interface %s (ifindex %d): %v", iface.Name, ifindex, err)
		}
		fmt.Printf("Set domains for interface %s (ifindex %d)\n", iface.Name, ifindex)
	}

	return nil
}

// configureFallbackDNS configures DNS by modifying /etc/resolv.conf without backup.
func configureFallbackDNS(domain, nsIPAddress, nsPort string) error {
	const resolvConfPath = "/etc/resolv.conf"

	// Read the existing resolv.conf.
	originalContent, err := os.ReadFile(resolvConfPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", resolvConfPath, err)
	}

	// Prepare new resolv.conf content.
	var newLines []string
	managedMarker := "# Managed by split-dns program"
	managedStart := "# BEGIN MANAGED BLOCK"
	managedEnd := "# END MANAGED BLOCK"

	// Split the original content into lines and remove any existing managed block.
	originalLines := strings.Split(string(originalContent), "\n")
	inManagedBlock := false
	for _, line := range originalLines {
		if strings.Contains(line, managedStart) {
			inManagedBlock = true
			continue
		}
		if strings.Contains(line, managedEnd) {
			inManagedBlock = false
			continue
		}
		if !inManagedBlock && line != "" {
			newLines = append(newLines, line)
		}
	}

	// Add the new managed block.
	newLines = append(newLines, managedMarker, managedStart)

	// Add the nameserver entry (with port if specified).
	nsEntry := fmt.Sprintf("nameserver %s", nsIPAddress)
	if nsPort != "" && nsPort != "53" { // Default DNS port is 53, omit if standard.
		nsEntry = fmt.Sprintf("nameserver %s:%s", nsIPAddress, nsPort)
	}
	newLines = append(newLines, nsEntry)

	// Add the domain as a search domain if not routing-only.
	if !strings.HasPrefix(domain, "~") {
		newLines = append(newLines, fmt.Sprintf("search %s", strings.TrimPrefix(domain, "~")))
	} else {
		fmt.Printf("Warning: Route-only domain %s cannot be configured in /etc/resolv.conf (interface-specific routing not supported)\n", domain)
	}

	newLines = append(newLines, managedEnd)

	// Write the new resolv.conf without backup.
	newContent := strings.Join(newLines, "\n") + "\n"
	if err := os.WriteFile(resolvConfPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %v", resolvConfPath, err)
	}

	fmt.Printf("Configured DNS via %s (fallback mode)\n", resolvConfPath)
	return nil
}
