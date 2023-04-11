package ui

import (
	"fmt"
	"strings"

	"github.com/ekkinox/yo/ai"
	"github.com/ekkinox/yo/config"
	"github.com/ekkinox/yo/history"
	"github.com/ekkinox/yo/run"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/viper"
)

type ReplState struct {
	configuring bool
	querying    bool
	confirming  bool
	executing   bool
	buffer      string
	command     string
	error       error
}

type ReplDimensions struct {
	width  int
	height int
}

type ReplComponents struct {
	prompt   *Prompt
	renderer *Renderer
	spinner  *Spinner
}

type Repl struct {
	state      ReplState
	dimensions ReplDimensions
	components ReplComponents
	config     *config.Config
	engine     *ai.Engine
	history    *history.History
}

func NewRepl() *Repl {
	return &Repl{
		state: ReplState{
			configuring: false,
			querying:    false,
			confirming:  false,
			executing:   false,
			buffer:      "",
			command:     "",
			error:       nil,
		},
		dimensions: ReplDimensions{
			150,
			150,
		},
		components: ReplComponents{
			prompt: NewPrompt(ExecPromptMode),
			renderer: NewRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(150),
			),
			spinner: NewSpinner(),
		},
		history: history.NewHistory(),
	}
}

func (r *Repl) Init() tea.Cmd {
	config, err := config.NewConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return r.startConfig()
		} else {
			return tea.Sequence(
				tea.Println(r.components.renderer.RenderError(err.Error())),
				tea.Quit,
			)
		}
	}

	return r.startRepl(config)
}

