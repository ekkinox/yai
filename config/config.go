package config

import (
	"fmt"
	"github.com/ekkinox/yo/system"
	"github.com/spf13/viper"
	"strings"
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
			temperature: viper.GetFloat64(openai_temperature),
			proxy:       viper.GetString(openai_proxy),
		},
		user: UserConfig{
			defaultPromptMode: viper.GetString(user_default_prompt_mode),
			preferences:       viper.GetString(user_preferences),
		},
		system: system,
	}, nil
}

func WriteConfig(key string) (*Config, error) {

	system := system.Analyse()

	viper.Set(openai_key, key)
	viper.SetDefault(openai_temperature, 0.2)
	viper.SetDefault(openai_proxy, "")
	viper.SetDefault(user_default_prompt_mode, "exec")
	viper.SetDefault(user_preferences, "")

	err := viper.SafeWriteConfigAs(system.GetConfigFile())
	if err != nil {
		return nil, err
	}

	return NewConfig()
}
