package context

import (
	"os"
	"runtime"
	"strings"

	"github.com/ekkinox/yo/runner"

	"github.com/mitchellh/go-homedir"
)

type ContextOutput struct {
	operatingSystem string
	distribution    string
	shell           string
	homeDirectory   string
	username        string
	editor          string
}

func (o *ContextOutput) GetOperatingSystem() string {
	return o.operatingSystem
}

func (o *ContextOutput) GetDistribution() string {
	return o.distribution
}

func (o *ContextOutput) GetShell() string {
	return o.shell
}

func (o *ContextOutput) GetHomeDirectory() string {
	return o.homeDirectory
}

func (o *ContextOutput) GetUsername() string {
	return o.username
}

func (o *ContextOutput) GetEditor() string {
	return o.editor
}

type Context struct {
	config string
	runner *runner.Runner
}

func NewContext() *Context {
	return &Context{
		config: "config",
		runner: runner.NewRunner(),
	}
}

func (c *Context) Analyse() ContextOutput {
	return ContextOutput{
		operatingSystem: c.GetOperatingSystem().String(),
		distribution:    c.GetDistribution(),
		shell:           c.GetShell(),
		homeDirectory:   c.GetHomeDirectory(),
		username:        c.GetUsername(),
		editor:          c.GetEditor(),
	}
}

func (c *Context) GetOperatingSystem() OperatingSystem {
	switch runtime.GOOS {
	case "linux":
		return LinuxOperatingSystem
	case "darwin":
		return MacOperatingSystem
	case "windows":
		return WindowsOperatingSystem
	default:
		return UnknownOperatingSystem
	}
}

func (c *Context) GetDistribution() string {
	dist, err := c.runner.Run("lsb_release", "-sd")
	if err != nil {
		return ""
	}

	return strings.Trim(strings.Trim(dist, "\n"), "\"")
}

func (c *Context) GetShell() string {
	shell, err := c.runner.Run("echo", os.Getenv("SHELL"))
	if err != nil {
		return ""
	}

	split := strings.Split(strings.Trim(strings.Trim(shell, "\n"), "\""), "/")

	return split[len(split)-1]

}

func (c *Context) GetHomeDirectory() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return homeDir
}

func (c *Context) GetUsername() string {
	name, err := c.runner.Run("echo", os.Getenv("USER"))
	if err != nil {
		return ""
	}

	return strings.Trim(name, "\n")
}

func (c *Context) GetEditor() string {
	name, err := c.runner.Run("echo", os.Getenv("EDITOR"))
	if err != nil {
		return "vim"
	}

	return strings.Trim(name, "\n")
}
