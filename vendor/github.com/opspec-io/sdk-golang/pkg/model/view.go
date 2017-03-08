package model

type PackageManifestView struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Version     string            `yaml:"version,omitempty"`
	Inputs      map[string]*Param `yaml:"inputs,omitempty"`
	Outputs     map[string]*Param `yaml:"outputs,omitempty"`
	Run         *Scg              `yaml:"run,omitempty"`
}

type PackageView struct {
	Description string
	Inputs      map[string]*Param
	Name        string
	Outputs     map[string]*Param
	Run         *Scg
	Version     string
}
