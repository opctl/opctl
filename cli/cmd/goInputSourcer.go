package cmd

import (
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
)

// A rejected explicit Go input must prompt again, rather than use its default.
type goInputSourcer struct {
	provided cliparamsatisfier.InputSourcer
	defaults inputsrc.InputSrc
	prompt   inputsrc.InputSrc
	seen     map[string]bool
}

func (s goInputSourcer) Source(name string) (*string, bool) {
	if !s.seen[name] {
		s.seen[name] = true
		if value, ok := s.provided.Source(name); ok {
			return value, true
		}
		if value, ok := s.defaults.ReadString(name); ok {
			return value, true
		}
	}
	return s.prompt.ReadString(name)
}
