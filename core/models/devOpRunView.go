package models

import "encoding/json"

type devOpRunView struct {
  DevOpName          string `json:"devOpName"`
  StartedAtEpochTime int64 `json:"startedAtEpochTime"`
  EndedAtEpochTime   int64 `json:"endedAtEpochTime"`
  ExitCode           int `json:"exitCode"`
}

type DevOpRunView struct {
  devOpRunView
}

func newDevOpRunView(
devOpName          string,
startedAtEpochTime int64,
endedAtEpochTime   int64,
exitCode           int,
) DevOpRunView {

  return DevOpRunView{
    devOpRunView{
      DevOpName:devOpName,
      StartedAtEpochTime:startedAtEpochTime,
      EndedAtEpochTime:endedAtEpochTime,
      ExitCode:exitCode,
    },
  }

}

func (this DevOpRunView) MarshalJSON(
) ([]byte, error) {
  return json.Marshal(this.devOpRunView)
}

func (this *DevOpRunView) UnmarshalJSON(
b []byte,
) error {
  return json.Unmarshal(b, &this.devOpRunView)
}

func (this DevOpRunView) DevOpName() string {
  return this.devOpRunView.DevOpName
}

func (this DevOpRunView) StartedAtEpochTime() int64 {
  return this.devOpRunView.StartedAtEpochTime
}

func (this DevOpRunView) EndedAtEpochTime() int64 {
  return this.devOpRunView.EndedAtEpochTime
}

func (this DevOpRunView) ExitCode() int {
  return this.devOpRunView.ExitCode
}