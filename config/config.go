package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const openai_url = "OPENAI_URL"
const openai_key = "OPENAI_KEY"
const openai_model = "OPENAI_MODEL"

const Openai_Key_Placeholder = "{REPLACE-ME}"

type Config struct {
	OpenAIUrl   string
	OpenAIKey   string
	OpenAIModel string
}

func InitConfig(cfgFile string) Config {

	viper.SetDefault(openai_url, "https://api.openai.com/v1/chat/completions")
	viper.SetDefault(openai_model, "gpt-3.5-turbo")

	homeDir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("hey")
		viper.AddConfigPath(fmt.Sprintf("%s/.config/", homeDir))
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			newCfgFile := fmt.Sprintf("%s/.config/hey.yaml", homeDir)
			fmt.Printf("Creating config file in: %s, please update your %s before running again.", newCfgFile, openai_key)
			viper.Set(openai_key, Openai_Key_Placeholder)
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
		OpenAIUrl:   viper.GetString(openai_url),
		OpenAIKey:   viper.GetString(openai_key),
		OpenAIModel: viper.GetString(openai_model),
	}
}
