package interpolater

import (
	"errors"
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	for _, tc := range []struct {
		expr string
		want []string
	}{
		{`a\$(escaped)-$(x)-$(dir/$(key))`, []string{`literal:a$(escaped)-`, `ref:$(x)`, `literal:-`, `ref:$(dir/$(key))`}},
		{`\\$(x)`, []string{`literal:\`, `ref:$(x)`}},
		{`$(unfinished`, []string{`literal:$(`, `literal:unfinished`}},
		{`$(outer$(inner)`, []string{`literal:$(`, `literal:outer`, `ref:$(inner)`}},
	} {
		var got []string
		err := Scan(tc.expr, func(s string) error { got = append(got, "literal:"+s); return nil }, func(s string) error { got = append(got, "ref:"+s); return nil })
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%q: %q, %v; want %q", tc.expr, got, err, tc.want)
		}
	}
	sentinel := errors.New("stop")
	if err := Scan("$(x)", func(string) error { return nil }, func(string) error { return sentinel }); !errors.Is(err, sentinel) {
		t.Fatal(err)
	}
}
