package opspec

//go:generate counterfeiter -o ./fakeFormat.go --fake-name fakeFormat ./ format

// a data format implements methods to convert between itself and native go types
type format interface {
  From(
  in interface{},
  ) (out []byte, err error)

  To(
  in []byte,
  out interface{},
  ) (err error)
}
