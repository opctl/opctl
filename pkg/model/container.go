package model

type Container struct {
  Cmd   []string `json:"cmd"`
  Env   []*ContainerEnvEntry `json:"env"`
  Fs    []*ContainerFsEntry `json:"fs"`
  Image string `json:"image"`
  Net   []*ContainerNetEntry `json:"net"`
  WorkDir string `json:"workDir"`
}

// entry in a containers env; an env var
type ContainerEnvEntry struct {
  Name  string `json:"name"`
  Value string `json:"value"`
}

// entry in a containers fs; a file/directory
type ContainerFsEntry struct {
  // reference to src of fs entry.
  // note: the use of an opaque string reference rather than an io.Reader, []byte, or path. Why? in order to remain fs-agnostic.
  SrcRef string `json:"srcRef"`
  // path of the file/directory in the container
  Path   string `json:"path"`
}

// entry in a containers network; a network socket
type ContainerNetEntry struct {
  Host        string `json:"host"`
  // aliases to give the network socket host in the container
  HostAliases []string `json:"hostAliases"`
  Port        uint     `json:"port"`
}
