package infra

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

type DatabaseConnection interface {
	Set(bucket, key, value string) error
	Get(bucket, key string) (string, error)
}

type DatabaseConfig struct {
	Database *bolt.DB
}

func NewDatabaseConfig(dbLoc string) *DatabaseConfig {
	db, err := bolt.Open(dbLoc, 0600, nil)
	if err != nil {
		log.Fatalf("DB not initialised: %v", err)
	}
	return &DatabaseConfig{
		Database: db,
	}
}

func (db *DatabaseConfig) Set(bucket, key, value string) error {
	return nil
}

func (db *DatabaseConfig) Get(bucket, key string) (string, error) {
	return "", nil
}
