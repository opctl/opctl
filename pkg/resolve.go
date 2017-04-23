package pkg

// Resolve resolves a local package according to opspec package resolution rules and returns it's absolute path.
func (this pkg) Resolve(
	basePath,
	pkgRef string,
) (string, bool) {
	return this.resolver.Resolve(basePath, pkgRef)
}
