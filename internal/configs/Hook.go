package configs

type Hook struct {
	GroupName string `env:"GROUP_NAME" envDefault:"myra-dns-solver"`
	Port      uint16 `env:"Port" envDefault:"8000"`
}
