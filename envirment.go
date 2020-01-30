package xin

//Mode app mode
type Mode int

// Predefine mode const
const (
	DevMode Mode = iota
	TestMode
	ReleaseMode
)

//Envirment envirment interface
type Envirment interface {
	//Return working mode
	Mode() Mode
}

//EnvConfig a config interface for envirment
type EnvConfig interface {
	//get working mode string
	Env() string
}

//EnvSetting implement Envirment interface
type EnvSetting struct {
	config EnvConfig
}

//Mode return app working mode
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

//NewEnvSetting generate new EnvSetting
func NewEnvSetting(config EnvConfig) *EnvSetting {
	return &EnvSetting{
		config: config,
	}
}
