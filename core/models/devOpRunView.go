package models

type DevOpRunView struct {
  Id                *string
  DevOpName         string
  StartedAtUnixTime int64
  EndedAtUnixTime   int64
  ExitCode          int
}