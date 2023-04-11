package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/ekkinox/yo/ai"
	"github.com/ekkinox/yo/config"
)

type CliState struct {
	configuring bool
	querying    bool
	confirming  bool
	executing   bool
	args        string
	buffer      string
	command     string
	error       error
}

type CliDimensions struct {
	width  int
	height int
}

type CliComponents struct {
	prompt   *Prompt
	renderer *Renderer
	spinner  *Spinner
}

type Cli struct {
	state      CliState
	dimensions CliDimensions
	components CliComponents
	config     *config.Config
	engine     *ai.Engine
}

func NewCli(args string) *Cli {
	return &Cli{
		state: CliState{
			configuring: false,
			querying:    false,
			confirming:  false,
			executing:   false,
			args:        args,
			buffer:      "",
			command:     "",
			error:       nil,
		},
		dimensions: CliDimensions{
			150,
			150,
		},
		components: CliComponents{
			prompt: NewPrompt(ExecPromptMode),
			renderer: NewRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(150),
			),
			spinner: NewSpinner(),
		},
	}
}

func (c *Cli) Init() tea.Cmd {
	return tea.Sequence(
		tea.Println("cli mode "+c.state.args),
		tea.Quit,
	)
}

func (c *Cli) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmds []tea.Cmd
	)

	return c, tea.Batch(cmds...)
}

func (c *Cli) View() string {
	return ""
}
