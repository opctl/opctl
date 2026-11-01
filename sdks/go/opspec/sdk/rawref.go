package sdk

import (
	"fmt"
	"regexp"
)

var wholeReference = regexp.MustCompile(`^\$\(.+\)$`)

// RawRef emits a whole $(...) reference verbatim, without visibility checks.
// It supports op-directory paths, subpaths, and nested references the SDK does
// not model. Value cannot be read during plan. Building expr from input values
// can inject references into the op graph; use typed handles when possible.
// Literal Cmd, Env, WorkDir, and image strings still escape every $.
func RawRef[T Handle](expr string) T {
	var value T
	setter, ok := any(&value).(interface{ setRaw(string) })
	if !ok {
		panic(inputError{fmt.Errorf("RawRef needs an opctl handle value")})
	}
	setter.setRaw(expr)
	return value
}

func (v *inputValue) setRaw(expr string) {
	if !wholeReference.MatchString(expr) {
		v.err = fmt.Errorf("RawRef needs a whole $(...) reference, got %q", expr)
		return
	}
	v.raw = expr
}

func (v inputValue) isRaw() bool { return v.raw != "" }
