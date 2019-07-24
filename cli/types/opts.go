package types

type RunOpts struct {
	ArgFile string
	Args    []string
}

type NodeCreateOpts struct {
	// DataDir configures path of dir used to store node data
	DataDir *string
}
