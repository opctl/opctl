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
  Env     []*EnvEntryInstance
  Fs      []*FsEntryInstance
  Image   string
  Net     []*NetEntryInstance
  WorkDir string
}

// instance of an entry in an env; an env var
type EnvEntryInstance struct {
  // name of env var
  Name  string
  // value of env var
  Value string
}

// instance of an entry in a fs; a file/directory
type FsEntryInstance struct {
  // id of volume
  VolumeId string
  // path in container
  Path     string
}

// instance of an entry in a network; a host
type NetEntryInstance struct {
  // id of network
  NetworkId string
  // hostname of container on network
  Hostname  string
}

// instance of an op call
type OpCallInstance struct {
  Ref     string
  // binds inputs of referenced op to in scope variables
  Args    *map[string]string
  // binds in scope variables to outputs of referenced op
  Results *map[string]string
}

// instance of a parallel call
type ParallelCallInstance []CallGraphInstance

// instance of a serial call
type SerialCallInstance []CallGraphInstance
