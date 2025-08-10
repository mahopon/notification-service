package infra

import (
	"fmt"
	"log"

	"github.com/mahopon/notification-service/internal/config"
	bolt "go.etcd.io/bbolt"
)

type DatabaseConnection interface {
	Set(bucket, key, value string) error
	Get(bucket, key string) (string, error)
}

type DatabaseConfig struct {
	Database *bolt.DB
}

func NewDatabaseConfig(cfg *config.DBConfig) *DatabaseConfig {
	dbLoc := cfg.Location
	db, err := bolt.Open(dbLoc, 0600, nil)
	if err != nil {
		log.Fatalf("DB not initialised: %v", err)
	}
	returnConfig := &DatabaseConfig{
		Database: db,
	}
	returnConfig.initSchemas()
	return returnConfig
}

func CloseDatabaseConfig(dbConfig *DatabaseConfig) {
	dbConfig.Database.Close()
}

func (db *DatabaseConfig) initSchemas() {
	err := db.Database.Update(func(tx *bolt.Tx) error {
		buckets := []string{"user_chat"}
		for _, b := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return fmt.Errorf("create bucket %q: %w", b, err)
			}
		}
		return nil
	})
	if err != nil {
		db.Database.Close()
		log.Fatal("ERROR: Unable to create schemas")
	}
}

func (db *DatabaseConfig) Set(bucket, key, value string) error {
	return db.Database.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(bucket))
		return b.Put([]byte(key), []byte(value))
	})
}

func (db *DatabaseConfig) Get(bucket, key string) (string, error) {
	var value string
	err := db.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %q not found", bucket)
		}
		v := b.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("key %q not found", key)
		}
		value = string(v)
		return nil
	})
	return value, err
}
