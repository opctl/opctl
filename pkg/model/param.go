package model

type Param struct {
  // backwards compatibility w/ opspec v 0.1.2
  V0_1_2Param `yaml:",inline,omitempty"`

  Dir       *DirParam `yaml:",omitempty"`
  File      *FileParam `yaml:",omitempty"`
  NetSocket *NetSocketParam `yaml:",omitempty"`
  String    *StringParam `yaml:",omitempty"`
}

type V0_1_2Param struct {
  Name        string `yaml:"name,omitempty"`
  Description string `yaml:"description,omitempty"`
  IsSecret    bool `yaml:"isSecret,omitempty"`
  Default     string `yaml:"default,omitempty"`
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
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  IsSecret    bool `yaml:"isSecret"`
  Default     string `yaml:"default,omitempty"`
}
