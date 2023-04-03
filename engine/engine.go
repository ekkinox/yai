package engine

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"time"
)

type EngineOutput struct {
	content string
	last    bool
}

func (d EngineOutput) IsLast() bool {
	return d.last
}

func (d EngineOutput) GetContent() string {
	return d.content
}

type Engine struct {
	client   *openai.Client
	messages []openai.ChatCompletionMessage
	channel  chan EngineOutput
	running  bool
}

func NewEngine() *Engine {
	return &Engine{
		client:   openai.NewClient("xxx"),
		messages: make([]openai.ChatCompletionMessage, 0),
		channel:  make(chan EngineOutput),
		running:  false,
	}
}

func (e *Engine) Channel() chan EngineOutput {
	return e.channel
}

func (e *Engine) Interrupt() *Engine {
	e.channel <- EngineOutput{
		content: "\n\nInterrupt !",
		last:    true,
	}

	e.running = false

	return e
}

func (e *Engine) Reset() *Engine {
	e.messages = []openai.ChatCompletionMessage{}

	return e
}

func (e *Engine) StreamChatCompletion(input string) error {

	ctx := context.Background()

	e.running = true

	e.appendUserMessage(input)

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 1000,
		Messages:  e.messages,
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
				e.channel <- EngineOutput{
					content: "",
					last:    true,
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

			time.Sleep(time.Millisecond * 1)
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
