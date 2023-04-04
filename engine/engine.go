package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/ekkinox/yo/config"
	"io"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const NORUN = "NORUN"

type EngineOutput struct {
	content    string
	last       bool
	executable bool
}

func (d EngineOutput) GetContent() string {
	return d.content
}

func (d EngineOutput) IsLast() bool {
	return d.last
}

func (d EngineOutput) IsExecutable() bool {
	return d.executable
}

type Engine struct {
	config   *config.Config
	client   *openai.Client
	messages []openai.ChatCompletionMessage
	channel  chan EngineOutput
	mode     EngineMode
	running  bool
}

func NewEngine(config *config.Config) *Engine {
	return &Engine{
		config:   config,
		client:   openai.NewClient(config.GetOpenAI().GetKey()),
		messages: make([]openai.ChatCompletionMessage, 0),
		channel:  make(chan EngineOutput),
		mode:     EngineModeFromString(config.GetUserPreferences().GetDefaultMode()),
		running:  false,
	}
}

func (e *Engine) Channel() chan EngineOutput {
	return e.channel
}

func (e *Engine) Interrupt() *Engine {
	e.channel <- EngineOutput{
		content:    "Interrupt",
		last:       true,
		executable: false,
	}

	e.running = false

	return e
}

func (e *Engine) Reset() *Engine {
	e.messages = []openai.ChatCompletionMessage{}

	return e
}

func (e *Engine) SetMode(mode EngineMode) *Engine {
	e.mode = mode

	return e
}

func (e *Engine) GetMode() EngineMode {
	return e.mode
}

func (e *Engine) StreamChatCompletion(input string) error {

	ctx := context.Background()

	e.running = true

	e.appendUserMessage(input)

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 1000,
		Messages:  e.prepareCompletionMessages(),
		Stream:    true,
	}

	stream, err := e.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		log.Printf("error on stream creation: %v", err)
		return err
	}
	defer stream.Close()

	var output string

	for {
		if e.running {
			resp, err := stream.Recv()

			if errors.Is(err, io.EOF) {
				executable := false
				if e.mode == RunEngineMode {
					if !strings.HasPrefix(output, NORUN) {
						executable = true
					}
				}

				e.channel <- EngineOutput{
					content:    "",
					last:       true,
					executable: executable,
				}
				e.running = false
				e.appendAssistantMessage(output)

				return nil
			}

			if err != nil {
				log.Printf("error on stream read: %v", err)
				e.running = false
				return err
			}

			delta := resp.Choices[0].Delta.Content

			output += delta

			e.channel <- EngineOutput{
				content: delta,
				last:    false,
			}

			//time.Sleep(time.Microsecond * 100)
		} else {
			stream.Close()

			return nil
		}
	}
}

func (e *Engine) appendUserMessage(content string) *Engine {
	e.messages = append(e.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})

	return e
}

func (e *Engine) appendAssistantMessage(content string) *Engine {
	e.messages = append(e.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})

	return e
}

func (e *Engine) prepareCompletionMessages() []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: e.prepareSystemMessageContent(),
		},
	}
	for _, m := range e.messages {
		messages = append(messages, m)
	}

	return messages
}

func (e *Engine) prepareSystemMessageContent() string {
	prompt := "You are Yo, an helpful AI command line assistant running in a terminal, created by github.com/ekkinox. "

	switch e.mode {
	case ChatEngineMode:
		prompt += "You will provide an answer for my input the most helpful possible, rendered in markdown format. "
	case RunEngineMode:
		prompt += "You will prepare a single command line that fulfills my input, and you will NEVER provide any explanation or descriptive text even if you did before in the discussion. "
		prompt += "This command line cannot have new lines, use instead separators like && and ;. "
		prompt += fmt.Sprintf("If you do NOT manage to generate a single command line, in this case prefix your answer with %s. ", NORUN)
	}

	prompt += "My context: "
	if e.config.GetContext().GetOperatingSystem() != "other" {
		prompt += fmt.Sprintf("my operating system is %s, ", e.config.GetContext().GetOperatingSystem())
	}
	if e.config.GetContext().GetDistribution() != "" {
		prompt += fmt.Sprintf("my distribution is %s, ", e.config.GetContext().GetDistribution())
	}
	if e.config.GetContext().GetHomeDirectory() != "" {
		prompt += fmt.Sprintf("my home directory is %s, ", e.config.GetContext().GetHomeDirectory())
	}
	if e.config.GetContext().GetShell() != "" {
		prompt += fmt.Sprintf("my shell is %s, ", e.config.GetContext().GetShell())
	}
	if e.config.GetContext().GetShell() != "" {
		prompt += fmt.Sprintf("my editor is %s, ", e.config.GetContext().GetEditor())
	}
	prompt += "take this into account. "

	if e.config.GetUserPreferences().GetContext() != "" {
		prompt += fmt.Sprintf("Also, %s ", e.config.GetUserPreferences().GetContext())
	}

	return prompt
}
