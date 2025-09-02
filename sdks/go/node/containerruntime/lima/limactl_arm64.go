//go:build arm64

package lima

import _ "embed"

//go:embed embeds/arm64/limactl
var limactlBytes []byte

//go:embed embeds/arm64/lima-guestagent.Linux-aarch64
var limaLinuxGuestAgentBytes []byte
