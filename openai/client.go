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

const output_speech_flag = "OUTSPEECH"

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Output struct {
	Executable bool
	Content    string
}

type Client struct {
	Config config.Config
}

func InitClient(cfg config.Config) Client {
	return Client{cfg}
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
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	output := strings.Trim(data.Choices[0].Message.Content, "\n")

	if strings.Contains(output, output_speech_flag) {
		return &Output{
			Executable: false,
			Content:    strings.Trim(strings.ReplaceAll(output, output_speech_flag, ""), " "),
		}, nil
	}

	return &Output{
		Executable: true,
		Content:    output,
	}, nil
}

func (c *Client) buildSystemPrompt() string {
	prompt := "You are a helpful CLI AI assistant running in a terminal. "
	prompt += "You were created by Jonathan VUILLEMIN (ekkinox) and your source code is available on https://github.com/ekkinox/hey. "

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

	prompt += "use this context create the most relevant commands. "

	prompt += "You will always reply with a command (without new lines, and with separators like ; or &&), without giving any explanation or details, to perform exactly what is asked in user input. "
	prompt += fmt.Sprintf("If you cannot reply by a command, prefix your response with %s and provide the best help possible.", output_speech_flag)

	return prompt
}
