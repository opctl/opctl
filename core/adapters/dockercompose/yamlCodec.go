package dockercompose

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

type yamlCodecImpl struct{}

func (yc yamlCodecImpl) toYaml(
in interface{},
) (pipelineFileBytes []byte, err error) {

  pipelineFileBytes, err= yaml.Marshal(in)

  return

}

func (yc yamlCodecImpl)  fromYaml(
in []byte,
out interface{},
) (err error) {

  err= yaml.Unmarshal(in, out)

  return

}