package clicolorer

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/fatih/color"
)

//counterfeiter:generate -o fakes/cliColorer.go . CliColorer
type CliColorer interface {
	// silently disables coloring
	Disable()

	// attention colors
	Attention(
		s string,
	) string

	// errors colors
	Error(
		s string,
	) string

	// info colors
	Info(
		s string,
	) string

	// success colors
	Success(s string) string
}

func New() CliColorer {
	attentionCliColorer := color.New(color.FgHiYellow, color.Bold)
	attentionCliColorer.EnableColor()

	errorCliColorer := color.New(color.FgHiRed, color.Bold)
	errorCliColorer.EnableColor()

	infoCliColorer := color.New(color.FgHiCyan, color.Bold)
	infoCliColorer.EnableColor()

	successCliColorer := color.New(color.FgHiGreen, color.Bold)
	successCliColorer.EnableColor()

	return &cliColorer{
		attentionCliColorer: attentionCliColorer,
		errorCliColorer:     errorCliColorer,
		infoCliColorer:      infoCliColorer,
		successCliColorer:   successCliColorer,
	}
}

type cliColorer struct {
	attentionCliColorer *color.Color
	errorCliColorer     *color.Color
	infoCliColorer      *color.Color
	successCliColorer   *color.Color
}

func (this *cliColorer) Disable() {
	this.attentionCliColorer.DisableColor()
	this.errorCliColorer.DisableColor()
	this.infoCliColorer.DisableColor()
	this.successCliColorer.DisableColor()
}

func (this cliColorer) Attention(
	s string,
) string {
	return this.attentionCliColorer.SprintfFunc()(s)
}

func (this cliColorer) Error(
	s string,
) string {
	return this.errorCliColorer.SprintFunc()(s)
}

func (this cliColorer) Info(
	s string,
) string {
	return this.infoCliColorer.SprintFunc()(s)
}

func (this cliColorer) Success(
	s string,
) string {
	return this.successCliColorer.SprintFunc()(s)
}
