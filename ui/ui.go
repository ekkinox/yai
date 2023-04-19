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

type UiState struct {
	error       error
	runMode     RunMode
	promptMode  PromptMode
	configuring bool
	querying    bool
	confirming  bool
	executing   bool
	args        string
	pipe        string
	buffer      string
	command     string
}

type UiDimensions struct {
	width  int
	height int
}

type UiComponents struct {
	prompt   *Prompt
	renderer *Renderer
	spinner  *Spinner
}

type Ui struct {
	state      UiState
	dimensions UiDimensions
	components UiComponents
	config     *config.Config
	engine     *ai.Engine
	history    *history.History
}

func NewUi(input *UiInput) *Ui {
	return &Ui{
		state: UiState{
			error:       nil,
			runMode:     input.GetRunMode(),
			promptMode:  input.GetPromptMode(),
			configuring: false,
			querying:    false,
			confirming:  false,
			executing:   false,
			args:        input.GetArgs(),
			pipe:        input.GetPipe(),
			buffer:      "",
			command:     "",
		},
		dimensions: UiDimensions{
			150,
			150,
		},
		components: UiComponents{
			prompt: NewPrompt(input.GetPromptMode()),
			renderer: NewRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(150),
			),
			spinner: NewSpinner(),
		},
		history: history.NewHistory(),
	}
}

func (u *Ui) Init() tea.Cmd {
	config, err := config.NewConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if u.state.runMode == ReplMode {
				return tea.Sequence(
					tea.ClearScreen,
					u.startConfig(),
				)
			} else {
				return u.startConfig()
			}
		} else {
			return tea.Sequence(
				tea.Println(u.components.renderer.RenderError(err.Error())),
				tea.Quit,
			)
		}
	}

	if u.state.runMode == ReplMode {
		return u.startRepl(config)
	} else {
		return u.startCli(config)
	}
}

