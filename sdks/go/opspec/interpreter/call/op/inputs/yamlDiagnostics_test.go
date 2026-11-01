package inputs

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestYAMLInputInterpretationPrintsDiagnostic(t *testing.T) {
	for _, sensitive := range []bool{false, true} {
		reader, writer, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		original := os.Stdout
		os.Stdout = writer
		_, failure := Interpret(map[string]interface{}{"unknown": "x"}, nil, "", nil, "", map[string]bool{"unknown": sensitive})
		writer.Close()
		os.Stdout = original
		output, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			t.Fatal(err)
		}
		if failure == nil {
			t.Fatal("expected input error")
		}
		if sensitive && len(output) != 0 {
			t.Fatal("sensitive diagnostic printed")
		}
		if !sensitive && !strings.Contains(string(output), "unknown ") {
			t.Fatal("legacy YAML diagnostic missing")
		}
	}
}
