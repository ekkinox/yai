package config

const openai_key = "OPENAI_KEY"
const openai_proxy = "OPENAI_PROXY"
const openai_temperature = "OPENAI_TEMPERATURE"

type AiConfig struct {
	key         string
	temperature float64
	proxy       string
}

func (c AiConfig) GetKey() string {
	return c.key
}

func (c AiConfig) GetTemperature() float64 {
	return c.temperature
}

func (c AiConfig) GetProxy() string {
	return c.proxy
}
