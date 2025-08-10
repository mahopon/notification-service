package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Mail     *MailConfig
	Telegram *TGConfig
	Database *DBConfig
}

type MailConfig struct {
	Host     string
	Email    string
	Password string
}

type TGConfig struct {
	Key   string
	Debug bool
}

type DBConfig struct {
	Location string
}

var Cfg *Config
var doOnce sync.Once

func Load(db_debug bool) (*Config, error) {
	doOnce.Do(func() {
		_ = godotenv.Load()

		if os.Getenv("DB_LOC") == "" {
			log.Fatal("No DB_LOC in env file")
		}

		cwd, _ := os.Getwd()
		dbConfig := &DBConfig{
			Location: cwd + os.Getenv("DB_LOC"),
		}

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
				Key:   os.Getenv("TG_KEY"),
				Debug: db_debug,
			}
		}

		Cfg = &Config{
			Database: dbConfig,
			Mail:     mailConfig,
			Telegram: telegramConfig,
		}
	})
	return Cfg, nil
}
