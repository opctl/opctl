package models

func NewEnvVar(
_default string,
name string,
) *EnvVar {

  return &EnvVar{
    Default:_default,
    Name:name,
  }

}

type EnvVar struct {
  Default string `yaml:"default"`
  Name    string `yaml:"name"`
}
