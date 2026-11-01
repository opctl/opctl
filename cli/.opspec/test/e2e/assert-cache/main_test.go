package main

import (
	"github.com/opctl/opctl/sdks/go/model"
	"testing"
)

func TestCacheLifecycleAssertions(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		builds, describes, plans int
		warm, wantError          bool
	}{
		{"cold", 1, 1, 1, false, false}, {"warm", 0, 0, 1, true, false}, {"recompile", 1, 2, 1, true, true}, {"missing plan", 0, 2, 0, true, true},
	} {
		events := []model.Event{}
		cached := 1
		if tc.warm {
			cached = 2
		}
		cachedName := "bootstrap describe (cached)"
		for i := 0; i < cached; i++ {
			events = append(events, model.Event{CallStarted: &model.CallStarted{Call: model.Call{Name: &cachedName}}})
		}
		for name, count := range map[string]int{"bootstrap build": tc.builds, "bootstrap describe": tc.describes, "bootstrap plan": tc.plans} {
			for i := 0; i < count; i++ {
				name := name
				events = append(events, model.Event{CallStarted: &model.CallStarted{Call: model.Call{Name: &name}}})
			}
		}
		if err := checkCacheEvents(events, tc.warm); (err != nil) != tc.wantError {
			t.Fatalf("%s: %v", tc.name, err)
		}
	}
}
