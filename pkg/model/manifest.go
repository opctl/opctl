package model

type Manifest struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version,omitempty"`
}

type CollectionManifest struct {
	Manifest `yaml:",inline"`
}

type OpManifest struct {
	Manifest `yaml:",inline"`
	Inputs   []*Param   `yaml:"inputs,omitempty"`
	Outputs  []*Param   `yaml:"outputs,omitempty"`
	Run      *CallGraph `yaml:"run,omitempty"`
}
