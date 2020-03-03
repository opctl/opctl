package internal

import (
	"fmt"
	"path/filepath"
)

type Ref struct {
	Name    string
	Version string
	OpPath  string
}

// ToPath constructs a filesystem path for a Ref, assuming the provided base path
func (pr Ref) ToPath(basePath string) string {
	crossPlatPath := filepath.FromSlash(fmt.Sprintf("%v#%v", pr.Name, pr.Version))
	return filepath.Join(basePath, crossPlatPath)
}
