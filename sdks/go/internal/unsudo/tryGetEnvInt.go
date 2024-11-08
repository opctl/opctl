package unsudo

import (
	"os"
	"strconv"
	"strings"
)

func tryGetEnvInt(
	envVarName string,
) *int {
	if envVar := os.Getenv(envVarName); envVar != "" {
		if envVarInt, err := strconv.Atoi(
			strings.TrimSpace(envVar),
		); err == nil {
			return &envVarInt
		}
	}
	return nil
}
