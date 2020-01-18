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
	config EnvConfig
}

func (e *EnvSetting) Mode() Mode {
	env := e.config.Env()
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

type EnvConfig interface {
	Env() string
}

func NewEnvSetting(config EnvConfig) *EnvSetting {
	return &EnvSetting{
		config: config,
	}
}
