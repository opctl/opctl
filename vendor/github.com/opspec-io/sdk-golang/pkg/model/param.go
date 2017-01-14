package model

// Parameter of an op
type Param struct {
	Dir    *DirParam    `yaml:"dir,omitempty"`
	File   *FileParam   `yaml:"file,omitempty"`
	Socket *SocketParam `yaml:"socket,omitempty"`
	String *StringParam `yaml:"string,omitempty"`
}

// Directory parameter of an op
type DirParam struct {
	Description string `yaml:"description"`
	IsSecret    bool   `yaml:"isSecret"`
	Name        string `yaml:"name"`
}

// File parameter of an op
type FileParam struct {
	Description string `yaml:"description"`
	IsSecret    bool   `yaml:"isSecret"`
	Name        string `yaml:"name"`
}

// Socket parameter of an op
type SocketParam struct {
	Description string `yaml:"description"`
	IsSecret    bool   `yaml:"isSecret"`
	Name        string `yaml:"name"`
}

// String parameter of an op
type StringParam struct {
	Constraints *StringConstraints `yaml:"constraints"`
	Default     string             `yaml:"default,omitempty"`
	Description string             `yaml:"description"`
	IsSecret    bool               `yaml:"isSecret"`
	Name        string             `yaml:"name"`
}
