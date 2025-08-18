package services

import (
	"errors"
)

type MockDB struct {
	store      map[string]string
	ShouldFail bool
}

func NewMockDB() *MockDB {
	return &MockDB{
		store: make(map[string]string),
	}
}

func (db *MockDB) Get(bucket, key string) (string, error) {
	if db.ShouldFail {
		return "", errors.New("db get failed")
	}
	val, ok := db.store[key]
	if !ok {
		return "", errors.New("not found")
	}
	return val, nil
}

func (db *MockDB) Set(bucket, key, value string) error {
	if db.ShouldFail {
		return errors.New("db set failed")
	}
	db.store[key] = value
	return nil
}
