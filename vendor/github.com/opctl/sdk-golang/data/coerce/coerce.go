// Package coerce implements typed data coercion
package coerce

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Coerce

// Coerce exposes use cases supported by the coerce package
type Coerce interface {
	toArrayer
	toBooleaner
	toFiler
	toNumberer
	toObjecter
	toStringer
}

// New returns an initialized Coerce instance
func New() Coerce {
	return struct {
		toArrayer
		toBooleaner
		toFiler
		toNumberer
		toObjecter
		toStringer
	}{
		newToArrayer(),
		newToBooleaner(),
		newToFiler(),
		newToNumberer(),
		newToObjecter(),
		newToStringer(),
	}
}
