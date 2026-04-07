package configs

type Myra struct {
	Api Api `envPrefix:"API_"`
}

type Api struct {
	URL    string `env:"URL"  envDefault:"http://localhost:8080"`
	Key    string `env:"KEY"`
	Secret string `env:"SECRET"`
	Token  string `env:"TOKEN"`
}
