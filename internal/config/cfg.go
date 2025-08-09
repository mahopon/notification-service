package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Mail     *MailConfig
	Telegram *TGConfig
}

type MailConfig struct {
	Host     string
	Email    string
	Password string
}

type TGConfig struct {
	Key string
}

var Cfg *Config
var doOnce sync.Once

func Load() (*Config, error) {
	doOnce.Do(func() {
		_ = godotenv.Load()

		var mailConfig *MailConfig = nil
		var telegramConfig *TGConfig = nil

		if os.Getenv("MAIL_HOST") != "" && os.Getenv("MAIL_USER") != "" && os.Getenv("MAIL_PASSWORD") != "" {
			mailConfig = &MailConfig{
				Host:     os.Getenv("MAIL_HOST"),
				Email:    os.Getenv("MAIL_USER"),
				Password: os.Getenv("MAIL_PASSWORD"),
			}
		}

		if os.Getenv("TG_KEY") != "" {
			telegramConfig = &TGConfig{
				Key: os.Getenv("TG_KEY"),
			}
		}

		Cfg = &Config{
			Mail:     mailConfig,
			Telegram: telegramConfig,
		}
	})
	return Cfg, nil
}
