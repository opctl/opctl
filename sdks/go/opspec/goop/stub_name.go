package goop

import (
	"go/token"
	"regexp"
)

var stubPackageName = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

func validStubName(name string) bool {
	if !stubPackageName.MatchString(name) || token.Lookup(name).IsKeyword() {
		return false
	}
	switch name {
	case "main", "internal", "vendor", "testdata":
		return false
	}
	return true
}
