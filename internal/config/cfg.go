package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Mail MailConfig
}

type MailConfig struct {
	Host     string
	Email    string
	Password string
}

var Cfg *Config
var doOnce sync.Once

func Load() (*Config, error) {
	doOnce.Do(func() {
		_ = godotenv.Load()
		Cfg = &Config{
			Mail: MailConfig{
				Host:     os.Getenv("MAIL_HOST"),
				Email:    os.Getenv("MAIL_USER"),
				Password: os.Getenv("MAIL_PASSWORD"),
			},
		}
	})
	return Cfg, nil
}
