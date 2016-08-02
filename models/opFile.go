package models

type OpFileRunInstruction struct {
  Container *Container `yaml:"container,omitempty"`
  SubOps    []SubOpRunInstruction `yaml:"subOps,omitempty"`
}

type OpFile struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  Inputs      []Parameter `yaml:"inputs,omitempty"`
  Outputs     []Parameter `yaml:"outputs,omitempty"`
  Run         OpFileRunInstruction `yaml:"run,omitempty"`
  Version     string `yaml:"version,omitempty"`
}
