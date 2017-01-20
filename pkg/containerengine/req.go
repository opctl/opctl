package containerengine

type StartContainerReq struct {
	Cmd         []string
	Dirs        map[string]string
	Env         map[string]string
	Files       map[string]string
	Image       string
	Net         map[string]string
	WorkDir     string
	ContainerId string
	OpGraphId   string
}
