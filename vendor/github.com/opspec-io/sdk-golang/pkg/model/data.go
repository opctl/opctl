package model

// typed data
type Data struct {
  NetSocket *NetSocketData `json:"netSocket"`
  // reference to a file of a fs
  File      string `json:"file"`
  // reference to a dir of a fs
  Dir       string `json:"dir"`
  String    string `json:"string"`
}

// network socket
type NetSocketData struct {
  Host string `json:"host"`
  Port uint   `json:"port"`
}
