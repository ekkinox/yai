package config

import (
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/ekkinox/yai/system"
	"github.com/spf13/viper"
)

type Config struct {
	ai     AiConfig
	user   UserConfig
	system *system.Analysis
}

func (c *Config) GetAiConfig() AiConfig {
	return c.ai
}

func (c *Config) GetUserConfig() UserConfig {
	return c.user
}

func (c *Config) GetSystemConfig() *system.Analysis {
	return c.system
}

func NewConfig() (*Config, error) {

	system := system.Analyse()

	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	viper.AddConfigPath(fmt.Sprintf("%s/.config/", system.GetHomeDirectory()))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		ai: AiConfig{
			key:         viper.GetString(openai_key),
			model:       viper.GetString(openai_model),
			proxy:       viper.GetString(openai_proxy),
			temperature: viper.GetFloat64(openai_temperature),
			maxTokens:   viper.GetInt(openai_max_tokens),
		},
		user: UserConfig{
			defaultPromptMode: viper.GetString(user_default_prompt_mode),
			preferences:       viper.GetString(user_preferences),
		},
		system: system,
	}, nil
}

func WriteConfig(key string, write bool) (*Config, error) {
	system := system.Analyse()
	// ai defaults
	viper.Set(openai_key, key)
	viper.Set(openai_model, openai.GPT3Dot5Turbo)
	viper.SetDefault(openai_proxy, "")
	viper.SetDefault(openai_temperature, 0.2)
	viper.SetDefault(openai_max_tokens, 1000)
	// user defaults
	viper.SetDefault(user_default_prompt_mode, "exec")
	viper.SetDefault(user_preferences, "")

	if write {
		err := viper.SafeWriteConfigAs(system.GetConfigFile())
		if err != nil {
			return nil, err
		}
	}

	return NewConfig()
}
