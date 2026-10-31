// Package sdk builds a Go op's op graph without running its calls.
package sdk

// Call is a call in a lazily built op graph.
type Call struct {
	container containerSpec
}

type containerSpec struct {
	Image   imageSpec         `json:"image"`
	Cmd     []string          `json:"cmd,omitempty"`
	EnvVars map[string]string `json:"envVars,omitempty"`
	WorkDir string            `json:"workDir,omitempty"`
}

type imageSpec struct {
	Ref string `json:"ref"`
}

// Container adds a container call using a literal image ref.
func Container(image string) Call {
	return Call{container: containerSpec{Image: imageSpec{Ref: image}}}
}

// Cmd sets the container's literal command and arguments.
func (c Call) Cmd(args ...string) Call {
	c.container.Cmd = append([]string(nil), args...)
	return c
}

// Env sets a literal environment variable.
func (c Call) Env(name, value string) Call {
	env := make(map[string]string, len(c.container.EnvVars)+1)
	for key, existing := range c.container.EnvVars {
		env[key] = existing
	}
	env[name] = value
	c.container.EnvVars = env
	return c
}

// WorkDir sets the container's literal working directory.
func (c Call) WorkDir(path string) Call {
	c.container.WorkDir = path
	return c
}
