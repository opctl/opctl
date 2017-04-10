package pkg

// Validate validates an opspec package
func (this pkg) Validate(
	pkgRef string,
) []error {
	return this.validator.Validate(pkgRef)
}
