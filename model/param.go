package model

// Parameter of an op
type Param struct {
	Dir    *DirParam    `yaml:"dir,omitempty"`
	File   *FileParam   `yaml:"file,omitempty"`
	Number *NumberParam `yaml:"number,omitempty"`
	Socket *SocketParam `yaml:"socket,omitempty"`
	String *StringParam `yaml:"string,omitempty"`
}

// Number parameter
type NumberParam struct {
	Constraints *NumberConstraints `yaml:"constraints,omitempty"`
	Default     *float64           `yaml:"default,omitempty"`
	Description string             `yaml:"description,omitempty"`
	IsSecret    bool               `yaml:"isSecret,omitempty"`
}

// Number parameter constraints
type NumberConstraints struct {
	// json struct tags used for validating via gojsonschema
	AllOf      []*NumberConstraints `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf      []*NumberConstraints `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Enum       []float64            `json:"enum,omitempty" yaml:"enum,omitempty"`
	Format     string               `json:"format,omitempty" yaml:"format,omitempty"`
	Maximum    float64              `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	Minimum    float64              `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	MultipleOf float64              `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Integer    bool                 `json:"integer,omitempty" yaml:"integer,omitempty"`
	Not        *NumberConstraints   `json:"not,omitempty" yaml:"not,omitempty"`
	OneOf      []*NumberConstraints `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
}

// Directory parameter
type DirParam struct {
	Default     *string `yaml:"default,omitempty"`
	Description string  `yaml:"description,omitempty"`
	IsSecret    bool    `yaml:"isSecret,omitempty"`
}

// File parameter
type FileParam struct {
	Default     *string `yaml:"default,omitempty"`
	Description string  `yaml:"description,omitempty"`
	IsSecret    bool    `yaml:"isSecret,omitempty"`
}

// Socket parameter
type SocketParam struct {
	Description string `yaml:"description,omitempty"`
	IsSecret    bool   `yaml:"isSecret,omitempty"`
}

// String parameter
type StringParam struct {
	Constraints *StringConstraints `yaml:"constraints,omitempty"`
	Default     *string            `yaml:"default,omitempty"`
	Description string             `yaml:"description,omitempty"`
	IsSecret    bool               `yaml:"isSecret,omitempty"`
}

// String parameter constraints
type StringConstraints struct {
	// json struct tags used for validating via gojsonschema
	AllOf     []*StringConstraints `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf     []*StringConstraints `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Enum      []string             `json:"enum,omitempty" yaml:"enum,omitempty"`
	Format    string               `json:"format,omitempty" yaml:"format,omitempty"`
	MaxLength int                  `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength int                  `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Not       *StringConstraints   `json:"not,omitempty" yaml:"not,omitempty"`
	OneOf     []*StringConstraints `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Pattern   string               `json:"pattern,omitempty" yaml:"pattern,omitempty"`
}
