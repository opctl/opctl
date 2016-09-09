package models

type OpManifest struct {
  Manifest `yaml:",inline"`
  Inputs []Param `yaml:"inputs,omitempty"`
  Run    *RunDeclaration `yaml:"run,omitempty"`
}
