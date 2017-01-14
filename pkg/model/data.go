package model

// typed data
type Data struct {
  NetSocket *NetSocket `json:"netSocket"`
  File      string `json:"file"`
  Dir       string `json:"dir"`
  String    string `json:"string"`
}

// network socket
type NetSocket struct {
  Host string `json:"host"`
  Port uint   `json:"port"`
}
