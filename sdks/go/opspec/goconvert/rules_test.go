package goconvert

import (
	"context"
	"strings"
	"testing"
)

func TestControlFlowRefusals(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"root empty", "name: test\nrun: {serial: [{if: [{eq: [a, b]}], container: {image: {ref: alpine}}}]}\n", "empty op graph"},
		{"secret predicate", "name: test\ninputs: {token: {string: {isSecret: true}}}\nrun: {serial: [{if: [{eq: ['$(token)', a]}], container: {image: {ref: alpine}}}]}\n", "plan-readable"},
		{"until", "name: test\nrun: {serialLoop: {until: [{eq: [a, b]}], run: {container: {image: {ref: alpine}}}}}\n", "until"},
		{"exists", "name: test\nrun: {serial: [{if: [{exists: '$(./x)'}], container: {image: {ref: alpine}}}]}\n", "runtime data"},
		{"produced path", "name: test\nrun: {serial: [{container: {image: {ref: alpine}, dirs: {/x: '$(out)'}}}, {container: {image: {ref: alpine}, dirs: {/x: '$(out/sub)'}}}]}\n", "produced"},
		{"parallel collision", "name: test\nrun: {parallel: [{container: {image: {ref: alpine}, stdout: '$(x)'}}, {container: {image: {ref: alpine}, stdout: '$(x)'}}]}\n", "two parallel siblings"},
		{"loop output", "name: test\noutputs: {x: {string: {}}}\nrun: {serialLoop: {range: [a], run: {container: {image: {ref: alpine}, stdout: '$(x)'}}}}\n", "inside a loop"},
		{"null default", "name: test\ninputs: {x: {string: {default: null}}}\nrun: {container: {image: {ref: alpine}}}\n", "null defaults"},
		{"wrong default", "name: test\ninputs: {x: {string: {default: 2}}}\nrun: {container: {image: {ref: alpine}}}\n", "wrong type"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, err := Analyze(context.Background(), t.TempDir(), []byte(tc.source), nil)
			if err != nil {
				if tc.name == "null default" && strings.Contains(err.Error(), "inputs.x.string.default") {
					return
				}
				t.Fatal(err)
			}
			for _, r := range a.Refusals {
				if strings.Contains(r.Reason, tc.want) {
					return
				}
			}
			t.Fatalf("want %q in %+v", tc.want, a.Refusals)
		})
	}
}
