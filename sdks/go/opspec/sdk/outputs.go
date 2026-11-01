package sdk

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
)

var planMu sync.Mutex
var activePlan *outputPlan

type outputPlan struct {
	visible map[string]bool
	counter int
	handles map[string]outputSite
}

type outputSite struct {
	kind, site     string
	explicitResult bool
}

// NewDir creates a fresh directory handle. The first container to mount it
// produces it; later serial calls may mount it, including after a parallel group.
// Parallel siblings cannot produce the same handle. Visible handles mounted in
// parallel are shared and mutable, so concurrent writes need coordination.
// Call NewDir only inside the function passed to Main or Plan. Build the graph
// deterministically in one goroutine, and iterate over sorted map keys.
func NewDir() Dir { return Dir{newOutput("dir")} }

// NewFile creates a fresh file handle with the same visibility rules as NewDir.
// Prefer directory outputs for tools that replace files by rename: a mounted
// file can be written in place, but replacing that mount point fails.
func NewFile() File { return File{newOutput("file")} }

func newOutput(kind string) inputValue {
	if activePlan == nil {
		name := "NewDir"
		if kind == "file" {
			name = "NewFile"
		}
		return inputValue{err: fmt.Errorf("%s must be called while building the op graph (inside the function passed to Main)", name)}
	}
	activePlan.counter++
	name := fmt.Sprintf("opctl_%s_%d", kind, activePlan.counter)
	activePlan.handles[name] = outputSite{kind: kind, site: outputCreationSite()}
	return inputValue{name: name}
}

func outputCreationSite() string {
	_, source, _, _ := runtime.Caller(0)
	var callers [16]uintptr
	frames := runtime.CallersFrames(callers[:runtime.Callers(2, callers[:])])
	for {
		frame, more := frames.Next()
		internal := false
		if filepath.Dir(frame.File) == filepath.Dir(source) {
			switch filepath.Base(frame.File) {
			case "outputs.go", "container.go", "op.go", "container_outputs.go", "typed_outputs.go":
				internal = true
			}
		}
		if !internal || !more {
			return fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		}
	}
}

func newCallResult(kind string) inputValue {
	v := newOutput(kind)
	v.result = true
	if v.err == nil {
		site := activePlan.handles[v.name]
		site.explicitResult = true
		activePlan.handles[v.name] = site
	}
	return v
}

// Mirror node/serialCaller.go and node/parallelCaller.go: serial accumulates
// outputs, while every parallel child starts with the same incoming scope.
func (p *outputPlan) validate(root *callSpec) error {
	p.visible = map[string]bool{}
	return p.walk(root, "run", p.visible, map[string]string{})
}

func copyVisible(visible map[string]bool) map[string]bool {
	result := make(map[string]bool, len(visible))
	for name := range visible {
		result[name] = true
	}
	return result
}

func (p *outputPlan) walk(call *callSpec, path string, visible map[string]bool, producers map[string]string) error {
	if call.Op != nil {
		op := call.Op
		location := fmt.Sprintf("%s (%s)", path, op.Ref)
		for _, use := range append(append([]handleUse{}, op.uses...), op.credUses...) {
			name, site, generated := p.output(use.ref)
			if generated && !use.raw && !visible[name] {
				return fmt.Errorf("%s created at %s is used in %s before its producer is visible", site.kind, site.site, location)
			}
		}
		for _, key := range sortedBindings(op.Outputs) {
			name, site, _ := p.output(op.Outputs[key])
			if earlier, exists := producers[name]; exists {
				return fmt.Errorf("%s created at %s is produced twice by %s and %s", site.kind, site.site, earlier, location)
			}
			producers[name] = location
			visible[name] = true
		}
		return nil
	}
	if call.Container != nil {
		ctr := call.Container
		location := fmt.Sprintf("%s (%s)", path, ctr.Image.Ref)
		for _, key := range sortedBindings(ctr.EnvVars) {
			name, site, fresh := p.output(ctr.EnvVars[key])
			if fresh && !ctr.rawBindings[bindingKey{"env", key}] && !visible[name] {
				return fmt.Errorf("%s created at %s is used by EnvFrom in %s before any earlier call mounts it", site.kind, site.site, location)
			}
		}
		produced := map[string]bool{}
		for _, mount := range []struct {
			kind     string
			bindings map[string]string
		}{{"dir", ctr.Dirs}, {"file", ctr.Files}} {
			for _, key := range sortedBindings(mount.bindings) {
				name, site, fresh := p.output(mount.bindings[key])
				if !fresh || visible[name] || ctr.rawBindings[bindingKey{mount.kind, key}] {
					continue
				}
				if site.explicitResult {
					return fmt.Errorf("%s call output created at %s is consumed by %s before its producer is visible", site.kind, site.site, location)
				}
				if produced[name] {
					return fmt.Errorf("%s created at %s is produced at two mount paths in %s", site.kind, site.site, location)
				}
				if earlier, ok := producers[name]; ok {
					return fmt.Errorf("%s created at %s is mounted by parallel calls %s and %s; parallel calls can't see each other's outputs, so produce it in an earlier serial call", site.kind, site.site, earlier, location)
				}
				produced[name] = true
				producers[name] = location
			}
		}
		refs := []string{ctr.Stdout, ctr.Stderr}
		for _, key := range sortedBindings(ctr.Outputs) {
			refs = append(refs, ctr.Outputs[key])
		}
		for _, ref := range refs {
			if ref == "" {
				continue
			}
			name, site, _ := p.output(ref)
			if earlier, exists := producers[name]; exists {
				return fmt.Errorf("%s created at %s is produced twice by %s and %s", site.kind, site.site, earlier, location)
			}
			producers[name] = location
			produced[name] = true
		}
		for name := range produced {
			visible[name] = true
		}
		return nil
	}
	for i, child := range call.Serial {
		if err := p.walk(child, fmt.Sprintf("%s.serial[%d]", path, i), visible, producers); err != nil {
			return err
		}
	}
	outgoing := copyVisible(visible)
	for i, child := range call.Parallel {
		childVisible := copyVisible(visible)
		if err := p.walk(child, fmt.Sprintf("%s.parallel[%d]", path, i), childVisible, producers); err != nil {
			return err
		}
		for name := range childVisible {
			outgoing[name] = true
		}
	}
	for name := range outgoing {
		visible[name] = true
	}
	return nil
}

func (p *outputPlan) output(ref string) (string, outputSite, bool) {
	if len(ref) < 4 || ref[:2] != "$(" || ref[len(ref)-1] != ')' {
		return "", outputSite{}, false
	}
	name := ref[2 : len(ref)-1]
	site, ok := p.handles[name]
	return name, site, ok
}

func sortedBindings(bindings map[string]string) []string {
	keys := make([]string, 0, len(bindings))
	for key := range bindings {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
