package model

type OpManifest struct {
  Manifest `yaml:",inline"`
  Inputs  []Param `yaml:"inputs,omitempty"`
  Outputs []Param `yaml:"outputs,omitempty"`
  Run     *RunDeclaration `yaml:"run,omitempty"`
}
