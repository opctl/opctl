package model

// Param represents a typed param of an op
type Param struct {
	Array       *ArrayParam   `json:"array,omitempty"`
	Boolean     *BooleanParam `json:"boolean,omitempty"`
	Description string        `json:"description,omitempty"`
	Dir         *DirParam     `json:"dir,omitempty"`
	File        *FileParam    `json:"file,omitempty"`
	Number      *NumberParam  `json:"number,omitempty"`
	Object      *ObjectParam  `json:"object,omitempty"`
	Socket      *SocketParam  `json:"socket,omitempty"`
	String      *StringParam  `json:"string,omitempty"`
}

// ArrayParam represents a parameter of type object
type ArrayParam struct {
	Constraints Constraints    `json:"constraints,omitempty"`
	Default     *[]interface{} `json:"default,omitempty"`
	// Deprecated: use Param.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// BooleanParam represents a parameter of type boolean
type BooleanParam struct {
	// Default is *bool instead of bool so we know if default was explicitly provided
	Default *bool `json:"default,omitempty"`
	// Deprecated: use Param.Description
	Description string `json:"description,omitempty"`
}

// DirParam represents a parameter of type directory
type DirParam struct {
	// Default is *string instead of string so we know if default was explicitly provided
	Default *string `json:"default,omitempty"`
	// Deprecated: use Param.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// FileParam represents a parameter of type file
type FileParam struct {
	// Default is *string instead of string so we know if default was explicitly provided
	Default *string `json:"default,omitempty"`
	// Deprecated: use Param.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// NumberParam represents a parameter of type number
type NumberParam struct {
	Constraints Constraints `json:"constraints,omitempty"`
	// Default is *float64 instead of float64 so we know if default was explicitly provided
	Default *float64 `json:"default,omitempty"`
	// Deprecated: use Param.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// ObjectParam represents a parameter of type object
type ObjectParam struct {
	Constraints Constraints             `json:"constraints,omitempty"`
	Default     *map[string]interface{} `json:"default,omitempty"`
	// Deprecated: use Param.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// SocketParam represents a parameter of type socket
type SocketParam struct {
	// Deprecated: use Param.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// StringParam represents a parameter of type string
type StringParam struct {
	Constraints Constraints `json:"constraints,omitempty"`
	// Default is *string instead of string so we know if default was explicitly provided
	Default *string `json:"default,omitempty"`
	// Deprecated: use Param.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

type Constraints map[string]interface{}
