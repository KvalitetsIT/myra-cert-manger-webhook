package configs

type Configuration struct {
	Myra Myra `envPrefix:"MYRA_"`
	Hook Hook `envPrefix:"HOOK_"`
}
