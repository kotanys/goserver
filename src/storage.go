package main

import (
	"fmt"
)

type StorageKey string
type StorageValue []byte

type Storage struct {
	log  *PersistentLogger
	data map[StorageKey]StorageValue
}

func NewStorage(log *PersistentLogger) *Storage {
	return &Storage{data: make(map[StorageKey]StorageValue), log: log}
}

func (s *Storage) Get(key StorageKey) (StorageValue, bool) {
	value, exist := s.data[key]
	return value, exist
}

func (s *Storage) Put(key StorageKey, value StorageValue) {
	s.data[key] = value
	if s.log != nil {
		op := PutOperation{K: key, V: value}
		s.log.AppendPut(op)
	}
}

func (s *Storage) Delete(key StorageKey) error {
	_, exist := s.data[key]
	if exist {
		delete(s.data, key)
		if s.log != nil {
			op := DeleteOperation{K: key}
			s.log.AppendDelete(op)
		}
		return nil
	} else {
		return fmt.Errorf("no key %v", key)
	}
}
