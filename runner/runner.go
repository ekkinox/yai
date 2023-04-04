package runner

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type RunnerOutput struct {
	Error error
}

func (o RunnerOutput) GetError() error {
	return o.Error
}

type Runner struct {
	config string
}

func NewRunner() *Runner {
	return &Runner{"config"}
}

func (r *Runner) RunInteractive(input string) tea.Cmd {

	time.Sleep(time.Microsecond * 100)
	c := exec.Command("bash", "-c", input)

	return tea.ExecProcess(c, func(error error) tea.Msg {
		time.Sleep(time.Microsecond * 100)

		return RunnerOutput{error}
	})
}

func (r *Runner) Run(cmd string, arg ...string) (string, error) {
	out, err := exec.Command(cmd, arg...).Output()
	if err != nil {
		message := fmt.Sprintf("error: %v", err)
		log.Println(message)

		return message, err
	}

	return string(out), nil
}
