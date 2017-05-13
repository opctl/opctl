package pkg

// Validate validates an opspec package
func (this _Pkg) Validate(
	pkgPath string,
) []error {
	return this.manifestValidator.Validate(pkgPath)
}
