package format

//go:generate counterfeiter -o ./fakeFormat.go --fake-name FakeFormat ./ Format

// a data format implements methods to convert between itself and native go types
type Format interface {
  From(
  in interface{},
  ) (out []byte, err error)

  To(
  in []byte,
  out interface{},
  ) (err error)
}
