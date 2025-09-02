package embedded

import _ "embed"

//go:embed embeds/linux-arm64/limactl
var limactlBytes []byte

//go:embed embeds/linux-arm64/lima-guestagent.Linux-aarch64
var limaLinuxGuestAgentBytes []byte
