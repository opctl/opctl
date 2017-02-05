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
	Inputs   map[string]*Param `yaml:"inputs,omitempty"`
	Outputs  map[string]*Param `yaml:"outputs,omitempty"`
	Run      *Scg              `yaml:"run,omitempty"`
}
