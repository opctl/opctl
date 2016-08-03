package models

func NewPortBinding(
endPort int,
protocol string,
startPort int,
) *PortBinding {

  return &PortBinding{
    EndPort:endPort,
    Protocol:protocol,
    StartPort:startPort,
  }

}

type PortBinding struct {
  EndPort int `yaml:"path"`
  Protocol string `yaml:"procotol"`
  StartPort int `yaml:"name"`
}
