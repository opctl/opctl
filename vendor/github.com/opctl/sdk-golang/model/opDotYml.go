package model

type OpDotYml struct {
	Description string            `yaml:"description"`
	Inputs      map[string]*Param `yaml:"inputs,omitempty"`
	Name        string            `yaml:"name"`
	Outputs     map[string]*Param `yaml:"outputs,omitempty"`
	Run         *SCG              `yaml:"run,omitempty"`
	Version     string            `yaml:"version,omitempty"`
}