func (r *Repl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmds       []tea.Cmd
		promptCmd  tea.Cmd
		spinnerCmd tea.Cmd
	)

	switch msg := msg.(type) {
	// spinner
	case spinner.TickMsg:
		if r.state.querying {
			r.components.spinner, spinnerCmd = r.components.spinner.Update(msg)
			cmds = append(
				cmds,
				spinnerCmd,
			)
		}
	// size
	case tea.WindowSizeMsg:
		r.dimensions.width = msg.Width
		r.dimensions.height = msg.Height
		r.components.renderer = NewRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(r.dimensions.width),
		)
	// keyboard
	case tea.KeyMsg:
		switch msg.Type {
		// quit
		case tea.KeyCtrlC:
			return r, tea.Quit
		// history
		case tea.KeyUp, tea.KeyDown:
			if !r.state.querying && !r.state.confirming {
				var input *string
				if msg.Type == tea.KeyUp {
					input = r.history.GetPrevious()
				} else {
					input = r.history.GetNext()
				}
				if input != nil {
					r.components.prompt.SetValue(*input)
					r.components.prompt, promptCmd = r.components.prompt.Update(msg)
					cmds = append(
						cmds,
						promptCmd,
					)
				}
			}
		// switch mode
		case tea.KeyTab:
			if !r.state.querying && !r.state.confirming {
				if r.engine.GetMode() == ai.ChatEngineMode {
					r.engine.SetMode(ai.ExecEngineMode)
					r.components.prompt.SetMode(ExecPromptMode)
				} else {
					r.engine.SetMode(ai.ChatEngineMode)
					r.components.prompt.SetMode(ChatPromptMode)
				}
				r.engine.Reset()
				r.components.prompt, promptCmd = r.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		// enter
		case tea.KeyEnter:
			if r.state.configuring {
				return r, r.finishConfig(r.components.prompt.GetValue())
			}
			if !r.state.querying && !r.state.confirming {
				input := r.components.prompt.GetValue()
				if input != "" {
					inputPrint := r.components.prompt.AsString()
					r.history.Add(input)
					r.components.prompt.SetValue("")
					r.components.prompt.Blur()
					r.components.prompt, promptCmd = r.components.prompt.Update(msg)
					if r.engine.GetMode() == ai.ChatEngineMode {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(inputPrint),
							r.startChatStream(input),
							r.awaitChatStream(),
						)
					} else {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(inputPrint),
							r.startExec(input),
							r.components.spinner.Tick,
						)
					}
				}
			}

		// clear
		case tea.KeyCtrlL:
			if !r.state.querying && !r.state.confirming {
				r.components.prompt, promptCmd = r.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		// reset
		case tea.KeyCtrlR:
			if !r.state.querying && !r.state.confirming {
				r.history.Reset()
				r.engine.Reset()
				r.components.prompt.SetValue("")
				r.components.prompt, promptCmd = r.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		// edit settings
		case tea.KeyCtrlS:
			if !r.state.querying && !r.state.confirming && !r.state.configuring && !r.state.executing {
				r.state.executing = true
				r.state.buffer = ""
				r.state.command = ""
				r.components.prompt.Blur()
				r.components.prompt, promptCmd = r.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					r.editSettings(),
				)
			}

		default:
			if r.state.confirming {
				if strings.ToLower(msg.String()) == "y" {
					r.state.confirming = false
					r.state.executing = true
					r.state.buffer = ""
					r.components.prompt.SetValue("")
					return r, tea.Sequence(
						promptCmd,
						r.execCommand(r.state.command),
					)
				} else {
					r.state.confirming = false
					r.state.executing = false
					r.state.buffer = ""
					r.components.prompt, promptCmd = r.components.prompt.Update(msg)
					r.components.prompt.SetValue("")
					r.components.prompt.Focus()
					cmds = append(
						cmds,
						promptCmd,
						tea.Println(fmt.Sprintf("\n%s\n", r.components.renderer.RenderWarning("[cancel]"))),
						textinput.Blink,
					)
				}
				r.state.command = ""
			} else {
				r.components.prompt.Focus()
				r.components.prompt, promptCmd = r.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		}
	// engine exec feedback
	case ai.EngineExecOutput:
		//r.state.querying = false
		var output string
		if msg.IsExecutable() {
			r.state.confirming = true
			r.state.command = msg.GetCommand()
			output = r.components.renderer.RenderContent(fmt.Sprintf("`%s`", r.state.command))
			output += fmt.Sprintf("  %s\n\n  confirm execution? [y/N]", r.components.renderer.RenderHelp(msg.GetExplanation()))
			r.components.prompt.Blur()
		} else {
			output = r.components.renderer.RenderContent(msg.GetExplanation())
			r.components.prompt.Focus()
		}
		r.components.prompt, promptCmd = r.components.prompt.Update(msg)
		return r, tea.Sequence(
			promptCmd,
			textinput.Blink,
			tea.Println(output),
		)
	// engine chat stream feedback
	case ai.EngineChatStreamOutput:
		if msg.IsLast() {
			output := r.components.renderer.RenderContent(r.state.buffer)
			r.state.buffer = ""
			r.components.prompt.Focus()
			return r, tea.Sequence(
				tea.Println(output),
				textinput.Blink,
			)
		} else {
			return r, r.awaitChatStream()
		}
	// runner feedback
	case run.RunOutput:
		r.state.querying = false
		r.components.prompt, promptCmd = r.components.prompt.Update(msg)
		r.components.prompt.Focus()
		output := r.components.renderer.RenderSuccess(fmt.Sprintf("\n%s\n", msg.GetSuccessMessage()))
		if msg.HasError() {
			output = r.components.renderer.RenderError(fmt.Sprintf("\n%s\n", msg.GetErrorMessage()))
		}
		return r, tea.Sequence(
			tea.Println(output),
			promptCmd,
			textinput.Blink,
		)
	// errors
	case error:
		r.state.error = msg
		return r, nil
	}

	return r, tea.Batch(cmds...)
}

func (r *Repl) View() string {
	if r.state.error != nil {
		return r.components.renderer.RenderError(fmt.Sprintf("[error] %s", r.state.error))
	}

	if r.state.configuring {
		return fmt.Sprintf(
			"%s\n%s",
			r.components.renderer.RenderContent(r.state.buffer),
			r.components.prompt.View(),
		)
	}

	if !r.state.querying && !r.state.confirming && !r.state.executing {
		return fmt.Sprintf("%s", r.components.prompt.View())
	}

	if r.engine.GetMode() == ai.ChatEngineMode {
		return r.components.renderer.RenderContent(r.state.buffer)
	} else {
		if r.state.querying {
			return r.components.spinner.View()
		} else {
			if !r.state.executing {
				return r.components.renderer.RenderContent(r.state.buffer)
			}
		}
	}

	return ""
}

