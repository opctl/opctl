package pkg

func (this pkg) Validate(
	pkgRef string,
) (errs []error) {
	errs = this.validator.Validate(pkgRef)
	return
}
