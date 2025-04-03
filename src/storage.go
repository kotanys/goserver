package main

import (
	"fmt"
)

type StorageKey string
type StorageValue []byte

type Storage struct {
	data map[StorageKey]StorageValue
}

func (s *Storage) Get(key StorageKey) (StorageValue, error) {
	value, exist := s.data[key]
	if exist {
		return value, nil
	} else {
		return nil, fmt.Errorf("no key %v", key)
	}
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
