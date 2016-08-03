package models

func NewMountPoint(
name string,
path string,
) *MountPoint {

  return &MountPoint{
    Name:name,
    Path:path,
  }

}

type MountPoint struct {
  Name string `yaml:"name"`
  Path string `yaml:"path"`
}
