package config

import (
	"fmt"
	"github.com/ekkinox/hey/detect"
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
}

type OpenAIConfig struct {
	Url   string
	Key   string
	Model string
}

func InitConfig(cfgFile string) Config {

	viper.SetDefault(env_openai_url, "https://api.openai.com/v1/chat/completions")
	viper.SetDefault(env_openai_model, "gpt-4.0-turbo")

	homeDir := detect.DetectHomeDir()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("hey")
		viper.AddConfigPath(fmt.Sprintf("%s/.config/", homeDir))
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			newCfgFile := fmt.Sprintf("%s/.config/hey.yaml", homeDir)
			fmt.Printf("Creating config file in: %s, please update your %s before running again.", newCfgFile, env_openai_key)
			viper.Set(env_openai_key, Openai_Key_Placeholder)
			err = viper.SafeWriteConfigAs(newCfgFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)
		} else {
			fmt.Println("Can't read config:", err)
			os.Exit(1)
		}
	}

	return Config{
		System: SystemConfig{
			OperatingSystem: detect.DetectOperatingSystem(),
			Distribution:    detect.DetectDistribution(),
			Shell:           detect.DetectShell(),
			HomeDir:         homeDir,
		},
		OpenAI: OpenAIConfig{
			Url:   viper.GetString(env_openai_url),
			Key:   viper.GetString(env_openai_key),
			Model: viper.GetString(env_openai_model),
		},
	}
}
