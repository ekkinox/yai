package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/ekkinox/yai/config"
	"github.com/ekkinox/yai/system"

	"github.com/sashabaranov/go-openai"
)

const noexec = "[noexec]"

type Engine struct {
	mode         EngineMode
	config       *config.Config
	client       *openai.Client
	execMessages []openai.ChatCompletionMessage
	chatMessages []openai.ChatCompletionMessage
	channel      chan EngineChatStreamOutput
	pipe         string
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
		channel:      make(chan EngineChatStreamOutput),
		pipe:         "",
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

func (e *Engine) GetChannel() chan EngineChatStreamOutput {
	return e.channel
}

func (e *Engine) SetPipe(pipe string) *Engine {
	e.pipe = pipe

	return e
}

func (e *Engine) Interrupt() *Engine {
	e.channel <- EngineChatStreamOutput{
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

func (e *Engine) ExecCompletion(input string) (*EngineExecOutput, error) {
	ctx := context.Background()

	e.running = true

	e.appendUserMessage(input)

	resp, err := e.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:     e.config.GetAiConfig().GetModel(),
			MaxTokens: e.config.GetAiConfig().GetMaxTokens(),
			Messages:  e.prepareCompletionMessages(),
		},
	)
	if err != nil {
		return nil, err
	}

	content := resp.Choices[0].Message.Content
	e.appendAssistantMessage(content)

	var output EngineExecOutput
	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		re := regexp.MustCompile(`\{.*?\}`)
		match := re.FindString(content)
		if match != "" {
			err = json.Unmarshal([]byte(match), &output)
			if err != nil {
				return nil, err
			}
		} else {
			output = EngineExecOutput{
				Command:     "",
				Explanation: content,
				Executable:  false,
			}
		}
	}

	return &output, nil
}

func (e *Engine) ChatStreamCompletion(input string) error {
	ctx := context.Background()

	e.running = true

	e.appendUserMessage(input)

	req := openai.ChatCompletionRequest{
		Model:     e.config.GetAiConfig().GetModel(),
		MaxTokens: e.config.GetAiConfig().GetMaxTokens(),
		Messages:  e.prepareCompletionMessages(),
		Stream:    true,
	}

	stream, err := e.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
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

				e.channel <- EngineChatStreamOutput{
					content:    "",
					last:       true,
					executable: executable,
				}
				e.running = false
				e.appendAssistantMessage(output)

				return nil
			}

			if err != nil {
				e.running = false
				return err
			}

			delta := resp.Choices[0].Delta.Content

			output += delta

			e.channel <- EngineChatStreamOutput{
				content: delta,
				last:    false,
			}

			// time.Sleep(time.Microsecond * 100)
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

	if e.pipe != "" {
		messages = append(
			messages,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: e.preparePipePrompt(),
			},
		)
	}

	if e.mode == ExecEngineMode {
		messages = append(messages, e.execMessages...)
	} else {
		messages = append(messages, e.chatMessages...)
	}

	return messages
}

func (e *Engine) preparePipePrompt() string {
	return fmt.Sprintf("I will work on the following input: %s", e.pipe)
}

func (e *Engine) prepareSystemPrompt() string {
	var bodyPart string
	if e.mode == ExecEngineMode {
		bodyPart = e.prepareSystemPromptExecPart()
	} else {
		bodyPart = e.prepareSystemPromptChatPart()
	}

	return fmt.Sprintf("%s\n%s", bodyPart, e.prepareSystemPromptContextPart())
}

func (e *Engine) prepareSystemPromptExecPart() string {
	return "Your are Yai, a powerful terminal assistant generating a JSON containing a command line for my input.\n" +
		"You will always reply using the following json structure: {\"cmd\":\"the command\", \"exp\": \"some explanation\", \"exec\": true}.\n" +
		"Your answer will always only contain the json structure, never add any advice or supplementary detail or information, even if I asked the same question before.\n" +
		"The field cmd will contain a single line command (don't use new lines, use separators like && and ; instead).\n" +
		"The field exp will contain an short explanation of the command if you managed to generate an executable command, otherwise it will contain the reason of your failure.\n" +
		"The field exec will contain true if you managed to generate an executable command, false otherwise." +
		"\n" +
		"Examples:\n" +
		"Me: list all files in my home dir\n" +
		"Yai: {\"cmd\":\"ls ~\", \"exp\": \"list all files in your home dir\", \"exec\\: true}\n" +
		"Me: list all pods of all namespaces\n" +
		"Yai: {\"cmd\":\"kubectl get pods --all-namespaces\", \"exp\": \"list pods form all k8s namespaces\", \"exec\": true}\n" +
		"Me: how are you ?\n" +
		"Yai: {\"cmd\":\"\", \"exp\": \"I'm good thanks but I cannot generate a command for this. Use the chat mode to discuss.\", \"exec\": false}"
}

func (e *Engine) prepareSystemPromptChatPart() string {
	return "You are Yai a powerful terminal assistant created by github.com/ekkinox.\n" +
		"You will answer in the most helpful possible way.\n" +
		"Always format your answer in markdown format.\n\n" +
		"For example:\n" +
		"Me: What is 2+2 ?\n" +
		"Yai: The answer for `2+2` is `4`\n" +
		"Me: +2 again ?\n" +
		"Yai: The answer is `6`\n"
}

func (e *Engine) prepareSystemPromptContextPart() string {
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
