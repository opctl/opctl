package model

type Param struct {
  Dir       *DirParam `yaml:",omitempty"`
  File      *FileParam `yaml:",omitempty"`
  NetSocket *NetSocketParam `yaml:",omitempty"`
  String    *StringParam `yaml:",omitempty"`
}

type DirParam struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
}

type FileParam struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
}

type NetSocketParam struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
}

type StringParam struct {
  Default     string `yaml:"default,omitempty"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
  MaxLength   int `yaml:"maxLength,omitempty"`
  MinLength   int `yaml:"minLength,omitempty"`
  Name        string `yaml:"name"`
  Pattern     string `yaml:"pattern,omitempty"`
}
