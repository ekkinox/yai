package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ekkinox/hey/config"
	"github.com/ekkinox/hey/detect"
)

type Client struct {
	Messages []Message
	Config   config.Config
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Output struct {
	Executable bool
	Content    string
}

func InitClient(cfg config.Config) *Client {

	return &Client{[]Message{}, cfg}
}

func (c *Client) Reset() *Client {
	c.Messages = []Message{}

	return c
}

func (c *Client) Send(input string) (*Output, error) {

	payload := fmt.Sprintf(
		`{"model":"%s","temperature":0.2,"messages":%s}`,
		c.Config.OpenAI.Model,
		c.buildMessagesPayload(input),
	)

	req, err := http.NewRequest(http.MethodPost, c.Config.OpenAI.Url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.OpenAI.Key))

	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error during api call, code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	var data apiResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	c.Messages = append(c.Messages, Message{Role: "user", Content: input})
	c.Messages = append(c.Messages, Message{Role: "assistant", Content: data.Choices[0].Message.Content})

	output := strings.Trim(data.Choices[0].Message.Content, "\n")

	if strings.Contains(output, "OUTSPEECH") {
		return &Output{
			Executable: false,
			Content:    strings.ReplaceAll(output, "OUTSPEECH", "=>"),
		}, nil
	}

	return &Output{
		Executable: true,
		Content:    output,
	}, nil
}

func (c *Client) buildMessagesPayload(input string) string {

	messages := []Message{
		{
			Role:    "system",
			Content: c.buildSystemPrompt(),
		},
	}

	for _, m := range c.Messages {
		messages = append(messages, m)
	}

	messages = append(messages, Message{Role: "user", Content: input})

	payload, err := json.Marshal(messages)
	if err != nil {
		return ""
	}

	return string(payload)
}

func (c *Client) buildSystemPrompt() string {
	prompt := "You are Hey, a AI command line assistant running in a terminal, created by Jonathan VUILLEMIN (ekkinox). "
	prompt += "You will ALWAYS try to reply to the user input with ONLY a single line command line, WITHOUT any explanation and using separators like ; or &&. "

	prompt += "The context is the following: "
	if c.Config.System.OperatingSystem != detect.OS_other {
		prompt += fmt.Sprintf("the operating system is %s, ", c.Config.System.OperatingSystem)
	}

	if c.Config.System.Distribution != "" {
		prompt += fmt.Sprintf("the distribution is %s, ", c.Config.System.Distribution)
	}

	if c.Config.System.Shell != "" {
		prompt += fmt.Sprintf("the shell is %s, ", c.Config.System.Shell)
	}

	if c.Config.System.OperatingSystem != "" {
		prompt += fmt.Sprintf("the home directory is %s, ", c.Config.System.HomeDir)
	}
	prompt += "reply accordingly. "

	prompt += "If you cannot make an response containing ONLY a single command line, prefix your response with OUTSPEECH."

	return prompt
}
