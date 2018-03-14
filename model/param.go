package model

// Param represents a typed param of an op
type Param struct {
	Array   *ArrayParam   `yaml:"array,omitempty"`
	Boolean *BooleanParam `yaml:"boolean,omitempty"`
	Dir     *DirParam     `yaml:"dir,omitempty"`
	File    *FileParam    `yaml:"file,omitempty"`
	Number  *NumberParam  `yaml:"number,omitempty"`
	Object  *ObjectParam  `yaml:"object,omitempty"`
	Socket  *SocketParam  `yaml:"socket,omitempty"`
	String  *StringParam  `yaml:"string,omitempty"`
}

// ArrayParam represents a parameter of type object
type ArrayParam struct {
	Constraints *ArrayConstraints `yaml:"constraints,omitempty"`
	Default     []interface{}     `yaml:"default,omitempty"`
	Description string            `yaml:"description,omitempty"`
	IsSecret    bool              `yaml:"isSecret,omitempty"`
}

// BooleanParam represents a parameter of type boolean
type BooleanParam struct {
	// Default is *bool instead of bool so we know if default was explicitly provided
	Default     *bool  `yaml:"default,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// ArrayConstraints represents constraints of an ObjectParam
type ArrayConstraints struct {
	// json struct tags used for validating via gojsonschema
	Items           interface{}      `json:"items,omitempty" yaml:"items,omitempty"`
	AdditionalItems *TypeConstraints `json:"additionalItems,omitempty" yaml:"additionalItems,omitempty"`
	MaxItems        int              `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems        int              `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems     bool             `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
}

// DirParam represents a parameter of type directory
type DirParam struct {
	// Default is *string instead of string so we know if default was explicitly provided
	Default     *string `yaml:"default,omitempty"`
	Description string  `yaml:"description,omitempty"`
	IsSecret    bool    `yaml:"isSecret,omitempty"`
}

// FileParam represents a parameter of type file
type FileParam struct {
	// Default is *string instead of string so we know if default was explicitly provided
	Default     *string `yaml:"default,omitempty"`
	Description string  `yaml:"description,omitempty"`
	IsSecret    bool    `yaml:"isSecret,omitempty"`
}

// NumberParam represents a parameter of type number
type NumberParam struct {
	Constraints *NumberConstraints `yaml:"constraints,omitempty"`
	// Default is *float64 instead of float64 so we know if default was explicitly provided
	Default     *float64 `yaml:"default,omitempty"`
	Description string   `yaml:"description,omitempty"`
	IsSecret    bool     `yaml:"isSecret,omitempty"`
}

// NumberConstraints represents constraints of a NumberParam
type NumberConstraints struct {
	// json struct tags used for validating via gojsonschema
	AllOf      []*NumberConstraints `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf      []*NumberConstraints `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Enum       []float64            `json:"enum,omitempty" yaml:"enum,omitempty"`
	Format     string               `json:"format,omitempty" yaml:"format,omitempty"`
	Maximum    float64              `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	Minimum    float64              `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	MultipleOf float64              `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Not        *NumberConstraints   `json:"not,omitempty" yaml:"not,omitempty"`
	OneOf      []*NumberConstraints `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
}

// ObjectParam represents a parameter of type object
type ObjectParam struct {
	Constraints *ObjectConstraints     `yaml:"constraints,omitempty"`
	Default     map[string]interface{} `yaml:"default,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	IsSecret    bool                   `yaml:"isSecret,omitempty"`
}

// ObjectConstraints represents constraints of a ObjectParam
type ObjectConstraints struct {
	// json struct tags used for validating via gojsonschema
	AdditionalProperties *TypeConstraints            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	AllOf                []*ObjectConstraints        `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf                []*ObjectConstraints        `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Description          bool                        `json:"description,omitempty" yaml:"description,omitempty"`
	Enum                 []map[string]interface{}    `json:"enum,omitempty" yaml:"enum,omitempty"`
	MaxProperties        int                         `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        int                         `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Not                  *ObjectConstraints          `json:"not,omitempty" yaml:"not,omitempty"`
	OneOf                []*ObjectConstraints        `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	PatternProperties    map[string]*TypeConstraints `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	Properties           map[string]*TypeConstraints `json:"properties,omitempty" yaml:"properties,omitempty"`
	Required             []string                    `json:"required,omitempty" yaml:"required,omitempty"`
	Title                bool                        `json:"title,omitempty" yaml:"title,omitempty"`
	// Type must be interface{} because can be either string or []string
	Type      interface{} `json:"type,omitempty" yaml:"type,omitempty"`
	WriteOnly bool        `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
}

// TypeConstraints represents data type constraints
type TypeConstraints struct {
	// json struct tags used for validating via gojsonschema

	// override Enum from typed constraints
	Enum []interface{} `json:"enum,omitempty" yaml:"enum,omitempty"`

	// include all typed constraints
	ArrayConstraints
	NumberConstraints
	ObjectConstraints
	StringConstraints
}

// SocketParam represents a parameter of type socket
type SocketParam struct {
	Description string `yaml:"description,omitempty"`
	IsSecret    bool   `yaml:"isSecret,omitempty"`
}

// StringParam represents a parameter of type string
type StringParam struct {
	Constraints *StringConstraints `yaml:"constraints,omitempty"`
	// Default is *string instead of string so we know if default was explicitly provided
	Default     *string `yaml:"default,omitempty"`
	Description string  `yaml:"description,omitempty"`
	IsSecret    bool    `yaml:"isSecret,omitempty"`
}

// StringConstraints represents constraints of a StringParam
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
