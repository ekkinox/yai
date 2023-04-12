package config

const (
	user_default_prompt_mode = "USER_DEFAULT_PROMPT_MODE"
	user_preferences         = "USER_PREFERENCES"
)

type UserConfig struct {
	defaultPromptMode string
	preferences       string
}

func (c UserConfig) GetDefaultPromptMode() string {
	return c.defaultPromptMode
}

func (c UserConfig) GetPreferences() string {
	return c.preferences
}
