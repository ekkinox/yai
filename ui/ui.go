package ui

import (
	"errors"
	"fmt"
	"github.com/ekkinox/yo/run"
	"github.com/ekkinox/yo/ui/prompt"
	"log"
	"strings"

	"github.com/ekkinox/yo/ai"
	"github.com/ekkinox/yo/config"
	"github.com/ekkinox/yo/history"
	"github.com/ekkinox/yo/ui/renderer"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/viper"
)

type UiState struct {
	configuring bool
	running     bool
	confirming  bool
	buffer      string
	command     string
	error       error
}

type UiDimensions struct {
	width  int
	height int
}

type UiComponents struct {
	prompt   *prompt.Prompt
	renderer *renderer.Renderer
}

type Ui struct {
	state      UiState
	dimensions UiDimensions
	components UiComponents
	config     *config.Config
	engine     *ai.Engine
	history    *history.History
}

func NewUi() *Ui {
	return &Ui{
		state: UiState{
			configuring: false,
			running:     false,
			confirming:  false,
			buffer:      "",
			command:     "",
			error:       nil,
		},
		dimensions: UiDimensions{
			150,
			150,
		},
		components: UiComponents{
			prompt: prompt.NewPrompt(prompt.ExecPromptMode),
			renderer: renderer.NewRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(150),
			),
		},
		history: history.NewHistory(),
	}
}

func (u *Ui) Init() tea.Cmd {
	config, err := config.NewConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return u.startConfiguration()
		} else {
			return tea.Sequence(
				tea.Println(u.components.renderer.RenderError(err.Error())),
				tea.Quit,
			)
		}
	}

	return u.startUi(config)
}

func (u *Ui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmds      []tea.Cmd
		promptCmd tea.Cmd
	)

	switch msg := msg.(type) {
	// size
	case tea.WindowSizeMsg:
		u.dimensions.width = msg.Width
		u.dimensions.height = msg.Height
		u.components.renderer = renderer.NewRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(u.dimensions.width),
		)
	// keys
	case tea.KeyMsg:
		switch msg.Type {
		// quit
		case tea.KeyCtrlC:
			return u, tea.Quit
		// history
		case tea.KeyUp, tea.KeyDown:
			if !u.state.running && !u.state.confirming {
				var input *string
				if msg.Type == tea.KeyUp {
					input = u.history.Previous()
				} else {
					input = u.history.Next()
				}
				if input != nil {
					u.components.prompt.Input.SetValue(*input)
					u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
					cmds = append(
						cmds,
						promptCmd,
					)
				}
			}
		// switch mode
		case tea.KeyTab:
			if !u.state.running && !u.state.confirming {
				if u.engine.GetMode() == ai.ChatEngineMode {
					u.engine.SetMode(ai.ExecEngineMode)
					u.components.prompt.ChangeMode(prompt.ExecPromptMode)
				} else {
					u.engine.SetMode(ai.ChatEngineMode)
					u.components.prompt.ChangeMode(prompt.ChatPromptMode)
				}
				u.engine.Reset()
				u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		// enter
		case tea.KeyEnter:
			if u.state.configuring {
				return u, u.finishConfiguration(u.components.prompt.Input.Value())
			}
			if !u.state.running && !u.state.confirming {
				input := u.components.prompt.Input.Value()
				inputStr := u.components.prompt.String()
				u.history.Add(input)
				u.components.prompt.Input.SetValue("")
				u.components.prompt.Input.Blur()
				u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.Println(inputStr),
					u.runEngine(input),
					u.awaitEngine(),
				)
			}

		// clear
		case tea.KeyCtrlL:
			if !u.state.running && !u.state.confirming {
				u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		// reset
		case tea.KeyCtrlR:
			if !u.state.running && !u.state.confirming {
				u.history.Reset()
				u.engine.Reset()
				u.components.prompt.Input.SetValue("")
				u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		// edit settings
		case tea.KeyCtrlS:
			if !u.state.running && !u.state.confirming && !u.state.configuring {
				u.state.buffer = ""
				u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					u.runCommand(
						fmt.Sprintf(
							"%s %s",
							u.config.GetSystemConfig().GetEditor(),
							u.config.GetSystemConfig().GetConfigFile(),
						),
					),
				)
			}

		default:
			if u.state.confirming {
				if strings.ToLower(msg.String()) == "y" {
					u.state.confirming = false
					u.state.buffer = ""
					u.components.prompt.Input.SetValue("")
					return u, tea.Sequence(
						promptCmd,
						u.runCommand(u.state.command),
					)
				} else {
					u.state.confirming = false
					u.state.buffer = ""
					u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
					u.components.prompt.Input.SetValue("")
					u.components.prompt.Input.Focus()
					cmds = append(
						cmds,
						promptCmd,
						tea.Println(fmt.Sprintf("\n%s\n", u.components.renderer.RenderWarning("[cancel]"))),
						textinput.Blink,
					)
				}
			} else {
				u.components.prompt.Input.Focus()
				u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		}
	// engine feedback
	case ai.EngineOutput:
		if msg.IsLast() {
			var output string
			if msg.IsExecutable() {
				u.state.confirming = true
				u.state.command = u.state.buffer
				output = u.components.renderer.RenderContent(fmt.Sprintf("`%s`", u.state.command))
				output += fmt.Sprintf("  Confirm execution? [y/N]")
				u.components.prompt.Input.Blur()
			} else {
				output = u.components.renderer.RenderContent(u.state.buffer)
				u.components.prompt.Input.Focus()
			}
			return u, tea.Sequence(
				promptCmd,
				tea.Println(output),
				textinput.Blink,
			)
		} else {
			return u, u.awaitEngine()
		}
	// runner feedback
	case run.RunOutput:
		u.state.running = false
		u.components.prompt.Input, promptCmd = u.components.prompt.Input.Update(msg)
		u.components.prompt.Input.Focus()
		output := u.components.renderer.RenderSuccess("\n[ok]\n")
		if msg.GetError() != nil {
			output = u.components.renderer.RenderError(fmt.Sprintf("\n[error] %s\n", msg.GetError()))
		}
		cmds = append(
			cmds,
			promptCmd,
			tea.Println(output),
			textinput.Blink,
		)
	// errors
	case error:
		u.state.error = msg
		return u, nil
	}

	return u, tea.Batch(cmds...)
}

