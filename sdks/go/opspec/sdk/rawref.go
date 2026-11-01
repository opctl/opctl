package sdk

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/opspec/sdk/refscan"
	"regexp"
)

var wholeReference = regexp.MustCompile(`^\$\(.+\)$`)

// RawRef emits an expression containing $(...) references verbatim, without visibility checks.
// It supports nested references and other expressions the SDK does not model;
// use OpDir, OpFile, and Dir subpaths for paths. Value cannot be read during plan. Building expr from input values
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
	if !refscan.Contains(expr) {
		v.err = fmt.Errorf("RawRef needs an unescaped $(...) reference, got %q", expr)
		return
	}
	v.raw = expr
}

func (v inputValue) isRaw() bool { return v.raw != "" }
