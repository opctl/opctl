package core

import (
  "gopkg.in/yaml.v2"
)

type yamlCodec interface {
  toYaml(
  in interface{},
  ) (pipelineFileBytes []byte, err error)

  fromYaml(
  in []byte,
  out interface{},
  ) (err error)
}

type _yamlCodec struct{}

func (this _yamlCodec) toYaml(
in interface{},
) (pipelineFileBytes []byte, err error) {

  pipelineFileBytes, err= yaml.Marshal(in)

  return

}

func (this _yamlCodec)  fromYaml(
in []byte,
out interface{},
) (err error) {

  err= yaml.Unmarshal(in, out)

  return

}