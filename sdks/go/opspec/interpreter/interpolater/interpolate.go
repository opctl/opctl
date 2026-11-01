package interpolater

import (
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"
	"github.com/opctl/opctl/sdks/go/opspec/sdk/refscan"
	"strings"
)

const (
	RefStart = "$("
	RefEnd   = ")"
)

// Scan visits unescaped literal text and whole references, preserving nested references.
func Scan(expr string, onLiteral, onRef func(string) error) error {
	return refscan.Scan(expr, onLiteral, onRef)
}

// Interpolate resolves opspec references using the shared interpolation scanner.
func Interpolate(expression string, scope map[string]*model.Value) (string, error) {
	return interpolate(expression, scope, false)
}

func interpolate(expression string, scope map[string]*model.Value, inside bool) (string, error) {
	var out strings.Builder
	scan := refscan.Scan
	if inside {
		scan = refscan.ScanBody
	}
	err := scan(expression, func(s string) error { out.WriteString(s); return nil }, func(ref string) error {
		body, err := interpolate(ref[2:len(ref)-1], scope, true)
		if err != nil {
			return err
		}
		v, err := reference.Interpret(opspec.NameToRef(body), scope, nil)
		if err != nil {
			return err
		}
		text, err := coerce.ToString(v)
		if err != nil {
			return err
		}
		out.WriteString(*text.String)
		return nil
	})
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
