package greeting

import (
	_ "embed"
	"strings"
)

//go:embed greeting.txt
var text string

func Text() string { return strings.TrimSpace(text) }
