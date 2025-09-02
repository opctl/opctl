package embedded

import _ "embed"

//go:embed embeds/darwin-arm64/limactl
var limactlBytes []byte

//go:embed embeds/darwin-arm64/lima-guestagent.Linux-aarch64
var limaLinuxGuestAgentBytes []byte
