package docker

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func ensureDnsServerRegistered() error {
	service, err := getNetworkService()
	if err != nil {
		return err
	}
	fmt.Println(service)

	dnsPayload, err := getDNSForPrimaryService(service)
	if err != nil {
		return err
	}
	fmt.Println(dnsPayload.DomainName, dnsPayload.ServerAddresses)

	dnsPayload.ServerAddresses = ensureLocalhostFirst(dnsPayload.ServerAddresses)
	return writeDNSForPrimaryService(service, dnsPayload)
}

func ensureLocalhostFirst(s []string) []string {
	r := []string{"127.0.0.1"}
	for _, h := range s {
		if h != "127.0.0.1" {
			r = append(r, h)
		}
	}
	return r
}

func execSCutilScript(ctx context.Context, script []string) ([]string, error) {
	cmd := exec.CommandContext(ctx, "scutil")

	stdIn, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	for _, line := range script {
		cmdLine := []byte(strings.TrimSpace(line) + "\n")
		n, err := stdIn.Write([]byte(strings.TrimSpace(line) + "\n"))
		if n != len(cmdLine) {
			return nil, fmt.Errorf("No all bytes written to scutil")
		}
		if err != nil {
			return nil, err
		}
	}

	buff := bufio.NewScanner(stdOut)
	var allText []string

	for buff.Scan() {
		allText = append(allText, buff.Text())
	}

	return allText, nil
}

func getNetworkService() (string, error) {
	script := []string {
		"open",
		"get State:/Network/Global/IPv4",
		"d.show",
		"close",
		"quit",
	}

	ipv4Settings, err := execSCutilScript(context.Background(), script)
	if err != nil {
		return "", err
	}

	for _, l := range ipv4Settings {
		parts := strings.Split(l, ":")
		if strings.TrimSpace(parts[0]) == "PrimaryService" {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("Didn't find primary service")
}

type DNSPayload struct {
	DomainName string
	ServerAddresses []string
}

func getDNSForPrimaryService(service string) (*DNSPayload, error) {
	script := []string {
		"open",
		fmt.Sprintf("get State:/Network/Service/%s/DNS", service),
		"d.show",
		"close",
		"quit",
	}

	dnsSettings, err := execSCutilScript(context.Background(), script)
	if err != nil {
		return nil, err
	}

	addresses := []string{}
	domain := ""
	for i := 1; i < len(dnsSettings); i++ {
		parts := strings.Split(dnsSettings[i], ":")
		if strings.TrimSpace(parts[0]) == "ServerAddresses" {
			i += 1
			for strings.TrimSpace(dnsSettings[i]) != "}" {
				dnsPart := strings.Split(dnsSettings[i], ":")
				addresses = append(addresses, strings.TrimSpace(dnsPart[1]))
				i ++
			}
		}
		if strings.TrimSpace(parts[0]) == "DomainName" {
			domain = strings.TrimSpace(parts[1])
		}
	}

	return &DNSPayload{
		DomainName:      domain,
		ServerAddresses: addresses,
	}, nil
}

func writeDNSForPrimaryService(service string, dns *DNSPayload) error {

	addServerAddresses := fmt.Sprintf("d.add ServerAddresses %s", strings.Join(dns.ServerAddresses, " "))
	addDomain := ""
	if dns.DomainName != "" {
		addDomain = fmt.Sprintf("d.add DomainName %s", dns.DomainName)
	}

	script := []string {
		// Set State of resolver
		fmt.Sprintf("get %s:/Network/Service/%s/DNS", "State", service),
		"d.remove SearchDomains",
		"d.remove ServerAddresses",
		addServerAddresses,
		addDomain,
		fmt.Sprintf("set %s:/Network/Service/%s/DNS", "State", service),

		// Set Setup State of Resolver (which is the current state ??? )
		fmt.Sprintf("get %s:/Network/Service/%s/DNS", "Setup", service),
		"d.remove SearchDomains",
		"d.remove ServerAddresses",
		addServerAddresses,
		addDomain,
		fmt.Sprintf("set %s:/Network/Service/%s/DNS", "Setup", service),

		// Done
		"exit",
	}
	_, err := execSCutilScript(context.Background(), script)
	if err != nil {
		return err
	}
	return nil
}