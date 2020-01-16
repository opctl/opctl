package model

type OpFile struct {
	Description string            `json:"description"`
	Inputs      map[string]*Param `json:"inputs,omitempty"`
	Name        string            `json:"name"`
	Outputs     map[string]*Param `json:"outputs,omitempty"`
	Run         *SCG              `json:"run,omitempty"`
	Version     string            `json:"version,omitempty"`
}
