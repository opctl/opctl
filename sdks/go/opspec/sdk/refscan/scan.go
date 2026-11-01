// Package refscan scans opspec interpolation without resolving references.
// It has no dependencies so the Go op SDK and interpreter share its grammar.
package refscan

import "strings"

// Scan visits unescaped literal text and whole $(...) references in order.
// References include their delimiters and preserve nested references verbatim.
func Scan(expr string, literal func(string) error, reference func(string) error) error {
	return scan(expr, literal, reference, false)
}

// ScanBody scans the inside of a reference, where backslashes are literal.
func ScanBody(expr string, literal, reference func(string) error) error {
	return scan(expr, literal, reference, true)
}

func scan(expr string, literal func(string) error, reference func(string) error, inside bool) error {
	var text strings.Builder
	flush := func() error {
		if text.Len() == 0 {
			return nil
		}
		s := text.String()
		text.Reset()
		return literal(s)
	}
	for i := 0; i < len(expr); {
		if !inside && expr[i] == '\\' {
			start := i
			for i < len(expr) && expr[i] == '\\' {
				i++
			}
			n := i - start
			if i < len(expr) && expr[i] == '$' {
				text.WriteString(strings.Repeat("\\", n/2))
				if n%2 == 1 {
					text.WriteByte('$')
					i++
					continue
				}
			} else {
				text.WriteString(expr[start:i])
				continue
			}
		}
		if strings.HasPrefix(expr[i:], "$(") {
			end := referenceEnd(expr, i+2)
			if end >= 0 {
				if err := flush(); err != nil {
					return err
				}
				if err := reference(expr[i : end+1]); err != nil {
					return err
				}
				i = end + 1
				continue
			}
			text.WriteString("$(")
			if err := flush(); err != nil {
				return err
			}
			return scan(expr[i+2:], literal, reference, true)
		}
		text.WriteByte(expr[i])
		i++
	}
	return flush()
}

func referenceEnd(expr string, start int) int {
	for i := start; i < len(expr); i++ {
		if strings.HasPrefix(expr[i:], "$(") {
			i = referenceEnd(expr, i+2)
			if i < 0 {
				return -1
			}
		} else if expr[i] == ')' {
			return i
		}
	}
	return -1
}

// Contains reports whether expr contains a nonempty, single-line reference.
func Contains(expr string) bool {
	found := false
	_ = Scan(expr, func(string) error { return nil }, func(ref string) error {
		if len(ref) > 3 && !strings.ContainsAny(ref, "\r\n") {
			found = true
		}
		return nil
	})
	return found
}
