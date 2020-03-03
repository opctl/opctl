// Package coerce implements typed data coercion
package coerce

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// Coerce exposes use cases supported by the coerce package
//counterfeiter:generate -o fakes/coerce.go . Coerce
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
