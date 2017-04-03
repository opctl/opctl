package clicolorer

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliColorer

import (
	"github.com/fatih/color"
)

type CliColorer interface {
	// silently disables coloring
	Disable()

	// attention colors
	Attention(
		format string,
		values ...interface{},
	) string

	// errors colors
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

func New() CliColorer {
	color.NoColor = false
	return &cliColorer{
		attentionCliColorer: color.New(color.FgHiYellow, color.Bold).SprintfFunc(),
		errorCliColorer:     color.New(color.FgHiRed, color.Bold).SprintfFunc(),
		infoCliColorer:      color.New(color.FgHiCyan, color.Bold).SprintfFunc(),
		successCliColorer:   color.New(color.FgHiGreen, color.Bold).SprintfFunc(),
	}
}

type cliColorer struct {
	attentionCliColorer func(format string, a ...interface{}) string
	errorCliColorer     func(format string, a ...interface{}) string
	infoCliColorer      func(format string, a ...interface{}) string
	successCliColorer   func(format string, a ...interface{}) string
}

func (this cliColorer) Disable() {
	color.NoColor = true
}

func (this cliColorer) Attention(
	format string,
	values ...interface{},
) string {
	return this.attentionCliColorer(format, values...)
}

func (this cliColorer) Error(
	format string,
	values ...interface{},
) string {
	return this.errorCliColorer(format, values...)
}

func (this cliColorer) Info(
	format string,
	values ...interface{},
) string {
	return this.infoCliColorer(format, values...)
}

func (this cliColorer) Success(
	format string,
	values ...interface{},
) string {
	return this.successCliColorer(format, values...)
}
