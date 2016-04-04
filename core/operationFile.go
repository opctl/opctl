package core

type operationFileSubOperation struct {
  Url string `yaml:"url"`
}

type operationFile struct {
  Name          *string `yaml:"name"`
  Description   string `yaml:"description"`
  SubOperations []operationFileSubOperation `yaml:"subOperations,omitempty"`
}