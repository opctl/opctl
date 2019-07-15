package unbracketed

//go:generate counterfeiter -o ./fakeParser.go --fake-name FakeParser ./ Parser

// Parser parses an identifier by consuming ref up to (but not including) the first unbracketed identifier char i.e '.', '[', '/' encountered
// returns the identifier and ref remainder
type Parser interface {
	Parse(ref string) (string, string)
}

func NewParser() Parser {
	return _parser{}
}

type _parser struct {
}

func (p _parser) Parse(
	ref string,
) (string, string) {

	for i := 0; i < len(ref); i++ {

		if ref[i] == '.' || ref[i] == '[' || ref[i] == '/' {
			// identifier ended by '.', '[', or '/'
			return ref[:i], ref[i:]
		}
	}

	return ref, ""
}
