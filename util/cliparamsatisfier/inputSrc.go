package cliparamsatisfier

//go:generate counterfeiter -o ./fakeInputSrc.go --fake-name FakeInputSrc ./ InputSrc

type InputSrc interface {
	Read(
		inputName string,
	) *string
}
