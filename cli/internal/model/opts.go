package model

type RunOpts struct {
	ArgFile string
	Args    []string
}

type NodeCreateOpts struct {
	// DataDir sets the path of dir used to store node data
	DataDir string
	// ListenAddress sets the HOST:PORT on which the node will listen
	ListenAddress string
}
