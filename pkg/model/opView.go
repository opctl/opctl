package model

type OpView struct {
	Description string
	Inputs      map[string]*Param
	Name        string
	Outputs     map[string]*Param
	Run         *Scg
	Version     string
}
