package core

import (
  "gopkg.in/yaml.v2"
)

type yamlCodec interface {
  toYaml(
  in interface{},
  ) (opFileBytes []byte, err error)

  fromYaml(
  in []byte,
  out interface{},
  ) (err error)
}

func newYamlCodec() yamlCodec {

  return &_yamlCodec{}

}

type _yamlCodec struct{}

func (this _yamlCodec) toYaml(
in interface{},
) (opFileBytes []byte, err error) {

  opFileBytes, err = yaml.Marshal(in)

  return

}

func (this _yamlCodec)  fromYaml(
in []byte,
out interface{},
) (err error) {

  err = yaml.Unmarshal(in, out)

  return

}