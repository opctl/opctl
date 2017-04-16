package pkg

// ResolveRef resolves a pkgRef according to opspec package resolution rules.
func (this pkg) ResolveRef(
	pkgRef string,
) string {
	return this.refResolver.Resolve(pkgRef)
}
