package context

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ekkinox/yo/runner"
	"github.com/mitchellh/go-homedir"
)

const app_name = "Yo"

type Context struct {
	appName         string
	operatingSystem string
	distribution    string
	shell           string
	homeDirectory   string
	username        string
	editor          string
	configFile      string
}

func (o *Context) GetAppName() string {
	return o.appName
}

func (o *Context) GetOperatingSystem() string {
	return o.operatingSystem
}

func (o *Context) GetDistribution() string {
	return o.distribution
}

func (o *Context) GetShell() string {
	return o.shell
}

func (o *Context) GetHomeDirectory() string {
	return o.homeDirectory
}

func (o *Context) GetUsername() string {
	return o.username
}

func (o *Context) GetEditor() string {
	return o.editor
}

func (o *Context) GetConfigFile() string {
	return o.configFile
}

type ContextAnalyzer struct {
	config string
	runner *runner.Runner
}

func NewContextAnalyzer() *ContextAnalyzer {
	return &ContextAnalyzer{
		config: "config",
		runner: runner.NewRunner(),
	}
}

func (a *ContextAnalyzer) Analyse() *Context {
	return &Context{
		appName:         app_name,
		operatingSystem: a.GetOperatingSystem().String(),
		distribution:    a.GetDistribution(),
		shell:           a.GetShell(),
		homeDirectory:   a.GetHomeDirectory(),
		username:        a.GetUsername(),
		editor:          a.GetEditor(),
		configFile:      a.GetConfigFile(),
	}
}

func (a *ContextAnalyzer) GetOperatingSystem() OperatingSystem {
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

func (a *ContextAnalyzer) GetDistribution() string {
	dist, err := a.runner.Run("lsb_release", "-sd")
	if err != nil {
		return ""
	}

	return strings.Trim(strings.Trim(dist, "\n"), "\"")
}

func (a *ContextAnalyzer) GetShell() string {
	shell, err := a.runner.Run("echo", os.Getenv("SHELL"))
	if err != nil {
		return ""
	}

	split := strings.Split(strings.Trim(strings.Trim(shell, "\n"), "\""), "/")

	return split[len(split)-1]

}

func (a *ContextAnalyzer) GetHomeDirectory() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return homeDir
}

func (a *ContextAnalyzer) GetUsername() string {
	name, err := a.runner.Run("echo", os.Getenv("USER"))
	if err != nil {
		return ""
	}

	return strings.Trim(name, "\n")
}

func (a *ContextAnalyzer) GetEditor() string {
	name, err := a.runner.Run("echo", os.Getenv("EDITOR"))
	if err != nil {
		return "nano"
	}

	return strings.Trim(name, "\n")
}

func (a *ContextAnalyzer) GetConfigFile() string {
	return fmt.Sprintf(
		"%s/.config/%s.json",
		a.GetHomeDirectory(),
		strings.ToLower(app_name),
	)
}
