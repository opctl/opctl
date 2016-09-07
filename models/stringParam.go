package models

type StringParam struct {
  Name        string `yaml:"name"`
  Default     string `yaml:"default"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
}
