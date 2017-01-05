package model

type StringConstraints struct {
	Length   *StringLengthConstraint    `yaml:"length,omitempty"`
	Patterns []*StringPatternConstraint `yaml:"regex,omitempty"`
}

type StringLengthConstraint struct {
	Min         int    `yaml:"min,omitempty"`
	Max         int    `yaml:"max,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type StringPatternConstraint struct {
	Regex       string `yaml:"regex"`
	Description string `yaml:"description,omitempty"`
}

type NetSocketConstraints struct {
	PortNumber *PortNumberNetSocketConstraint `yaml:"port"`
}

type PortNumberNetSocketConstraint struct {
	Number      uint   `yaml:"number"`
	Description string `yaml:"description"`
}
