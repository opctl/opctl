package sdk

import "fmt"

type groupCall struct {
	parallel bool
	children []Call
}

// Serial runs calls in order. Later calls see earlier calls' outputs.
// Empty groups are omitted; an empty root graph is a plan error.
func Serial(calls ...Call) Call { return groupCall{children: append([]Call(nil), calls...)} }

// Parallel runs calls concurrently. Siblings see only the group's incoming scope.
// Their outputs become visible to later serial calls. Empty groups are omitted.
func Parallel(calls ...Call) Call {
	return groupCall{parallel: true, children: append([]Call(nil), calls...)}
}

func (groupCall) call() {}

func (g groupCall) render() (*callSpec, error) {
	children := make([]*callSpec, 0, len(g.children))
	for _, child := range g.children {
		if child == nil {
			return nil, fmt.Errorf("Go op plan must return a call")
		}
		spec, err := child.render()
		if err != nil {
			return nil, err
		}
		if spec != nil {
			children = append(children, spec)
		}
	}
	if len(children) == 0 {
		return nil, nil
	}
	if g.parallel {
		return &callSpec{Parallel: children}, nil
	}
	return &callSpec{Serial: children}, nil
}
