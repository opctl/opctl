package manifest

// Validate validates the pkg manifest at path
func (this _Manifest) Validate(
	path string,
) []error {
	return this.validator.Validate(path)
}
