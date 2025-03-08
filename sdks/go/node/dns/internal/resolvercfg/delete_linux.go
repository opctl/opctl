package resolvercfg

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/iguanesolutions/go-systemd/v5/resolved"
)

// Delete deletes the split DNS configuration without requiring prior state or backups.
func Delete(ctx context.Context) error {
	// Try deleting from systemd-resolved first.
	fmt.Println("Attempting to delete split DNS configuration from systemd-resolved...")
	if err := deleteSplitDNS(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to delete split DNS configuration from systemd-resolved: %v\n", err)
	} else {
		fmt.Println("Split DNS configuration deleted successfully from systemd-resolved (if any).")
	}

	// Try deleting from /etc/resolv.conf.
	fmt.Println("Attempting to delete DNS configuration from /etc/resolv.conf...")
	if err := deleteFallbackDNS(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to delete DNS configuration from /etc/resolv.conf: %v\n", err)
	} else {
		fmt.Println("DNS configuration deleted successfully from /etc/resolv.conf (if any).")
	}

	return nil
}

// deleteSplitDNS deletes the split DNS settings applied via systemd-resolved for all interfaces.
func deleteSplitDNS() error {
	// Create a new D-Bus connection to systemd-resolved.
	conn, err := resolved.NewConn()
	if err != nil {
		fmt.Println("Failed to connect to systemd-resolved; assuming no settings to delete:", err)
		return nil
	}
	defer conn.Close()

	// Use a background context for D-Bus calls.
	ctx := context.Background()

	// Get all network interfaces.
	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to list network interfaces: %v", err)
	}

	// Reset DNS and domains for each interface.
	for _, iface := range interfaces {
		ifindex := iface.Index

		// Reset DNS servers for the interface (set to empty list).
		if err := conn.SetLinkDNS(ctx, ifindex, []resolved.LinkDNS{}); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to reset DNS for interface %s (ifindex %d): %v\n", iface.Name, ifindex, err)
			continue
		}
		fmt.Printf("Reset DNS for interface %s (ifindex %d)\n", iface.Name, ifindex)

		// Reset domains for the interface (set to empty list).
		if err := conn.SetLinkDomains(ctx, ifindex, []resolved.LinkDomain{}); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to reset domains for interface %s (ifindex %d): %v\n", iface.Name, ifindex, err)
			continue
		}
		fmt.Printf("Reset domains for interface %s (ifindex %d)\n", iface.Name, ifindex)
	}

	return nil
}

var etcResolvConfPath = "/etc/resolv.conf"

// deleteFallbackDNS deletes the DNS settings applied via /etc/resolv.conf by removing the managed block.
func deleteFallbackDNS() error {

	// Read the current resolv.conf.
	currentContent, err := os.ReadFile(etcResolvConfPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", etcResolvConfPath, err)
	}

	// Check if there is a managed block to remove.
	managedStart := "# BEGIN MANAGED BLOCK"
	managedEnd := "# END MANAGED BLOCK"
	currentLines := strings.Split(string(currentContent), "\n")
	var newLines []string
	inManagedBlock := false
	hasManagedBlock := false

	for _, line := range currentLines {
		if strings.Contains(line, managedStart) {
			inManagedBlock = true
			hasManagedBlock = true
			continue
		}
		if strings.Contains(line, managedEnd) {
			inManagedBlock = false
			continue
		}
		if !inManagedBlock {
			newLines = append(newLines, line)
		}
	}

	if !hasManagedBlock {
		fmt.Println("No managed block found in /etc/resolv.conf; assuming no settings to delete")
		return nil
	}

	// Write the updated resolv.conf without the managed block.
	newContent := strings.Join(newLines, "\n")
	if strings.TrimSpace(newContent) != "" {
		newContent += "\n"
	}
	if err := os.WriteFile(etcResolvConfPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %v", etcResolvConfPath, err)
	}

	fmt.Printf("Removed managed block from %s\n", etcResolvConfPath)
	return nil
}
