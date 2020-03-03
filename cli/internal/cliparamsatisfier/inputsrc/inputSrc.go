package inputsrc

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/inputSrc.go . InputSrc
type InputSrc interface {
	// ReadString returns the value (if any) and true/false to indicate whether the read was successful
	ReadString(
		inputName string,
	) (*string, bool)
}
