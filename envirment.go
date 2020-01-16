package xin

type Mode int

// Predefine mode const
const (
	DevMode Mode = iota
	TestMode
	ReleaseMode
)

type Envirment interface {
	Mode() Mode
}

type EnvSetting struct {
	config *Config
}

func (e *EnvSetting) Mode() Mode {
	env := e.config.Viper().GetString("env")
	switch env {
	case "dev", "", "debug":
		return DevMode
	case "test":
		return TestMode
	case "release":
		return ReleaseMode
	default:
		//TODO:  remove panic
		panic("Unknown env string")
	}
}

func NewEnvSetting(config *Config) *EnvSetting {
	return &EnvSetting{
		config: config,
	}
}
