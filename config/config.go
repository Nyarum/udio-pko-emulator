package config

type Config struct {
	IsDebug bool
}

func NewConfig(isDebug bool) Config {
	return Config{
		IsDebug: isDebug,
	}
}