func (u *Ui) View() string {
	if u.state.error != nil {
		return u.components.renderer.RenderError(fmt.Sprintf("[error] %s", u.state.error))
	}

	if u.state.configuring {
		return fmt.Sprintf(
			"%s\n%s",
			u.components.renderer.RenderContent(u.state.buffer),
			u.components.prompt.Input.View(),
		)
	}

	if !u.state.running && !u.state.confirming {
		return fmt.Sprintf("%s", u.components.prompt.Input.View())
	}

	if u.engine.GetMode() == ai.ChatEngineMode {
		return u.components.renderer.RenderContent(u.state.buffer)
	} else {
		if u.state.running {
			return u.components.renderer.RenderContent(u.state.buffer)
		}
	}

	return ""
}

func (u *Ui) startUi(config *config.Config) tea.Cmd {
	return func() tea.Msg {
		u.config = config
		engine, err := ai.NewEngine(ai.ExecEngineMode, config)
		if err != nil {
			log.Printf("error: %v", err)
		}
		u.engine = engine
		u.components.prompt = prompt.NewPrompt(prompt.FromString(config.GetUserConfig().GetDefaultMode()))

		return tea.Sequence(
			tea.ClearScreen,
			textinput.Blink,
		)
	}
}

func (u *Ui) startConfiguration() tea.Cmd {
	return func() tea.Msg {
		u.state.configuring = true
		u.state.running = false
		u.state.confirming = false

		u.state.buffer = "**Yo**, welcome! 👋  \n\n"
		u.state.buffer += "I cannot find a configuration file, please enter an **OpenAI API key** "
		u.state.buffer += "from https://platform.openai.com/account/api-keys so I can generate it for you."

		u.state.command = ""

		u.components.prompt = prompt.NewPrompt(prompt.ConfigPromptMode)

		return nil
	}
}

func (u *Ui) finishConfiguration(key string) tea.Cmd {
	return func() tea.Msg {
		u.state.configuring = false

		config, err := config.WriteConfig(key)
		if err != nil {
			log.Println(fmt.Sprintf("error: %v", err))
			return err
		}

		return u.startUi(config)
	}
}

func (u *Ui) runEngine(input string) tea.Cmd {
	return func() tea.Msg {
		if u.state.running {
			log.Printf("engine already running")
			return errors.New("engine already running")
		}

		u.state.running = true
		u.state.confirming = false
		u.state.buffer = ""
		u.state.command = ""

		err := u.engine.StreamChatCompletion(input)
		if err != nil {
			log.Printf("StreamChatCompletion error: %v", err)
			return err
		}

		return nil
	}
}

func (u *Ui) awaitEngine() tea.Cmd {
	return func() tea.Msg {
		var output ai.EngineOutput

		output = <-u.engine.GetChannel()
		u.state.buffer += output.GetContent()
		u.state.running = !output.IsLast()

		return output
	}
}

func (u *Ui) runCommand(input string) tea.Cmd {

	u.state.running = true
	u.state.confirming = false

	c := run.PrepareInteractiveCommand(input)

	return tea.ExecProcess(c, func(error error) tea.Msg {
		u.state.running = false
		u.state.command = ""

		return run.RunOutput{Error: error}
	})
}
