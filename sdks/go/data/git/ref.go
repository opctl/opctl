package git

import (
	"fmt"
	"path/filepath"
)

type ref struct {
	Name    string
	Version string
	OpPath  string
}

// ToPath constructs a filesystem path for a Ref, assuming the provided base path
func (pr ref) ToPath(basePath string) string {
	crossPlatPath := filepath.FromSlash(fmt.Sprintf("%v#%v", pr.Name, pr.Version))
	return filepath.Join(basePath, crossPlatPath)
}
