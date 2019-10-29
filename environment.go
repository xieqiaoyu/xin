package xin

// Predefine mode const
const (
	Dev = iota
	Test
	Release
)

//Mode Get current app env mode in config
func Mode() int {
	env := Config().GetString("env")
	switch env {
	case "dev", "", "debug":
		return Dev
	case "test":
		return Test
	case "release":
		return Release
	}
	//TODO: remove panic
	panic("Unsupport environment: " + env)
}
