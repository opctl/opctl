package model

// Value represents a typed value
type Value struct {
	Array  []interface{}          `json:"array,omitempty"`
	Dir    *string                `json:"dir,omitempty"`
	File   *string                `json:"file,omitempty"`
	Number *float64               `json:"number,omitempty"`
	Object map[string]interface{} `json:"object,omitempty"`
	Socket *string                `json:"socket,omitempty"`
	String *string                `json:"string,omitempty"`
}
