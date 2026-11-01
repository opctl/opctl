// Package sdk builds a Go op's op graph without running its calls.
package sdk

import (
	"fmt"
	"strings"
)

// Call is a call in a lazily built op graph.
type Call interface {
	call()
	render() (*callSpec, error)
}

// ContainerCall builds a container call. Its methods return independent copies.
type ContainerCall struct {
	container containerSpec
	err       error
}
type containerSpec struct {
	Image   imageSpec         `json:"image"`
	Cmd     []string          `json:"cmd,omitempty"`
	EnvVars map[string]string `json:"envVars,omitempty"`
	WorkDir string            `json:"workDir,omitempty"`
	Dirs    map[string]string `json:"dirs,omitempty"`
	Files   map[string]string `json:"files,omitempty"`
	Sockets map[string]string `json:"sockets,omitempty"`
}
type imageSpec struct {
	Ref string `json:"ref"`
}

// Container adds a container call using a literal image ref.
func Container(image string) ContainerCall {
	return ContainerCall{container: containerSpec{Image: imageSpec{Ref: escapeLiteral(image)}}}
}

// Cmd sets the container's literal command and arguments.
func (c ContainerCall) Cmd(args ...string) ContainerCall {
	c.container.Cmd = make([]string, len(args))
	for i, arg := range args {
		c.container.Cmd[i] = escapeLiteral(arg)
	}
	return c
}

// Env sets a literal environment variable.
func (c ContainerCall) Env(name, value string) ContainerCall {
	return c.binding(name, escapeLiteral(value), "env")
}

// EnvFrom binds an input to an environment variable.
func (c ContainerCall) EnvFrom(name string, value Handle) ContainerCall {
	return c.bindHandle(name, value, "env")
}

// Dir mounts a directory input at a container path.
func (c ContainerCall) Dir(path string, value Dir) ContainerCall {
	return c.bindHandle(path, value, "dir")
}

// File mounts a file input at a container path.
func (c ContainerCall) File(path string, value File) ContainerCall {
	return c.bindHandle(path, value, "file")
}

// Socket binds a socket input to a container address.
func (c ContainerCall) Socket(address string, value Socket) ContainerCall {
	return c.bindHandle(address, value, "socket")
}

// WorkDir sets the container's literal working directory.
func (c ContainerCall) WorkDir(path string) ContainerCall {
	c.container.WorkDir = escapeLiteral(path)
	return c
}
func (c ContainerCall) bindHandle(name string, value Handle, kind string) (result ContainerCall) {
	result = c
	defer func() {
		if failure := recover(); failure != nil {
			if err, ok := failure.(inputError); ok {
				if result.err == nil {
					result.err = err.error
				}
			} else {
				panic(failure)
			}
		}
	}()
	if value == nil {
		if c.err == nil {
			c.err = fmt.Errorf("handle is not bound to an input")
		}
		return c
	}
	ref := value.handle()
	return c.binding(name, ref, kind)
}
func (c ContainerCall) binding(name, value string, kind string) ContainerCall {
	if strings.Contains(name, "$") {
		if c.err == nil {
			c.err = fmt.Errorf("binding name %q must not contain $", name)
		}
		return c
	}
	var existing map[string]string
	switch kind {
	case "env":
		existing = c.container.EnvVars
	case "dir":
		existing = c.container.Dirs
	case "file":
		existing = c.container.Files
	case "socket":
		existing = c.container.Sockets
	}
	values := make(map[string]string, len(existing)+1)
	for k, v := range existing {
		values[k] = v
	}
	values[name] = value
	switch kind {
	case "env":
		c.container.EnvVars = values
	case "dir":
		c.container.Dirs = values
	case "file":
		c.container.Files = values
	case "socket":
		c.container.Sockets = values
	}
	return c
}
func escapeLiteral(value string) string {
	var out strings.Builder
	slashes := 0
	for _, r := range value {
		if r == '\\' {
			slashes++
			continue
		}
		if r == '$' {
			out.WriteString(strings.Repeat("\\", 2*slashes+1))
		} else {
			out.WriteString(strings.Repeat("\\", slashes))
		}
		out.WriteRune(r)
		slashes = 0
	}
	out.WriteString(strings.Repeat("\\", slashes))
	return out.String()
}

func (ContainerCall) call() {}

func (c ContainerCall) render() (*callSpec, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.container.Image.Ref == "" {
		return nil, fmt.Errorf("Go op plan must return a container call with an image")
	}
	return &callSpec{Container: &c.container}, nil
}