func (u *Ui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds       []tea.Cmd
		promptCmd  tea.Cmd
		spinnerCmd tea.Cmd
	)

	switch msg := msg.(type) {
	// spinner
	case spinner.TickMsg:
		if u.state.querying {
			u.components.spinner, spinnerCmd = u.components.spinner.Update(msg)
			cmds = append(
				cmds,
				spinnerCmd,
			)
		}
	// size
	case tea.WindowSizeMsg:
		u.dimensions.width = msg.Width
		u.dimensions.height = msg.Height
		u.components.renderer = NewRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(u.dimensions.width),
		)
	// keyboard
	case tea.KeyMsg:
		switch msg.Type {
		// quit
		case tea.KeyCtrlC:
			return u, tea.Quit
		// history
		case tea.KeyUp, tea.KeyDown:
			if !u.state.querying && !u.state.confirming {
				var input *string
				if msg.Type == tea.KeyUp {
					input = u.history.GetPrevious()
				} else {
					input = u.history.GetNext()
				}
				if input != nil {
					u.components.prompt.SetValue(*input)
					u.components.prompt, promptCmd = u.components.prompt.Update(msg)
					cmds = append(
						cmds,
						promptCmd,
					)
				}
			}
		// switch mode
		case tea.KeyTab:
			if !u.state.querying && !u.state.confirming {
				if u.state.promptMode == ChatPromptMode {
					u.state.promptMode = ExecPromptMode
					u.components.prompt.SetMode(ExecPromptMode)
					u.engine.SetMode(ai.ExecEngineMode)
				} else {
					u.state.promptMode = ChatPromptMode
					u.components.prompt.SetMode(ChatPromptMode)
					u.engine.SetMode(ai.ChatEngineMode)
				}
				u.engine.Reset()
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		// enter
		case tea.KeyEnter:
			if u.state.configuring {
				return u, u.finishConfig(u.components.prompt.GetValue())
			}
			if !u.state.querying && !u.state.confirming {
				input := u.components.prompt.GetValue()
				if input != "" {
					inputPrint := u.components.prompt.AsString()
					u.history.Add(input)
					u.components.prompt.SetValue("")
					u.components.prompt.Blur()
					u.components.prompt, promptCmd = u.components.prompt.Update(msg)
					if u.state.promptMode == ChatPromptMode {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(inputPrint),
							u.startChatStream(input),
							u.awaitChatStream(),
						)
					} else {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(inputPrint),
							u.startExec(input),
							u.components.spinner.Tick,
						)
					}
				}
			}

		// help
		case tea.KeyCtrlH:
			if !u.state.configuring && !u.state.querying && !u.state.confirming {
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.Println(u.components.renderer.RenderContent(u.components.renderer.RenderHelpMessage())),
					textinput.Blink,
				)
			}

		// clear
		case tea.KeyCtrlL:
			if !u.state.querying && !u.state.confirming {
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		// reset
		case tea.KeyCtrlR:
			if !u.state.querying && !u.state.confirming {
				u.history.Reset()
				u.engine.Reset()
				u.components.prompt.SetValue("")
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}

		// edit settings
		case tea.KeyCtrlS:
			if !u.state.querying && !u.state.confirming && !u.state.configuring && !u.state.executing {
				u.state.executing = true
				u.state.buffer = ""
				u.state.command = ""
				u.components.prompt.Blur()
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					u.editSettings(),
				)
			}

		default:
			if u.state.confirming {
				if strings.ToLower(msg.String()) == "y" {
					u.state.confirming = false
					u.state.executing = true
					u.state.buffer = ""
					u.components.prompt.SetValue("")
					return u, tea.Sequence(
						promptCmd,
						u.execCommand(u.state.command),
					)
				} else {
					u.state.confirming = false
					u.state.executing = false
					u.state.buffer = ""
					u.components.prompt, promptCmd = u.components.prompt.Update(msg)
					u.components.prompt.SetValue("")
					u.components.prompt.Focus()
					if u.state.runMode == ReplMode {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(fmt.Sprintf("\n%s\n", u.components.renderer.RenderWarning("[cancel]"))),
							textinput.Blink,
						)
					} else {
						return u, tea.Sequence(
							promptCmd,
							tea.Println(fmt.Sprintf("\n%s\n", u.components.renderer.RenderWarning("[cancel]"))),
							tea.Quit,
						)
					}
				}
				u.state.command = ""
			} else {
				u.components.prompt.Focus()
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		}
	// engine exec feedback
	case ai.EngineExecOutput:
		var output string
		if msg.IsExecutable() {
			u.state.confirming = true
			u.state.command = msg.GetCommand()
			output = u.components.renderer.RenderContent(fmt.Sprintf("`%s`", u.state.command))
			output += fmt.Sprintf("  %s\n\n  confirm execution? [y/N]", u.components.renderer.RenderHelp(msg.GetExplanation()))
			u.components.prompt.Blur()
		} else {
			output = u.components.renderer.RenderContent(msg.GetExplanation())
			u.components.prompt.Focus()
			if u.state.runMode == CliMode {
				return u, tea.Sequence(
					tea.Println(output),
					tea.Quit,
				)
			}
		}
		u.components.prompt, promptCmd = u.components.prompt.Update(msg)
		return u, tea.Sequence(
			promptCmd,
			textinput.Blink,
			tea.Println(output),
		)
	// engine chat stream feedback
	case ai.EngineChatStreamOutput:
		if msg.IsLast() {
			output := u.components.renderer.RenderContent(u.state.buffer)
			u.state.buffer = ""
			u.components.prompt.Focus()
			if u.state.runMode == CliMode {
				return u, tea.Sequence(
					tea.Println(output),
					tea.Quit,
				)
			} else {
				return u, tea.Sequence(
					tea.Println(output),
					textinput.Blink,
				)
			}
		} else {
			return u, u.awaitChatStream()
		}
	// runner feedback
	case run.RunOutput:
		u.state.querying = false
		u.components.prompt, promptCmd = u.components.prompt.Update(msg)
		u.components.prompt.Focus()
		output := u.components.renderer.RenderSuccess(fmt.Sprintf("\n%s\n", msg.GetSuccessMessage()))
		if msg.HasError() {
			output = u.components.renderer.RenderError(fmt.Sprintf("\n%s\n", msg.GetErrorMessage()))
		}
		if u.state.runMode == CliMode {
			return u, tea.Sequence(
				tea.Println(output),
				tea.Quit,
			)
		} else {
			return u, tea.Sequence(
				tea.Println(output),
				promptCmd,
				textinput.Blink,
			)
		}
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
			u.components.prompt.View(),
		)
	}

	if !u.state.querying && !u.state.confirming && !u.state.executing {
		return u.components.prompt.View()
	}

	if u.state.promptMode == ChatPromptMode {
		return u.components.renderer.RenderContent(u.state.buffer)
	} else {
		if u.state.querying {
			return u.components.spinner.View()
		} else {
			if !u.state.executing {
				return u.components.renderer.RenderContent(u.state.buffer)
			}
		}
	}

	return ""
}

func (u *Ui) startRepl(config *config.Config) tea.Cmd {
	return tea.Sequence(
		tea.ClearScreen,
		tea.Println(u.components.renderer.RenderContent(u.components.renderer.RenderHelpMessage())),
		textinput.Blink,
		func() tea.Msg {
			u.config = config

			if u.state.promptMode == DefaultPromptMode {
				u.state.promptMode = GetPromptModeFromString(config.GetUserConfig().GetDefaultPromptMode())
			}

			engineMode := ai.ExecEngineMode
			if u.state.promptMode == ChatPromptMode {
				engineMode = ai.ChatEngineMode
			}

			engine, err := ai.NewEngine(engineMode, config)
			if err != nil {
				return err
			}

			if u.state.pipe != "" {
				engine.SetPipe(u.state.pipe)
			}

			u.engine = engine
			u.state.buffer = "Welcome \n\n"
			u.state.command = ""
			u.components.prompt = NewPrompt(u.state.promptMode)

			return nil
		},
	)
}

