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
	return &cliColorer{
		attentionCliColorer: color.New(color.FgHiYellow, color.Bold),
		errorCliColorer:     color.New(color.FgHiRed, color.Bold),
		infoCliColorer:      color.New(color.FgHiCyan, color.Bold),
		successCliColorer:   color.New(color.FgHiGreen, color.Bold),
	}
}

type cliColorer struct {
	attentionCliColorer *color.Color
	errorCliColorer     *color.Color
	infoCliColorer      *color.Color
	successCliColorer   *color.Color
}

func (this cliColorer) Disable() {
	this.attentionCliColorer.DisableColor()
	this.errorCliColorer.DisableColor()
	this.infoCliColorer.DisableColor()
	this.successCliColorer.DisableColor()
}

func (this cliColorer) Attention(
	format string,
	values ...interface{},
) string {
	return this.attentionCliColorer.SprintfFunc()(format, values...)
}

func (this cliColorer) Error(
	format string,
	values ...interface{},
) string {
	return this.errorCliColorer.SprintfFunc()(format, values...)
}

func (this cliColorer) Info(
	format string,
	values ...interface{},
) string {
	return this.infoCliColorer.SprintfFunc()(format, values...)
}

func (this cliColorer) Success(
	format string,
	values ...interface{},
) string {
	return this.successCliColorer.SprintfFunc()(format, values...)
}
