package models

type Param struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
  String      *StringParam `yaml:"string,omitempty"`
}
