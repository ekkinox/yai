package config

import (
	"fmt"
	"github.com/ekkinox/hey/detect"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"os"
)

const env_openai_url = "OPENAI_URL"
const env_openai_key = "OPENAI_KEY"
const env_openai_model = "OPENAI_MODEL"

const Openai_Key_Placeholder = "{REPLACE-ME}"

type Config struct {
	System SystemConfig
	OpenAI OpenAIConfig
}

type SystemConfig struct {
	OperatingSystem string
	Distribution    string
	Shell           string
	HomeDir         string
	Username        string
}

type OpenAIConfig struct {
	Url   string
	Key   string
	Model string
}

func InitConfig() Config {

	homeDir := detect.DetectHomeDir()
	username := detect.DetectUsername()

	viper.SetDefault(env_openai_url, "https://api.openai.com/v1/chat/completions")
	viper.SetDefault(env_openai_model, "gpt-3.5-turbo")

	viper.SetConfigName("hey")
	viper.AddConfigPath(fmt.Sprintf("%s/.config/", homeDir))

	if err := viper.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			fmt.Printf("Hey %s!\nApparently it is the first time you ask for my help, and for this I will need an OpenAI API key.\n", username)
			prompt := promptui.Prompt{
				Label: "OpenAI API key",
			}
			key, err := prompt.Run()
			if err != nil {
				color.HiRed("Cannot read key.", err)
				os.Exit(1)
			}

			viper.Set(env_openai_key, key)

			newCfgFile := fmt.Sprintf("%s/.config/hey.json", homeDir)
			fmt.Printf("Creating config file in: %s.\n\n", newCfgFile)

			err = viper.SafeWriteConfigAs(newCfgFile)
			if err != nil {
				color.HiRed("Cannot save config file.", err)
				os.Exit(1)
			}
		} else {
			color.HiRed("Cannot read config.", err)
			os.Exit(1)
		}
	}

	return Config{
		System: SystemConfig{
			OperatingSystem: detect.DetectOperatingSystem(),
			Distribution:    detect.DetectDistribution(),
			Shell:           detect.DetectShell(),
			HomeDir:         homeDir,
			Username:        username,
		},
		OpenAI: OpenAIConfig{
			Url:   viper.GetString(env_openai_url),
			Key:   viper.GetString(env_openai_key),
			Model: viper.GetString(env_openai_model),
		},
	}
}
