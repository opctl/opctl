package opspec

//go:generate counterfeiter -o ./fakeYamlCodec.go --fake-name fakeYamlCodec ./ yamlCodec

import (
  "gopkg.in/yaml.v2"
)

type yamlCodec interface {
  ToYaml(
  in interface{},
  ) (opFileBytes []byte, err error)

  FromYaml(
  in []byte,
  out interface{},
  ) (err error)
}

func newYamlCodec() yamlCodec {

  return &_yamlCodec{}

}

type _yamlCodec struct{}

func (this _yamlCodec) ToYaml(
in interface{},
) (opFileBytes []byte, err error) {

  opFileBytes, err = yaml.Marshal(in)

  return

}

func (this _yamlCodec)  FromYaml(
in []byte,
out interface{},
) (err error) {

  err = yaml.Unmarshal(in, out)

  return

}
