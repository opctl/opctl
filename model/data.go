package model

// Value represents a typed value
type Value struct {
	Dir    *string  `json:"dir,omitempty"`
	File   *string  `json:"file,omitempty"`
	Number *float64 `json:"number,omitempty"`
	Socket *string  `json:"socket,omitempty"`
	String *string  `json:"string,omitempty"`
}
