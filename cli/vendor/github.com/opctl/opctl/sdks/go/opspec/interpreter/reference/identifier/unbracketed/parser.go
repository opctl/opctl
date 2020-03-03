package unbracketed

// Parser parses an identifier by consuming ref up to (but not including) the first unbracketed identifier char i.e '.', '[', '/' encountered
// returns the identifier and ref remainder
//counterfeiter:generate -o fakes/parser.go . Parser
type Parser interface {
	parser
}

// parser is an internal version of Parser so fakes don't cause cyclic deps
//counterfeiter:generate -o internal/fakes/parser.go . parser
type parser interface {
	Parse(ref string) (string, string)
}

func NewParser() Parser {
	return _parser{}
}

type _parser struct{}

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
