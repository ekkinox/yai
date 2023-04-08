package ai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ekkinox/yo/config"
	"github.com/ekkinox/yo/system"

	"github.com/sashabaranov/go-openai"
)

const noexec = "[noexec]"

type Engine struct {
	mode         EngineMode
	config       *config.Config
	client       *openai.Client
	execMessages []openai.ChatCompletionMessage
	chatMessages []openai.ChatCompletionMessage
	channel      chan EngineOutput
	running      bool
}

func NewEngine(mode EngineMode, config *config.Config) (*Engine, error) {

	var client *openai.Client

	if config.GetAiConfig().GetProxy() != "" {

		clientConfig := openai.DefaultConfig(config.GetAiConfig().GetKey())

		proxyUrl, err := url.Parse(config.GetAiConfig().GetProxy())
		if err != nil {
			return nil, err
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}

		clientConfig.HTTPClient = &http.Client{
			Transport: transport,
		}

		client = openai.NewClientWithConfig(clientConfig)
	} else {
		client = openai.NewClient(config.GetAiConfig().GetKey())
	}

	return &Engine{
		mode:         mode,
		config:       config,
		client:       client,
		execMessages: make([]openai.ChatCompletionMessage, 0),
		chatMessages: make([]openai.ChatCompletionMessage, 0),
		channel:      make(chan EngineOutput),
		running:      false,
	}, nil
}

func (e *Engine) SetMode(mode EngineMode) *Engine {
	e.mode = mode

	return e
}

func (e *Engine) GetMode() EngineMode {
	return e.mode
}

func (e *Engine) GetChannel() chan EngineOutput {
	return e.channel
}

func (e *Engine) Interrupt() *Engine {
	e.channel <- EngineOutput{
		content:    "[Interrupt]",
		last:       true,
		interrupt:  true,
		executable: false,
	}

	e.running = false

	return e
}

func (e *Engine) Clear() *Engine {
	if e.mode == ExecEngineMode {
		e.execMessages = []openai.ChatCompletionMessage{}
	} else {
		e.chatMessages = []openai.ChatCompletionMessage{}
	}

	return e
}

func (e *Engine) Reset() *Engine {
	e.execMessages = []openai.ChatCompletionMessage{}
	e.chatMessages = []openai.ChatCompletionMessage{}

	return e
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
				if e.mode == ExecEngineMode {
					if !strings.HasPrefix(output, noexec) && !strings.Contains(output, "\n") {
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
	if e.mode == ExecEngineMode {
		e.execMessages = append(e.execMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	} else {
		e.chatMessages = append(e.chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	}

	return e
}

func (e *Engine) appendAssistantMessage(content string) *Engine {
	if e.mode == ExecEngineMode {
		e.execMessages = append(e.execMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
	} else {
		e.chatMessages = append(e.chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
	}

	return e
}

func (e *Engine) prepareCompletionMessages() []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: e.prepareSystemPrompt(),
		},
	}

	if e.mode == ExecEngineMode {
		for _, m := range e.execMessages {
			messages = append(messages, m)
		}
	} else {
		for _, m := range e.chatMessages {
			messages = append(messages, m)
		}
	}

	return messages
}

func (e *Engine) prepareSystemPrompt() string {

	var bodyPart string
	if e.mode == ExecEngineMode {
		bodyPart = e.prepareSystemPromptExecPart()
	} else {
		bodyPart = e.prepareSystemPromptChatPart()
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		e.prepareSystemPromptCommonPart(),
		bodyPart,
		e.prepareSystemPromptSystemPart(),
	)
}

func (e *Engine) prepareSystemPromptCommonPart() string {
	return "You are Yo, an AI command line assistant running in a terminal, created by github.com/ekkinox.\n"
}

func (e *Engine) prepareSystemPromptExecPart() string {
	return "You will always generate one single command line (no \n, use ; and && instead) that I can run in my terminal.\n" +
		"Never add any explanation or details, even if I made a [chat] query previously.\n" +
		"If you absolutely cannot generate and reply only with a command line, " +
		fmt.Sprintf("reply by %s instead, and never add any explanation or details.\n\n", noexec) +
		"For example:\n" +
		"Me: List all files in /home\n" +
		"Yo: ls -l /home\n" +
		"Me: Now count them\n" +
		"Yo: ls -l /home | wc -l\n" +
		"Me: Does god exists ?\n" +
		"Yo: [noexec]\n" +
		"Me: Start docker\n" +
		"Yo: \n"
}

func (e *Engine) prepareSystemPromptChatPart() string {
	return "You will answer in the most helpful possible way, always rendered in markdown format.\n\n" +
		"For example:\n" +
		"Me: What is 2+2 ?\n" +
		"Yo: The answer for `2+2` is `4`\n" +
		"Me: What is a dog ?\n" +
		"Yo: \n"
}

func (e *Engine) prepareSystemPromptSystemPart() string {
	part := "My context: "
	if e.config.GetSystemConfig().GetOperatingSystem() != system.UnknownOperatingSystem {
		part += fmt.Sprintf("my operating system is %s, ", e.config.GetSystemConfig().GetOperatingSystem().String())
	}
	if e.config.GetSystemConfig().GetDistribution() != "" {
		part += fmt.Sprintf("my distribution is %s, ", e.config.GetSystemConfig().GetDistribution())
	}
	if e.config.GetSystemConfig().GetHomeDirectory() != "" {
		part += fmt.Sprintf("my home directory is %s, ", e.config.GetSystemConfig().GetHomeDirectory())
	}
	if e.config.GetSystemConfig().GetShell() != "" {
		part += fmt.Sprintf("my shell is %s, ", e.config.GetSystemConfig().GetShell())
	}
	if e.config.GetSystemConfig().GetShell() != "" {
		part += fmt.Sprintf("my editor is %s, ", e.config.GetSystemConfig().GetEditor())
	}
	part += "take this into account. "

	if e.config.GetUserConfig().GetPreferences() != "" {
		part += fmt.Sprintf("Also, %s.", e.config.GetUserConfig().GetPreferences())
	}

	return part
}
