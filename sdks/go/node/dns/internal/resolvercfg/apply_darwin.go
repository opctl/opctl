package resolvercfg

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

var scutilKeyPrefix = "State:/Network/Service/opctl_"

// validateScutilInput ensures a value is safe to interpolate into scutil stdin commands.
// scutil reads commands line-by-line, so newlines would allow command injection.
func validateScutilInput(name, value string) error {
	if strings.ContainsAny(value, "\n\r") {
		return fmt.Errorf("invalid %s %q: must not contain newlines", name, value)
	}
	return nil
}

// Apply to the current system
func Apply(
	ctx context.Context,
	domain,
	nsIPAddress,
	nsPort string,
) error {
	for name, value := range map[string]string{
		"domain":    domain,
		"address":   nsIPAddress,
		"port":      nsPort,
	} {
		if err := validateScutilInput(name, value); err != nil {
			return err
		}
	}

	key := scutilKeyPrefix + domain + "/DNS"

	var sb strings.Builder
	sb.WriteString("d.init\n")
	sb.WriteString(fmt.Sprintf("d.add ServerAddresses * %s\n", nsIPAddress))
	if nsPort != "53" {
		sb.WriteString(fmt.Sprintf("d.add ServerPort # %s\n", nsPort))
	}
	sb.WriteString(fmt.Sprintf("d.add SupplementalMatchDomains * %s\n", domain))
	sb.WriteString(fmt.Sprintf("set %s\n", key))

	cmd := exec.CommandContext(
		ctx,
		"scutil",
	)
	cmd.Stdin = strings.NewReader(sb.String())

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply resolver config via scutil: %w, %s", err, string(outputBytes))
	}

	// scutil exits 0 even on permission errors, so check output
	if strings.Contains(string(outputBytes), "Permission denied") {
		return fmt.Errorf("failed to apply resolver config via scutil: permission denied (run as root)")
	}

	return clearCaches(ctx)
}
