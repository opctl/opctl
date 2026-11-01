package opfile

import (
	"fmt"
	"testing"
)

func TestContainerOutputDeclarationsSchema(t *testing.T) {
	for _, tc := range []struct {
		declaration string
		valid       bool
	}{
		{`{"count":{"number":{"constraints":{"minimum":1}}}}`, true},
		{`{"count":{"number":{}},"unbound":{"string":{"isSecret":true}}}`, true},
		{`{}`, false}, {`null`, false}, {`{"bad.key":{"string":{}}}`, false}, {`{"count":null}`, false},
		{`{"count":{"number":{},"string":{}}}`, false}, {`{"count":{"number":{"default":0}}}`, false},
		{`{"count":{"number":{"default":null}}}`, false}, {`{"count":{"dir":{}}}`, false}, {`{"count":{"file":{}}}`, false}, {`{"count":{"socket":{}}}`, false},
		{`{"count":{"boolean":{"isSecret":true}}}`, false}, {`{"count":{"number":{},"unknown":true}}`, false},
	} {
		data := fmt.Sprintf(`{"name":"typed","run":{"container":{"image":{"ref":"alpine"},"outputs":{"count":"$(count)"},"outputParams":%s}}}`, tc.declaration)
		_, err := Unmarshal([]byte(data))
		if (err == nil) != tc.valid {
			t.Fatalf("declaration %s: valid=%v, error=%v", tc.declaration, tc.valid, err)
		}
	}
	_, err := Unmarshal([]byte(`{"name":"typed","run":{"container":{"image":{"ref":"alpine"},"outputParams":{"count":{"number":{}}}}}}`))
	if err == nil {
		t.Fatal("declarations without bindings accepted")
	}
}
