package config

import (
	"fmt"
	"github.com/ekkinox/yo/context"
	"github.com/spf13/viper"
	"log"
	"strings"
)

const openai_key = "OPENAI_KEY"
const openai_temperature = "OPENAI_TEMPERATURE"
const user_default_mode = "USER_DEFAULT_MODE"
const user_context = "USER_CONTEXT"

type ConfigOutput struct{}

type ConfigOpenAI struct {
	key         string
	temperature float64
}

func (co ConfigOpenAI) GetKey() string {
	return co.key
}

func (co ConfigOpenAI) GetTemperature() float64 {
	return co.temperature
}

type ConfigUserPreferences struct {
	defaultMode string
	context     string
}

func (cu ConfigUserPreferences) GetDefaultMode() string {
	return cu.defaultMode
}

func (cu ConfigUserPreferences) GetContext() string {
	return cu.context
}

type Config struct {
	openAI          ConfigOpenAI
	userPreferences ConfigUserPreferences
	context         *context.Context
}

func (c *Config) GetOpenAI() ConfigOpenAI {
	return c.openAI
}

func (c *Config) GetUserPreferences() ConfigUserPreferences {
	return c.userPreferences
}

func (c *Config) GetContext() *context.Context {
	return c.context
}

func NewConfig() (*Config, error) {

	context := context.NewContextAnalyzer().Analyse()

	viper.SetConfigName(strings.ToLower(context.GetAppName()))
	viper.AddConfigPath(fmt.Sprintf("%s/.config/", context.GetHomeDirectory()))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		openAI: ConfigOpenAI{
			key:         viper.GetString(openai_key),
			temperature: viper.GetFloat64(openai_temperature),
		},
		userPreferences: ConfigUserPreferences{
			defaultMode: viper.GetString(user_default_mode),
			context:     viper.GetString(user_context),
		},
		context: context,
	}, nil
}

func WriteConfig(key string) (*Config, error) {

	context := context.NewContextAnalyzer().Analyse()

	viper.Set(openai_key, key)
	viper.SetDefault(openai_temperature, 0.2)
	viper.SetDefault(user_default_mode, "chat")
	viper.SetDefault(user_context, "")

	err := viper.SafeWriteConfigAs(context.GetConfigFile())
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	return NewConfig()
}
