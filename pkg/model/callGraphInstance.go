package model

type CallGraphInstance struct {
  CallGraphInstanceId string
  Container           *ContainerCallInstance
  Op                  *OpCallInstance
  Parallel            *ParallelCallInstance
  Serial              *SerialCallInstance
}

// instance of a container call
type ContainerCallInstance struct {
  Cmd     string
  Env     []*ContainerInstanceEnvEntry
  Fs      []*ContainerInstanceFsEntry
  Image   string
  Net     []*ContainerInstanceNetEntry
  WorkDir string
}

// instance of an op call
type OpCallInstance struct {
  Ref     string
  // binds inputs of referenced op to in scope variables
  Args    map[string]string
  // binds in scope variables to outputs of referenced op
  Results map[string]string
}

// instance of a parallel call
type ParallelCallInstance []CallGraphInstance

// instance of a serial call
type SerialCallInstance []CallGraphInstance
