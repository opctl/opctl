package sdk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// StubCall builds calls for generated op stubs. Use CallOp in hand-written code.
// Local refs are module-relative, and unset handles use the callee's defaults.
func StubCall(stub, ref string, local bool, inputs map[string]Handle, required ...string) OpCall {
	if local {
		cwd, err := os.Getwd()
		if err != nil {
			return OpCall{err: fmt.Errorf("op stub %s: %w", stub, err)}
		}
		root := cwd
		for {
			if _, err = os.Stat(filepath.Join(root, "go.mod")); err == nil {
				break
			}
			if !os.IsNotExist(err) {
				return OpCall{err: fmt.Errorf("op stub %s: %w", stub, err)}
			}
			parent := filepath.Dir(root)
			if parent == root {
				return OpCall{err: fmt.Errorf("op stub %s: needs a go.mod in the op directory or a parent", stub)}
			}
			root = parent
		}
		ref, err = filepath.Rel(cwd, filepath.Join(root, filepath.FromSlash(ref)))
		if err != nil {
			return OpCall{err: fmt.Errorf("op stub %s: %w", stub, err)}
		}
		ref = filepath.ToSlash(ref)
		if !strings.HasPrefix(ref, "../") {
			ref = "./" + ref
		}
	}
	values := map[string]any{}
	for name, handle := range inputs {
		if h, ok := handle.(interface{ bindingValue() inputValue }); ok {
			v := h.bindingValue()
			if v.name != "" || v.raw != "" || v.literal || v.err != nil {
				values[name] = handle
			}
		}
	}
	c := CallOp(ref).Inputs(values)
	for _, name := range required {
		if _, ok := values[name]; !ok && c.err == nil {
			c.err = fmt.Errorf("op stub %s: input %q is required", stub, name)
		}
	}
	return c
}
