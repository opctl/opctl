package opspec

import (
  "encoding/json"
)

func newJsonFormat(
) format {

  return &_jsonFormat{}

}

type _jsonFormat struct{}

func (this _jsonFormat) From(
in interface{},
) (out []byte, err error) {

  out, err = json.Marshal(in)
  return

}

func (this _jsonFormat)  To(
in []byte,
out interface{},
) (err error) {

  err = json.Unmarshal(in, out)

  return

}
