package tui

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/ekkinox/yo/config"
	"github.com/ekkinox/yo/engine"
	"github.com/ekkinox/yo/history"
	"github.com/ekkinox/yo/runner"
	"github.com/ekkinox/yo/tui/components"
	"github.com/spf13/viper"
	"log"
	"os/exec"
	"strings"
)

type TuiState struct {
	configuring bool
	running     bool
	confirming  bool
	buffer      string
	command     string
	error       error
}

type TuiDimensions struct {
	width  int
	height int
}

type TuiComponents struct {
	prompt   textinput.Model
	renderer *components.Renderer
}

type Tui struct {
	state      TuiState
	dimensions TuiDimensions
	components TuiComponents
	config     *config.Config
	engine     *engine.Engine
	runner     *runner.Runner
	history    *history.History
}

func NewTui() *Tui {
	return &Tui{
		state: TuiState{
			configuring: false,
			running:     false,
			confirming:  false,
			buffer:      "",
			command:     "",
			error:       nil,
		},
		dimensions: TuiDimensions{
			150,
			150,
		},
		components: TuiComponents{
			prompt: components.NewPrompt(engine.ChatEngineMode),
			renderer: components.NewRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(150),
			),
		},
		runner:  runner.NewRunner(),
		history: history.NewHistory(),
	}
}

func (t *Tui) Init() tea.Cmd {
	config, err := config.NewConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return t.startConfiguration()
		} else {
			return tea.Sequence(
				tea.Println(t.components.renderer.RenderError(err.Error())),
				tea.Quit,
			)
		}
	}

	t.config = config
	t.engine = engine.NewEngine(config)
	t.components.prompt = components.NewPrompt(engine.EngineModeFromString(config.GetUserPreferences().GetDefaultMode()))

	return tea.Sequence(
		tea.ClearScreen,
		textinput.Blink,
	)
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmds      []tea.Cmd
		promptCmd tea.Cmd
	)

	switch msg := msg.(type) {
	// size
	case tea.WindowSizeMsg:
		t.dimensions.width = msg.Width
		t.dimensions.height = msg.Height
		t.components.renderer = components.NewRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(t.dimensions.width),
		)
	// keys
	case tea.KeyMsg:
		switch msg.Type {
		// quit
		case tea.KeyCtrlC:
			return t, tea.Quit
		// history
		case tea.KeyUp, tea.KeyDown:
			if !t.state.running && !t.state.confirming {
				var input *string
				if msg.Type == tea.KeyUp {
					input = t.history.Previous()
				} else {
					input = t.history.Next()
				}
				if input != nil {
					t.components.prompt.SetValue(*input)
					t.components.prompt, promptCmd = t.components.prompt.Update(msg)
					cmds = append(
						cmds,
						promptCmd,
					)
				}
			}
		// switch mode
		case tea.KeyTab:
			if !t.state.running && !t.state.confirming {
				if t.engine.GetMode() == engine.ChatEngineMode {
					t.engine.SetMode(engine.RunEngineMode)
					t.components.prompt = components.UpdatePrompt(t.components.prompt, engine.RunEngineMode)
				} else {
					t.engine.SetMode(engine.ChatEngineMode)
					t.components.prompt = components.UpdatePrompt(t.components.prompt, engine.ChatEngineMode)
				}
				t.engine.Reset()
				t.components.prompt, promptCmd = t.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		// enter
		case tea.KeyEnter:
			if t.state.configuring {
				return t, t.finishConfiguration(t.components.prompt.Value())
			}
			if !t.state.running && !t.state.confirming {
				input := t.components.prompt.Value()
				t.history.Add(input)
				t.components.prompt.SetValue("")
				t.components.prompt.Blur()
				t.components.prompt, promptCmd = t.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.Println(components.RenderPromptText(input, t.engine.GetMode())),
					t.startEngine(input),
					t.awaitEngine(),
				)
			}

		// clear
		case tea.KeyCtrlL:
			if !t.state.running && !t.state.confirming {
				t.components.prompt, promptCmd = t.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		// reset
		case tea.KeyCtrlR:
			if !t.state.running && !t.state.confirming {
				t.history.Reset()
				t.engine.Reset()
				t.components.prompt.SetValue("")
				t.components.prompt, promptCmd = t.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		// edit settings
		case tea.KeyCtrlS:
			if !t.state.running && !t.state.confirming {
				t.state.buffer = ""
				t.components.prompt, promptCmd = t.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					t.startCommand(
						fmt.Sprintf(
							"%s %s",
							t.config.GetContext().GetEditor(),
							t.config.GetContext().GetConfigFile(),
						),
					),
				)
			}

		default:
			if t.state.confirming {
				if strings.ToLower(msg.String()) == "y" {
					t.state.confirming = false
					t.state.buffer = ""
					t.components.prompt.SetValue("")
					return t, tea.Sequence(
						promptCmd,
						t.startCommand(t.state.command),
					)
				} else {
					t.state.confirming = false
					t.state.buffer = ""
					t.components.prompt, promptCmd = t.components.prompt.Update(msg)
					t.components.prompt.SetValue("")
					t.components.prompt.Focus()
					cmds = append(
						cmds,
						promptCmd,
						tea.Println(fmt.Sprintf("\n%s\n", t.components.renderer.RenderWarning("[cancel]"))),
						textinput.Blink,
					)
				}
			} else {
				t.components.prompt.Focus()
				t.components.prompt, promptCmd = t.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		}
	// engine feedback
	case engine.EngineOutput:
		if msg.IsLast() {
			var output string
			if msg.IsExecutable() {
				t.state.confirming = true
				t.state.command = t.state.buffer
				output = t.components.renderer.RenderContent(fmt.Sprintf("`%s`", t.state.buffer))
				output += fmt.Sprintf("  Confirm execution? [y/N]")
				t.components.prompt.Blur()
			} else {
				output = t.components.renderer.RenderContent(t.state.buffer)
				t.components.prompt.Focus()
			}
			return t, tea.Sequence(
				promptCmd,
				tea.Println(output),
				textinput.Blink,
			)
		} else {
			return t, t.awaitEngine()
		}
	// runner feedback
	case runner.RunnerOutput:
		t.state.running = false
		t.components.prompt, promptCmd = t.components.prompt.Update(msg)
		t.components.prompt.Focus()
		output := t.components.renderer.RenderSuccess("\n[ok]\n")
		if msg.GetError() != nil {
			output = t.components.renderer.RenderError(fmt.Sprintf("\n[error] %s\n", msg.GetError()))
		}
		cmds = append(
			cmds,
			promptCmd,
			tea.Println(output),
			textinput.Blink,
		)
	// errors
	case error:
		t.state.error = msg
		return t, nil
	}

	return t, tea.Batch(cmds...)
}

