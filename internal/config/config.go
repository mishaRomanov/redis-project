package config

import (
	"github.com/caarlos0/env/v10"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// config struct
type Config struct {
	Port     string `env:"PORT"`
	Password string `env:"PASSWORD" envDefault:""`
	DB       int    `env:"DB" envDefault:"0"`
}

// func that inits config and returns it
func Init() (Config, error) {
	var cfg Config
	//parsing env variables
	err := env.Parse(&cfg)
	if err != nil {
		logrus.Fatalf("error while parsing environmental variables: %w/n", err)
		return Config{}, err
	}
	return cfg, nil
}

// Func that creates a middleware for JWT authenticationq
func TokenConfig() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("key_for_client"),
		ErrorHandler: func(ctx echo.Context, err error) error {
			logrus.Infof("Unauthorized request")
			//returns status unauthorized
			return ctx.String(401, "Unauthorized. Invalid token")
		},
	})
}
