package config

import "github.com/ilyakaznacheev/cleanenv"

type Service struct {
	Name                 string `env:"NAME" env-default:"shortener"`
	Host                 string `env:"HOST" env-default:"0.0.0.0"`
	Port                 int    `env:"PORT" env-default:"8080"`
	MaxGeneratorAttempts int    `env:"MAX_GENERATE_ATTEMPTS" env-default:"3"`
	InMemory             bool   `env:"IN_MEMORY_MODE" env-default:"false"`
}

type Postgres struct {
	Host     string `env:"HOST"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
	SSLMode  string `env:"SSL_MODE"`
	Port     int    `env:"PORT"`
	MaxConns int    `env:"MAX_CONNS" env-default:"20"`
	MinConns int    `env:"MIN_CONNS" env-default:"2"`
}

type Generator struct {
	Alphabet string `env:"ALPHABET" env-required:"true"`
	Len      int    `env:"LEN" env-required:"true"`
}

type Config struct {
	Postgres  Postgres  `env-prefix:"DB_"`
	Service   Service   `env-prefix:"SERVICE_"`
	Generator Generator `env-prefix:"GEN_"`
}

func Load() (Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
