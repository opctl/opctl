package attach

import "slices"

// focus tracks running containers in start order.
type focus struct {
	running []string
	current string
}

func (f *focus) Current() string { return f.current }

func (f *focus) Start(id string) {
	if !slices.Contains(f.running, id) {
		f.running = append(f.running, id)
		f.current = id
	}
}

func (f *focus) End(id string) {
	if i := slices.Index(f.running, id); i >= 0 {
		f.running = slices.Delete(f.running, i, i+1)
	}
	if f.current == id {
		f.current = ""
		if len(f.running) > 0 {
			f.current = f.running[len(f.running)-1]
		}
	}
}

func (f *focus) Cycle() {
	if i := slices.Index(f.running, f.current); i > 0 {
		f.current = f.running[i-1]
	} else if i == -1 && len(f.running) > 0 {
		f.current = f.running[len(f.running)-1]
	} else {
		f.current = ""
	}
}
