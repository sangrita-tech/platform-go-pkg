package logger

type Configs struct {
	Level      string            `yaml:"level" env:"LEVEL" env-default:"info"`
	DevMode    bool              `yaml:"devMode" env:"DEV_MODE" env-default:"false"`
	BaseFields map[string]string `yaml:"baseFields" env:"BASE_FIELDS"`
}
