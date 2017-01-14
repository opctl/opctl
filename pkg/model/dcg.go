package model

// dynamic call graph; see https://en.wikipedia.org/wiki/Call_graph

type DcgContainer struct {
  Cmd     []string `json:"cmd"`
  // format: containerPath => hostPath
  Dirs    map[string]string `json:"dirs"`
  Env     []*DcgContainerEnvEntry `json:"env"`
  // format: containerPath => hostPath
  Files   map[string]string `json:"files"`
  Image   string `json:"image"`
  Net     []*DcgContainerNetEntry `json:"net"`
  WorkDir string `json:"workDir"`
}

// entry in a containers env; an env var
type DcgContainerEnvEntry struct {
  Name  string `json:"name"`
  Value string `json:"value"`
}

// entry in a containers network; a network socket
type DcgContainerNetEntry struct {
  Host        string `json:"host"`
  // aliases to give the network socket host in the container
  HostAliases []string `json:"hostAliases"`
  Port        uint     `json:"port"`
}
