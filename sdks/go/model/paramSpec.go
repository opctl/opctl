package model

// ParamSpec represents a typed ParamSpec of an op
type ParamSpec struct {
	Array       *ArrayParamSpec   `json:"array,omitempty"`
	Boolean     *BooleanParamSpec `json:"boolean,omitempty"`
	Description string            `json:"description,omitempty"`
	Dir         *DirParamSpec     `json:"dir,omitempty"`
	File        *FileParamSpec    `json:"file,omitempty"`
	Number      *NumberParamSpec  `json:"number,omitempty"`
	Object      *ObjectParamSpec  `json:"object,omitempty"`
	Socket      *SocketParamSpec  `json:"socket,omitempty"`
	String      *StringParamSpec  `json:"string,omitempty"`
}

// ArrayParamSpec represents a Parameter of type object
type ArrayParamSpec struct {
	Constraints Constraints `json:"constraints,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	// Deprecated: use ParamSpec.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// BooleanParamSpec represents a Parameter of type boolean
type BooleanParamSpec struct {
	Default interface{} `json:"default,omitempty"`
	// Deprecated: use ParamSpec.Description
	Description string `json:"description,omitempty"`
}

// DirParamSpec represents a Parameter of type directory
type DirParamSpec struct {
	Default interface{} `json:"default,omitempty"`
	// Deprecated: use ParamSpec.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// FileParamSpec represents a Parameter of type file
type FileParamSpec struct {
	Default interface{} `json:"default,omitempty"`
	// Deprecated: use ParamSpec.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// NumberParamSpec represents a Parameter of type number
type NumberParamSpec struct {
	Constraints Constraints `json:"constraints,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	// Deprecated: use ParamSpec.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// ObjectParamSpec represents a Parameter of type object
type ObjectParamSpec struct {
	Constraints Constraints `json:"constraints,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	// Deprecated: use ParamSpec.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// SocketParamSpec represents a Parameter of type socket
type SocketParamSpec struct {
	// Deprecated: use ParamSpec.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

// StringParamSpec represents a Parameter of type string
type StringParamSpec struct {
	Constraints Constraints `json:"constraints,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	// Deprecated: use ParamSpec.Description
	Description string `json:"description,omitempty"`
	IsSecret    bool   `json:"isSecret,omitempty"`
}

type Constraints map[string]interface{}
