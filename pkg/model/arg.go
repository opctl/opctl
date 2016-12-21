package model

// arg for an op
type Arg struct {
  NetSocket *NetSocketArg `json:"netSocket"`
  // reference to a file of a fs
  File      string `json:"file"`
  // reference to a dir of a fs
  Dir       string `json:"dir"`
  String    string `json:"string"`
}

// network socket arg for an op
type NetSocketArg struct {
  Host string `json:"host"`
  Port uint `json:"port"`
}
