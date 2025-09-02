package embedded

import _ "embed"

//go:embed embeds/linux-amd64/limactl
var limactlBytes []byte

//go:embed embeds/linux-amd64/lima-guestagent.Linux-x86_64
var limaLinuxGuestAgentBytes []byte