func (u *Ui) startCli(config *config.Config) tea.Cmd {
	u.config = config

	if u.state.promptMode == DefaultPromptMode {
		u.state.promptMode = GetPromptModeFromString(config.GetUserConfig().GetDefaultPromptMode())
	}

	engineMode := ai.ExecEngineMode
	if u.state.promptMode == ChatPromptMode {
		engineMode = ai.ChatEngineMode
	}

	engine, err := ai.NewEngine(engineMode, config)
	if err != nil {
		u.state.error = err
		return nil
	}

	if u.state.pipe != "" {
		engine.SetPipe(u.state.pipe)
	}

	u.engine = engine
	u.state.querying = true
	u.state.confirming = false
	u.state.buffer = ""
	u.state.command = ""

	if u.state.promptMode == ExecPromptMode {
		return tea.Batch(
			u.components.spinner.Tick,
			func() tea.Msg {
				output, err := u.engine.ExecCompletion(u.state.args)
				u.state.querying = false
				if err != nil {
					return err
				}

				return *output
			},
		)
	} else {
		return tea.Batch(
			u.startChatStream(u.state.args),
			u.awaitChatStream(),
		)
	}
}

func (u *Ui) startConfig() tea.Cmd {
	return func() tea.Msg {
		u.state.configuring = true
		u.state.querying = false
		u.state.confirming = false
		u.state.executing = false

		u.state.buffer = u.components.renderer.RenderConfigMessage()
		u.state.command = ""
		u.components.prompt = NewPrompt(ConfigPromptMode)

		return nil
	}
}

func (u *Ui) finishConfig(key string) tea.Cmd {
	u.state.configuring = false

	config, err := config.WriteConfig(key, true)
	if err != nil {
		u.state.error = err
		return nil
	}

	u.config = config
	engine, err := ai.NewEngine(ai.ExecEngineMode, config)
	if err != nil {
		u.state.error = err
		return nil
	}

	if u.state.pipe != "" {
		engine.SetPipe(u.state.pipe)
	}

	u.engine = engine

	if u.state.runMode == ReplMode {
		return tea.Sequence(
			tea.ClearScreen,
			tea.Println(u.components.renderer.RenderSuccess("\n[settings ok]\n")),
			textinput.Blink,
			func() tea.Msg {
				u.state.buffer = ""
				u.state.command = ""
				u.components.prompt = NewPrompt(ExecPromptMode)

				return nil
			},
		)
	} else {
		if u.state.promptMode == ExecPromptMode {
			u.state.querying = true
			u.state.configuring = false
			u.state.buffer = ""
			return tea.Sequence(
				tea.Println(u.components.renderer.RenderSuccess("\n[settings ok]")),
				u.components.spinner.Tick,
				func() tea.Msg {
					output, err := u.engine.ExecCompletion(u.state.args)
					u.state.querying = false
					if err != nil {
						return err
					}

					return *output
				},
			)
		} else {
			return tea.Batch(
				u.startChatStream(u.state.args),
				u.awaitChatStream(),
			)
		}
	}
}

func (u *Ui) startExec(input string) tea.Cmd {
	return func() tea.Msg {
		u.state.querying = true
		u.state.confirming = false
		u.state.buffer = ""
		u.state.command = ""

		output, err := u.engine.ExecCompletion(input)
		u.state.querying = false
		if err != nil {
			return err
		}

		return *output
	}
}

func (u *Ui) startChatStream(input string) tea.Cmd {
	return func() tea.Msg {
		u.state.querying = true
		u.state.executing = false
		u.state.confirming = false
		u.state.buffer = ""
		u.state.command = ""

		err := u.engine.ChatStreamCompletion(input)
		if err != nil {
			return err
		}

		return nil
	}
}

func (u *Ui) awaitChatStream() tea.Cmd {
	return func() tea.Msg {
		output := <-u.engine.GetChannel()
		u.state.buffer += output.GetContent()
		u.state.querying = !output.IsLast()

		return output
	}
}

func (u *Ui) execCommand(input string) tea.Cmd {
	u.state.querying = false
	u.state.confirming = false
	u.state.executing = true

	c := run.PrepareInteractiveCommand(input)

	return tea.ExecProcess(c, func(error error) tea.Msg {
		u.state.executing = false
		u.state.command = ""

		return run.NewRunOutput(error, "[error]", "[ok]")
	})
}

func (u *Ui) editSettings() tea.Cmd {
	u.state.querying = false
	u.state.confirming = false
	u.state.executing = true

	c := run.PrepareEditSettingsCommand(fmt.Sprintf(
		"%s %s",
		u.config.GetSystemConfig().GetEditor(),
		u.config.GetSystemConfig().GetConfigFile(),
	))

	return tea.ExecProcess(c, func(error error) tea.Msg {
		u.state.executing = false
		u.state.command = ""

		if error != nil {
			return run.NewRunOutput(error, "[settings error]", "")
		}

		config, error := config.NewConfig()
		if error != nil {
			return run.NewRunOutput(error, "[settings error]", "")
		}

		u.config = config
		engine, error := ai.NewEngine(ai.ExecEngineMode, config)
		if u.state.pipe != "" {
			engine.SetPipe(u.state.pipe)
		}
		if error != nil {
			return run.NewRunOutput(error, "[settings error]", "")
		}
		u.engine = engine

		return run.NewRunOutput(nil, "", "[settings ok]")
	})
}
