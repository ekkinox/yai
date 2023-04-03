package runner

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type RunnerOutput struct {
	error error
}

func (o RunnerOutput) GetError() error {
	return o.error
}

type Runner struct {
	dummy string
}

func NewRunner() *Runner {
	return &Runner{"run"}
}

func (r *Runner) RunCommand(input string) tea.Cmd {
	c := exec.Command("bash", "-c", input)
	return tea.ExecProcess(c, func(error error) tea.Msg {
		return RunnerOutput{error}
	})
}
