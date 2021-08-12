package callprogress

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/opctl/opctl/cli/internal/clitext"
	"github.com/opctl/opctl/sdks/go/model"
)

// CallGraph maintains a record of the current state of an op
type CallGraph struct {
	rootNode *callGraphNode
	errors   []error
}

type callGraphNode struct {
	call      *model.Call
	startTime *time.Time
	endTime   *time.Time
	state     string
	children  []*callGraphNode
}

func newCallGraphNode(call *model.Call, timestamp time.Time) *callGraphNode {
	return &callGraphNode{
		call:      call,
		startTime: &timestamp,
		children:  []*callGraphNode{},
	}
}

var errNotFoundInGraph = errors.New("not found in graph")

const skippedState = "skipped"

func (n *callGraphNode) insert(call *model.Call, startTime time.Time, initialState string) error {
	if call.ParentID == nil {
		return fmt.Errorf("missing parent ID for %s", call.ID)
	}
	if n.call.ID == *call.ParentID {
		node := newCallGraphNode(call, startTime)
		node.state = initialState
		n.children = append(n.children, node)
		return nil
	}
	for _, child := range n.children {
		err := child.insert(call, startTime, initialState)
		if err == nil {
			return nil
		}
	}
	return errNotFoundInGraph
}

func (n *callGraphNode) find(call *model.Call) *callGraphNode {
	if call.ID == n.call.ID {
		return n
	}
	for _, child := range n.children {
		if c := child.find(call); c != nil {
			return c
		}
	}
	return nil
}

func (n *callGraphNode) isLeaf() bool {
	return len(n.children) == 0
}

func (n *callGraphNode) countChildren() int {
	count := 0
	for _, child := range n.children {
		if child.isLeaf() {
			count++
		} else {
			count += child.countChildren()
		}
	}
	return count
}

func (n callGraphNode) String(loader LoadingSpinner, now time.Time, collapseCompleted bool) string {
	var str strings.Builder

	// Graph node indicator
	if n.isLeaf() {
		str.WriteString("◉")
	} else {
		str.WriteString("◎")
	}

	// Leading "status"
	switch n.state {
	case model.OpOutcomeSucceeded:
		str.WriteString(success.Sprint(" ☑"))
	case model.OpOutcomeFailed:
		str.WriteString(failed.Sprint(" ⚠"))
	case model.OpOutcomeKilled:
		str.WriteString("️ ☒")
	case skippedState:
		str.WriteString(" ☐")
	case "":
		// only display loading spinner on leaf nodes
		if n.isLeaf() {
			str.WriteString(" " + loader.String())
		}
	default:
		str.WriteString(n.state)
	}

	call := *n.call

	// "Named" ops
	if call.Name != nil {
		str.WriteString(" " + highlighted.Sprint(*call.Name))
	}

	// Main node description
	var desc string
	if call.Container != nil {
		desc = muted.Sprint(call.Container.ContainerID[:8]) + " "
		if call.Container.Name != nil {
			desc += highlighted.Sprint(*call.Container.Name)
		} else {
			desc += *call.Container.Image.Ref
		}
	} else if call.Op != nil {
		desc = highlighted.Sprint(clitext.FromOpRef(call.Op.OpPath))
	} else if call.Parallel != nil {
		desc = "parallel"
	} else if call.ParallelLoop != nil {
		desc = "parallel loop"
	} else if call.Serial != nil {
		desc = "serial"
	} else if call.SerialLoop != nil {
		desc = "serial loop"
	}

	collapsed := n.state == model.OpOutcomeSucceeded && !n.isLeaf() && collapseCompleted

	if call.If != nil {
		str.WriteString(" if")
		// this means it was skipped
		if desc == "" {
			str.WriteString(" " + muted.Sprint("skipped"))
		} else {
			str.WriteString("\n")
			if n.isLeaf() || collapsed {
				str.WriteString(" ")
			} else {
				str.WriteString("│")
			}
		}
	}

	if desc != "" {
		str.WriteString(" " + desc)
	}

	// Time elapsed
	if n.startTime != nil && n.state != skippedState {
		if n.endTime != nil { // if done
			str.WriteString(" " + n.endTime.Sub(*n.startTime).String())
		} else if n.isLeaf() { // only display live time for leaf nodes, like loading spinner
			// don't show milliseconds - they're not really understandable
			str.WriteString(" " + now.Sub(*n.startTime).Round(time.Second).String())
		}
	}

	// Add the command invoked by a container if it's not named
	if call.Container != nil && call.Container.Name == nil && len(call.Container.Cmd) > 0 {
		str.WriteString(" " + muted.Sprint(strings.ReplaceAll(strings.Join(call.Container.Cmd, " "), "\n", "\\n")))
	}

	// Collapsed nodes
	if collapsed {
		str.WriteString(" ")
		childCount := n.countChildren()
		if childCount == 1 {
			str.WriteString(muted.Sprint("(1 child)"))
		} else {
			str.WriteString(muted.Sprintf("(%d children)", childCount))
		}
		return str.String()
	}

	// Children
	childLen := len(n.children)
	for i, child := range n.children {
		childLines := strings.Split(child.String(loader, now, collapseCompleted), "\n")
		for j, part := range childLines {
			if j == 0 {
				if i < childLen-1 {
					str.WriteString(fmt.Sprintf("\n├─%s", part))
				} else {
					str.WriteString(fmt.Sprintf("\n└─%s", part))
				}
			} else if i < childLen-1 {
				str.WriteString(fmt.Sprintf("\n│ %s", part))
			} else {
				str.WriteString(fmt.Sprintf("\n  %s", part))
			}
		}
	}

	return str.String()
}

// String returns a visual representation of the current state of the call graph
func (g CallGraph) String(loader LoadingSpinner, now time.Time, collapseCompleted bool) string {
	var str strings.Builder
	if g.rootNode == nil {
		return "Empty call graph"
	}
	str.WriteString(g.rootNode.String(loader, now, collapseCompleted))
	for _, err := range g.errors {
		str.WriteString("\n" + warning.Sprint("⚠️  ") + err.Error())
	}
	return str.String()
}

// HandleEvent accepts an opctl event and updates the call graph appropriately
func (g *CallGraph) HandleEvent(event *model.Event) error {
	if event.CallStarted != nil {
		if event.CallStarted.Call.ParentID == nil {
			if g.rootNode == nil {
				g.rootNode = newCallGraphNode(&event.CallStarted.Call, event.Timestamp)
				return nil
			}
			return errors.New("parent node already set")
		}
		return g.rootNode.insert(&event.CallStarted.Call, event.Timestamp, "")
	} else if event.CallEnded != nil {
		if g.rootNode == nil {
			return nil
		}
		node := g.rootNode.find(&event.CallEnded.Call)
		if node == nil {
			return g.rootNode.insert(&event.CallEnded.Call, event.Timestamp, skippedState)
		}
		node.endTime = &event.Timestamp
		node.state = event.CallEnded.Outcome
	}
	return nil
}
