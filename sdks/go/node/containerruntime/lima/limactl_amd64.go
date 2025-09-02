//go:build amd64

package lima

import _ "embed"

//go:embed embeds/amd64/limactl
var limactlBytes []byte

//go:embed embeds/amd64/lima-guestagent.Linux-amd64
var limaLinuxGuestAgentBytes []byte
