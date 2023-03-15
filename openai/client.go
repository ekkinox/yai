package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ekkinox/hey/config"
)

const system_message = "You are a helpful assistant. You will generate zsh commands based on user input. Your response should contain ONLY the command and NO explanation. Do NOT ever use newlines to separate commands, instead use ; or &&. The current working directory is /home/jonathan"

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Client struct {
	Url   string
	Model string
	Key   string
}

func InitClient(cfg config.Config) (*Client, error) {
	if cfg.OpenAIKey == "" || cfg.OpenAIKey == config.Openai_Key_Placeholder {
		return nil, fmt.Errorf("openai api key is not defined")
	}

	return &Client{
		Url:   cfg.OpenAIUrl,
		Model: cfg.OpenAIModel,
		Key:   cfg.OpenAIKey,
	}, nil
}

func (c *Client) Send(input string) (string, error) {

	payload := fmt.Sprintf(
		`{"model":"%s","messages":[{"role":"system","content":"%s"},{"role":"user","content":"%s"}]}`,
		c.Model,
		system_message,
		input,
	)

	req, err := http.NewRequest(http.MethodPost, c.Url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Key))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("openai: could not make request")
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("openai: could not read response body")
		return "", err
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("openai: could not unmarshal JSON")
		return "", err
	}

	return strings.Trim(data.Choices[0].Message.Content, "\n"), nil
}