func (r *Repl) startRepl(config *config.Config) tea.Cmd {
	return func() tea.Msg {
		r.config = config
		engine, err := ai.NewEngine(ai.ExecEngineMode, config)
		if err != nil {
			return err
		}

		r.engine = engine
		r.state.buffer = "Welcome \n\n"
		r.state.command = ""
		r.components.prompt = NewPrompt(ExecPromptMode)

		return tea.Sequence(
			tea.ClearScreen,
			textinput.Blink,
		)
	}
}

func (r *Repl) startConfig() tea.Cmd {
	return func() tea.Msg {
		r.state.configuring = true
		r.state.querying = false
		r.state.confirming = false
		r.state.executing = false

		r.state.buffer = "**Yo**, welcome! 👋  \n\n"
		r.state.buffer += "I cannot find a configuration file, please enter an **OpenAI API key** "
		r.state.buffer += "from https://platform.openai.com/account/api-keys so I can generate it for you."

		r.state.command = ""

		r.components.prompt = NewPrompt(ConfigPromptMode)

		return nil
	}
}

func (r *Repl) finishConfig(key string) tea.Cmd {
	return func() tea.Msg {
		r.state.configuring = false

		config, err := config.WriteConfig(key)
		if err != nil {
			return err
		}

		r.config = config
		engine, err := ai.NewEngine(ai.ExecEngineMode, config)
		if err != nil {
			return err
		}

		r.engine = engine
		r.state.buffer = ""
		r.state.command = ""
		r.components.prompt = NewPrompt(ExecPromptMode)

		return tea.Sequence(
			tea.ClearScreen,
			textinput.Blink,
		)
	}
}

func (r *Repl) startExec(input string) tea.Cmd {
	return func() tea.Msg {
		r.state.querying = true
		r.state.confirming = false
		r.state.buffer = ""
		r.state.command = ""

		output, err := r.engine.ExecCompletion(input)
		r.state.querying = false
		if err != nil {
			return err
		}

		return *output
	}
}

func (r *Repl) startChatStream(input string) tea.Cmd {
	return func() tea.Msg {
		r.state.querying = true
		r.state.executing = false
		r.state.confirming = false
		r.state.buffer = ""
		r.state.command = ""

		err := r.engine.ChatStreamCompletion(input)
		if err != nil {
			return err
		}

		return nil
	}
}

func (r *Repl) awaitChatStream() tea.Cmd {
	return func() tea.Msg {
		var output ai.EngineChatStreamOutput

		output = <-r.engine.GetChannel()
		r.state.buffer += output.GetContent()
		r.state.querying = !output.IsLast()

		return output
	}
}

func (r *Repl) execCommand(input string) tea.Cmd {

	r.state.querying = false
	r.state.confirming = false
	r.state.executing = true

	c := run.PrepareInteractiveCommand(input)

	return tea.ExecProcess(c, func(error error) tea.Msg {
		r.state.executing = false
		r.state.command = ""

		return run.NewRunOutput(error, "[error]", "[ok]")
	})
}

func (r *Repl) editSettings() tea.Cmd {

	r.state.querying = false
	r.state.confirming = false
	r.state.executing = true

	c := run.PrepareEditSettingsCommand(fmt.Sprintf(
		"%s %s",
		r.config.GetSystemConfig().GetEditor(),
		r.config.GetSystemConfig().GetConfigFile(),
	))

	return tea.ExecProcess(c, func(error error) tea.Msg {
		r.state.executing = false
		r.state.command = ""

		if error != nil {
			return run.NewRunOutput(error, "[settings edition error]", "")
		}

		config, error := config.NewConfig()
		if error != nil {
			return run.NewRunOutput(error, "[settings edition error]", "")
		}

		r.config = config
		engine, error := ai.NewEngine(ai.ExecEngineMode, config)
		if error != nil {
			return run.NewRunOutput(error, "[settings edition error]", "")
		}
		r.engine = engine

		return run.NewRunOutput(nil, "", "[settings edition success]")
	})
}
