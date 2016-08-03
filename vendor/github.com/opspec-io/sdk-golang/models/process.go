package models

func NewProcess(
cmd string,
envVars []EnvVar,
workDir string,
) *Process {

  return &Process{
    Cmd:cmd,
    EnvVars:envVars,
    WorkDir:workDir,
  }

}

type Process struct {
  Cmd     string `yaml:"cmd"`
  EnvVars []EnvVar `yaml:"envVars"`
  WorkDir string `yaml:"workDir"`
}
