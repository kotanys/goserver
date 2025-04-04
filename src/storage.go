package main

import (
	"fmt"
)

type StorageKey string
type StorageValue []byte

type Storage struct {
	data map[StorageKey]StorageValue
}

func NewStorage() *Storage {
	return &Storage{data: make(map[StorageKey]StorageValue)}
}

func (s *Storage) Get(key StorageKey) (StorageValue, bool) {
	value, exist := s.data[key]
	return value, exist
}

func (s *Storage) Put(key StorageKey, value StorageValue) {
	s.data[key] = value
}

func (s *Storage) Delete(key StorageKey) error {
	_, exist := s.data[key]
	if exist {
		delete(s.data, key)
		return nil
	} else {
		return fmt.Errorf("no key %v", key)
	}
}
