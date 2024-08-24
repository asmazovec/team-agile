package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// AppConfig application runtime configuration.
type AppConfig struct {
	AppShutdownTimeout time.Duration `env:"APP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
	HTTPPrimaryServer  ServerConfig  `envPrefix:"HTTP_"`
}

// ServerConfig HTTP server config.
type ServerConfig struct {
	Address           string        `env:"ADDRESS" envDefault:":8080"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"10s"`
}

// Origin default value will never break builder.
type Origin func(*AppConfig) error

// MustRead reads application config from a set of origins in presented order.
// On field collisions it will use the latest value.
func MustRead(origins ...Origin) AppConfig {
	cfg := AppConfig{}
	for _, opt := range origins {
		if opt == nil {
			continue
		}
		err := opt(&cfg)
		if err != nil {
			panic(err)
		}
	}
	return cfg
}

// FromEnv reads values from an env variables.
func FromEnv(path string) Origin {
	return func(cfg *AppConfig) error {
		if path != "" {
			err := godotenv.Load(path)
			if err != nil {
				return err
			}
		}
		if cfg == nil {
			return nil
		}
		return env.Parse(cfg)
	}
}