func (t *Tui) View() string {
	if t.state.error != nil {
		return t.components.renderer.RenderError(fmt.Sprintf("[error] %s", t.state.error))
	}

	if t.state.configuring {
		return fmt.Sprintf(
			"%s\n%s",
			t.components.renderer.RenderContent(t.state.buffer),
			t.components.prompt.View(),
		)
	}

	if !t.state.running && !t.state.confirming {
		return fmt.Sprintf("%s", t.components.prompt.View())
	}

	if t.engine.GetMode() == engine.ChatEngineMode {
		return t.components.renderer.RenderContent(t.state.buffer)
	} else {
		if t.state.running {
			return t.components.renderer.RenderContent(t.state.buffer)
		}
	}

	return ""
}

func (t *Tui) startConfiguration() tea.Cmd {
	return func() tea.Msg {
		t.state.configuring = true
		t.state.running = false
		t.state.confirming = false

		t.state.buffer = "**Yo**, welcome! 👋  \n\n"
		t.state.buffer += "I cannot find a configuration file, please enter an **OpenAI API key** "
		t.state.buffer += "from https://platform.openai.com/account/api-keys so I can generate it for you."

		t.state.command = ""

		t.components.prompt = components.NewConfigPrompt()

		return nil
	}
}

func (t *Tui) finishConfiguration(key string) tea.Cmd {
	return func() tea.Msg {
		t.state.configuring = false

		config, err := config.WriteConfig(key)
		if err != nil {
			log.Println(fmt.Sprintf("error: %v", err))
			return err
		}

		t.config = config
		t.engine = engine.NewEngine(config)
		t.state.buffer = ""
		t.components.prompt = components.NewPrompt(engine.ChatEngineMode)

		return tea.Sequence(
			tea.Println(fmt.Sprintf("\n\nConfig generated in %s.\n\n", t.config.GetContext().GetConfigFile())),
			textinput.Blink,
		)
	}
}

func (t *Tui) startEngine(input string) tea.Cmd {
	return func() tea.Msg {
		if t.state.running {
			log.Printf("engine already running")
			return errors.New("engine already running")
		}

		t.state.running = true
		t.state.confirming = false
		t.state.buffer = ""
		t.state.command = ""

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

func (t *Tui) startCommand(input string) tea.Cmd {

	t.state.running = true
	t.state.confirming = false

	c := exec.Command("bash", "-c", fmt.Sprintf("%s; echo \"\n\"", strings.TrimRight(input, ";")))

	return tea.ExecProcess(c, func(error error) tea.Msg {
		t.state.running = false
		t.state.command = ""

		return runner.RunnerOutput{Error: error}
	})
}
