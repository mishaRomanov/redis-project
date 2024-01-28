package config

import (
	"github.com/caarlos0/env/v10"
	"log"
)

// config struct
type Config struct {
	Port     string `env:"PORT"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB"`
}

//func that inits config and returns it
func Init() (Config, error) {
	var cfg Config
	//parsing env variables
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("error while parsing environmental variables: %w/n", err)
		return Config{}, err
	}
	return cfg, nil
}
