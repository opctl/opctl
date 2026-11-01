package node

import (
	"archive/tar"
	"bytes"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
)

type bootstrapAuth struct {
	archive []byte
	private string
}
type bootstrapAuthScope struct {
	host, path string
	creds      model.Creds
	explicit   bool
}

func authScope(resource string) (bootstrapAuthScope, bool) {
	if strings.ContainsAny(resource, "#?@*[]\\\x00\r\n\t ") {
		return bootstrapAuthScope{}, false
	}
	if !strings.Contains(resource, "://") {
		resource = "https://" + resource
	}
	parsed, err := url.Parse(resource)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Hostname() == "" {
		return bootstrapAuthScope{}, false
	}
	host := strings.ToLower(parsed.Host)
	for _, r := range parsed.Hostname() {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '.' || r == '-') {
			return bootstrapAuthScope{}, false
		}
	}
	if strings.Contains(host, ":") && parsed.Port() == "" {
		return bootstrapAuthScope{}, false
	}
	for _, r := range parsed.Path {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("/._-", r)) {
			return bootstrapAuthScope{}, false
		}
	}
	if !strings.Contains(parsed.Hostname(), ".") && parsed.Hostname() != "localhost" {
		return bootstrapAuthScope{}, false
	}
	if parsed.Port() != "" {
		for _, r := range parsed.Port() {
			if r < '0' || r > '9' {
				return bootstrapAuthScope{}, false
			}
		}
	}
	scopePath := strings.Trim(strings.ToLower(parsed.Path), "/")
	if strings.ContainsAny(scopePath, ":%") || (scopePath != "" && path.Clean(scopePath) != scopePath) || scopePath == ".." || strings.HasPrefix(scopePath, "../") {
		return bootstrapAuthScope{}, false
	}
	return bootstrapAuthScope{host: host, path: strings.TrimSuffix(scopePath, ".git")}, true
}
func shellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'" }
func prepareBootstrapAuth(auths []model.Auth, resource string, explicit *model.Creds) (bootstrapAuth, error) {
	scopes := []bootstrapAuthScope{}
	for _, auth := range auths {
		if scope, ok := authScope(auth.Resources); ok {
			scope.creds = auth.Creds
			scopes = append(scopes, scope)
		}
	}
	if explicit != nil {
		if scope, ok := authScope(resource); ok && scope.path != "" {
			scope.creds, scope.explicit = *explicit, true
			scopes = append(scopes, scope)
		}
	}
	if len(scopes) == 0 {
		return bootstrapAuth{}, nil
	}
	for _, scope := range scopes {
		if strings.ContainsAny(scope.creds.Username+scope.creds.Password, "\x00\r\n") {
			return bootstrapAuth{}, fmt.Errorf("Go bootstrap HTTPS credentials cannot contain line breaks or NUL")
		}
	}
	sort.SliceStable(scopes, func(i, j int) bool {
		a, b := scopes[i], scopes[j]
		if a.explicit != b.explicit {
			return a.explicit
		}
		if len(a.path) != len(b.path) {
			return len(a.path) > len(b.path)
		}
		return a.host+"/"+a.path < b.host+"/"+b.path
	})
	helper := `#!/bin/sh
set -eu
[ "${1:-get}" = get ] || exit 0
protocol= host= path=
while IFS='=' read -r key value; do
  [ -n "$key" ] || break
  case "$key" in protocol) protocol=$value;; host) host=$value;; path) path=$value;; esac
done
[ "$protocol" = https ] || exit 0
host=$(printf '%s' "$host" | tr '[:upper:]' '[:lower:]')
path=$(printf '%s' "$path" | tr '[:upper:]' '[:lower:]')
path=${path#/}
path=${path%.git}
case "/$path/" in */../*|*/./*|*//*) exit 0;; esac
case "$path" in *[!a-z0-9/._-]*) exit 0;; esac
`
	prefixes := map[string]bool{}
	netrcHosts := map[string]bool{}
	var netrc strings.Builder
	for _, scope := range scopes {
		prefix := scope.host
		if scope.path != "" {
			prefix += "/" + scope.path
		}
		prefixes[prefix] = true
		helper += "if [ \"$host\" = " + shellQuote(scope.host) + " ]; then\n"
		if scope.path != "" {
			helper += "case \"$path\" in " + shellQuote(scope.path) + "|" + shellQuote(scope.path) + "/*)\n"
		}
		helper += "printf '%s\\n' " + shellQuote("username="+scope.creds.Username) + " " + shellQuote("password="+scope.creds.Password) + "\nexit 0\n"
		if scope.path != "" {
			helper += ";; esac\n"
		}
		helper += "fi\n"
		// Older Go netrc readers cannot represent whitespace or quoted/escaped tokens.
		eligible := scope.path == "" && !strings.Contains(scope.host, ":") && scope.creds.Username != "" && scope.creds.Password != "" && !strings.ContainsAny(scope.creds.Username+scope.creds.Password, "\"'\\") && !strings.ContainsFunc(scope.creds.Username+scope.creds.Password, unicode.IsSpace)
		if eligible && !netrcHosts[scope.host] {
			fmt.Fprintf(&netrc, "machine %s login %s password %s\n", scope.host, scope.creds.Username, scope.creds.Password)
			netrcHosts[scope.host] = true
		}
	}
	helper += "exit 0\n"
	config := `[credential]
    helper =
    helper = "!sh /run/opctl-auth/helper"
    useHttpPath = true
[core]
    askPass =
`
	files := map[string]string{"helper": helper, "gitconfig": config}
	if netrc.Len() > 0 {
		files[".netrc"] = netrc.String()
	}
	names := []string{}
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	var archive bytes.Buffer
	writer := tar.NewWriter(&archive)
	for _, name := range names {
		data := files[name]
		if err := writer.WriteHeader(&tar.Header{Name: name, Mode: 0600, Size: int64(len(data))}); err != nil {
			return bootstrapAuth{}, err
		}
		if _, err := writer.Write([]byte(data)); err != nil {
			return bootstrapAuth{}, err
		}
	}
	if err := writer.Close(); err != nil {
		return bootstrapAuth{}, err
	}
	if archive.Len() > containerruntime.BootstrapInputLimit {
		return bootstrapAuth{}, fmt.Errorf("Go bootstrap auth archive exceeds 1 MiB")
	}
	sorted := []string{}
	for prefix := range prefixes {
		sorted = append(sorted, prefix)
	}
	sort.Strings(sorted)
	return bootstrapAuth{archive: archive.Bytes(), private: strings.Join(sorted, ",")}, nil
}
