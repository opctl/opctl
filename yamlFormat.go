package opspec

import (
  "gopkg.in/yaml.v2"
)

func newYamlFormat(
) format {

  return &_yamlFormat{}

}

type _yamlFormat struct{}

func (this _yamlFormat) From(
in interface{},
) (out []byte, err error) {

  out, err = yaml.Marshal(in)
  return

}

func (this _yamlFormat)  To(
in []byte,
out interface{},
) (err error) {

  err = yaml.Unmarshal(in, out)
  return

}
