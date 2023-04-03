package tui

import (
	"errors"
	"fmt"
	"github.com/ekkinox/yo/engine"
	"github.com/ekkinox/yo/history"
	"github.com/ekkinox/yo/runner"
	"github.com/ekkinox/yo/tui/components"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type TuiState struct {
	running bool
	buffer  string
	error   error
}

type TuiComponents struct {
	prompt   textinput.Model
	renderer *components.Renderer
}

type Tui struct {
	state      TuiState
	components TuiComponents
	history    *history.History
	engine     *engine.Engine
	runner     *runner.Runner
}

func NewTui() *Tui {
	i := textinput.New()
	i.Placeholder = "Enter something"
	i.Focus()

	return &Tui{
		state: TuiState{
			running: false,
			buffer:  "",
			error:   nil,
		},
		components: TuiComponents{
			prompt: components.NewPrompt(),
			renderer: components.NewRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(100),
			),
		},
		history: history.NewHistory(),
		engine:  engine.NewEngine(),
		runner:  runner.NewRunner(),
	}
}

func (t *Tui) Init() tea.Cmd {
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmds      []tea.Cmd
		promptCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return t, tea.Quit
		case tea.KeyEnter:
			input := t.components.prompt.Value()
			t.history.Add(input)
			t.components.prompt, promptCmd = t.components.prompt.Update(msg)
			cmds = append(
				cmds,
				promptCmd,
				tea.Printf("\n> %s\n", input),
				t.startEngine(input),
				t.awaitEngine(),
			)
		case tea.KeyCtrlR:
			input := t.components.prompt.Value()
			t.history.Add(input)
			t.components.prompt, promptCmd = t.components.prompt.Update(msg)
			cmds = append(
				cmds,
				promptCmd,
				t.runner.RunCommand(input),
			)
		default:
			t.components.prompt, promptCmd = t.components.prompt.Update(msg)
			cmds = append(
				cmds,
				promptCmd,
			)
		}
	case engine.EngineOutput:
		if msg.IsLast() {
			t.components.prompt.Focus()
			t.components.prompt, promptCmd = t.components.prompt.Update(msg)
			return t, tea.Sequence(
				promptCmd,
				tea.Printf(t.components.renderer.Render(t.state.buffer)),
			)
		} else {
			return t, t.awaitEngine()
		}
	case runner.RunnerOutput:
		if msg.GetError() != nil {
			t.state.error = msg.GetError()
			return t, tea.Quit
		}
		t.components.prompt.Focus()
		t.components.prompt, promptCmd = t.components.prompt.Update(msg)
		return t, promptCmd

	case error:
		t.state.error = msg
		return t, nil
	}

	return t, tea.Batch(cmds...)
}

func (t *Tui) View() string {
	if t.state.error != nil {
		return "/!\\ " + t.state.error.Error()
	}

	if t.state.running {
		return t.components.renderer.Render(t.state.buffer)
	} else {
		return fmt.Sprintf("\n%s", t.components.prompt.View())
	}
}

func (t *Tui) startEngine(input string) tea.Cmd {
	return func() tea.Msg {
		if t.state.running {
			log.Printf("engine already running")
			return errors.New("engine already running")
		}

		t.state.running = true
		t.state.buffer = ""

		err := t.engine.StreamChatCompletion(input)
		if err != nil {
			log.Printf("StreamChatCompletion error: %v", err)
			return err
		}

		return nil
	}
}

func (t *Tui) awaitEngine() tea.Cmd {
	return func() tea.Msg {
		var output engine.EngineOutput
		output = <-t.engine.Channel()
		t.state.buffer += output.GetContent()
		t.state.running = !output.IsLast()

		return output
	}
}
