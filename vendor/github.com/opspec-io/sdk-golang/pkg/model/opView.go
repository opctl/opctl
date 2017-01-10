package model

type OpView struct {
	Description string
	Inputs      []*Param
	Name        string
	Outputs     []*Param
	Run         *Scg
	Version     string
}
