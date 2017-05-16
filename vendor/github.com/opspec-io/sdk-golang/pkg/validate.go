package pkg

import "path/filepath"

// Validate validates an opspec package
func (this _Pkg) Validate(
	pkgPath string,
) []error {
	return this.manifest.Validate(filepath.Join(pkgPath, OpDotYmlFileName))
}
