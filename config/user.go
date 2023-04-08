package config

const user_default_mode = "USER_DEFAULT_MODE"
const user_preferences = "USER_PREFERENCES"

type UserConfig struct {
	defaultMode string
	preferences string
}

func (c UserConfig) GetDefaultMode() string {
	return c.defaultMode
}

func (c UserConfig) GetPreferences() string {
	return c.preferences
}
