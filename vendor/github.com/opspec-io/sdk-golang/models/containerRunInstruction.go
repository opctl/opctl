package models

func NewContainerRunInstruction(
container *Container,
) *ContainerRunInstruction {

  return &ContainerRunInstruction{
    Container:container,
  }

}

type ContainerRunInstruction struct {
  Container *Container `yaml:"container"`
}
