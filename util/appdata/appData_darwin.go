package appdata

import (
	"fmt"
)

const (
	globalPathTemplate  = "/Library/Application Support"
	perUserPathTemplate = "%v/Library/Application Support"
)

func (this appData) GlobalPath() string {
	return globalPathTemplate
}

func (this appData) PerUserPath() string {
	return fmt.Sprintf(
		perUserPathTemplate,
		this.getUserHomePath(),
	)
}
