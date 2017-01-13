package colorer

//go:generate counterfeiter -o ./fakeColorer.go --fake-name FakeColorer ./ Colorer

import (
	"github.com/fatih/color"
)

type Colorer interface {
	// silently disables coloring
	Disable()

	// attention colors
	Attention(
		format string,
		values ...interface{},
	) string

	// errors collors
	Error(
		format string,
		values ...interface{},
	) string

	// info colors
	Info(
		format string,
		values ...interface{},
	) string

	// success colors
	Success(
		format string,
		values ...interface{},
	) string
}

func New() Colorer {
  color.NoColor = false
	return &colorer{
		attentionColorer: color.New(color.FgHiYellow, color.Bold).SprintfFunc(),
		errorColorer:     color.New(color.FgHiRed, color.Bold).SprintfFunc(),
		infoColorer:      color.New(color.FgHiCyan, color.Bold).SprintfFunc(),
		successColorer:   color.New(color.FgHiGreen, color.Bold).SprintfFunc(),
	}
}

type colorer struct {
	attentionColorer func(format string, a ...interface{}) string
	errorColorer     func(format string, a ...interface{}) string
	infoColorer      func(format string, a ...interface{}) string
	successColorer   func(format string, a ...interface{}) string
}

func (this colorer) Disable() {
	color.NoColor = true
}

func (this colorer) Attention(
	format string,
	values ...interface{},
) string {
	return this.attentionColorer(format, values...)
}

func (this colorer) Error(
	format string,
	values ...interface{},
) string {
	return this.errorColorer(format, values...)
}

func (this colorer) Info(
	format string,
	values ...interface{},
) string {
	return this.infoColorer(format, values...)
}

func (this colorer) Success(
	format string,
	values ...interface{},
) string {
	return this.successColorer(format, values...)
}
