package pkg

// Validate validates an opspec package
func (this pkg) Validate(
	pkgPath string,
) []error {
	return this.manifestValidator.Validate(pkgPath)
}
