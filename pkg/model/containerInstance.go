package model

// entry in a container instances env; an env var
type ContainerInstanceEnvEntry struct {
  // name of the env var in the container
  Name  string
  // value of the env var in the container
  Value string
}

// entry in a container instances fs; a file/directory
type ContainerInstanceFsEntry struct {
  // reference to src of fs entry.
  // note: the use of an opaque string reference rather than an io.Reader, []byte, or path. Why? in order to remain fs-agnostic.
  SrcRef string
  // path of the file/directory in the container
  Path   string
}

// entry in a container instances network; a network socket
type ContainerInstanceNetEntry struct {
  Host        string `json:"host"`
  // aliases to give the network socket host in the container
  HostAliases []string `json:"hostAliases"`
  Port        uint `json:"port"`
}
