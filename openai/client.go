package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ekkinox/hey/config"
	"github.com/ekkinox/hey/detect"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MessageList struct {
	Messages []Message `json:"messages"`
}

type Client struct {
	History MessageList
	Config  config.Config
}

type Output struct {
	Executable bool
	Content    string
}

func InitClient(cfg config.Config) Client {

	return Client{MessageList{}, cfg}
}

func (c Client) Reset() Client {
	c.History = MessageList{}

	return c
}

func (c Client) Send(input string) (*Output, error) {

	payload := fmt.Sprintf(
		`{"model":"%s","messages":[{"role":"system","content":"%s"},{"role":"user","content":"%s"}]}`,
		c.Config.OpenAI.Model,
		c.buildSystemPrompt(),
		input,
	)

	req, err := http.NewRequest(http.MethodPost, c.Config.OpenAI.Url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.OpenAI.Key))

	client := &http.Client{}
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

	output := strings.Trim(data.Choices[0].Message.Content, "\n")

	if output[0:1] == "[" && output[len(output)-1:] == "]" {

		output = strings.TrimPrefix(output, "[")
		output = strings.TrimSuffix(output, "]")

		return &Output{
			Executable: true,
			Content:    output,
		}, nil
	}

	return &Output{
		Executable: false,
		Content:    output,
	}, nil
}

func (c *Client) buildSystemPrompt() string {
	prompt := "You are Hey, a AI command line assistant running in a terminal, created by Jonathan VUILLEMIN (ekkinox). "
	prompt += "You will ALWAYS try to answer to the user input with ONLY a single line command line, WITHOUT any explanation, surrounded by [], and using separators like ; or &&. "

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

	prompt += "If you really cannot answer with a command line, reply with ONLY your answer EVEN if it is not a command line."

	return prompt
}
