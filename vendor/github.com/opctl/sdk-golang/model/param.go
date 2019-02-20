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
	Constraints Constraints   `yaml:"constraints,omitempty"`
	Default     []interface{} `yaml:"default,omitempty"`
	Description string        `yaml:"description,omitempty"`
	IsSecret    bool          `yaml:"isSecret,omitempty"`
}

// BooleanParam represents a parameter of type boolean
type BooleanParam struct {
	// Default is *bool instead of bool so we know if default was explicitly provided
	Default     *bool  `yaml:"default,omitempty"`
	Description string `yaml:"description,omitempty"`
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
	Constraints Constraints `yaml:"constraints,omitempty"`
	// Default is *float64 instead of float64 so we know if default was explicitly provided
	Default     *float64 `yaml:"default,omitempty"`
	Description string   `yaml:"description,omitempty"`
	IsSecret    bool     `yaml:"isSecret,omitempty"`
}

// ObjectParam represents a parameter of type object
type ObjectParam struct {
	Constraints Constraints            `yaml:"constraints,omitempty"`
	Default     map[string]interface{} `yaml:"default,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	IsSecret    bool                   `yaml:"isSecret,omitempty"`
}

// SocketParam represents a parameter of type socket
type SocketParam struct {
	Description string `yaml:"description,omitempty"`
	IsSecret    bool   `yaml:"isSecret,omitempty"`
}

// StringParam represents a parameter of type string
type StringParam struct {
	Constraints Constraints `yaml:"constraints,omitempty"`
	// Default is *string instead of string so we know if default was explicitly provided
	Default     *string `yaml:"default,omitempty"`
	Description string  `yaml:"description,omitempty"`
	IsSecret    bool    `yaml:"isSecret,omitempty"`
}

type Constraints map[string]interface{}
