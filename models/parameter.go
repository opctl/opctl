package models

func NewParameter(
name string,
_default string,
description string,
isSecret bool,
) *Parameter {

  return &Parameter{
    Name:name,
    Default:_default,
    Description:description,
    IsSecret:isSecret,
  }

}

type Parameter struct {
  Name        string `yaml:"name"`
  Default     string `yaml:"default"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
}
