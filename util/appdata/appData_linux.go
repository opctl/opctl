package appdata

import (
	"fmt"
)

const (
	globalPathTemplate  = "/var/lib"
	perUserPathTemplate = "%v"
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
