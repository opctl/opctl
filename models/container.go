package models

func NewContainer(
image string,
mounts []MountPoint,
ports []PortBinding,
process Process,
) *Container {

  return &Container{
    Image:image,
    Mounts:mounts,
    Ports:ports,
    Process:process,
  }

}

type Container struct {
  Image   string `yaml:"image"`
  Mounts  []MountPoint `yaml:"mounts"`
  Ports   []PortBinding `yaml:"ports"`
  Process Process `yaml:"process"`
}
