package models

func NewParameter(
name string,
description string,
isSecret bool,
) *Parameter {

  return &Parameter{
    Name:name,
    Description:description,
    IsSecret:isSecret,
  }

}

type Parameter struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
}
