package core

type operationFileSubOperation struct {
  Name string `yaml:"name"`
}

type operationFile struct {
  Description   string `yaml:"description"`
  SubOperations []operationFileSubOperation `yaml:"subOperations,omitempty"`
}