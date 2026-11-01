package goop

import (
	"fmt"
	"strings"
	"testing"
)

func TestTypedContainerOutputSchemaGate(t *testing.T) {
	for _, version := range []int{1, 2, 3, 4} {
		data := fmt.Sprintf(`{"contractVersion":1,"opspecSchema":%d,"opspec":{"name":"typed","run":{"serial":[{"container":{"image":{"ref":"alpine"},"outputs":{"count":"$(count)"},"outputParams":{"count":{"number":{}}}}}]}}}`, version)
		_, err := ParsePlan([]byte(data))
		if version == 4 && err != nil {
			t.Fatal(err)
		}
		if version < 4 && (err == nil || !strings.Contains(err.Error(), "schema 4")) {
			t.Fatalf("mislabeled typed graph was not rejected: %v", err)
		}
	}
}
