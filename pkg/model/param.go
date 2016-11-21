package model

type Param struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
  Default     string `yaml:"default,omitempty"`
}
