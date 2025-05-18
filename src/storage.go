package main

import "fmt"

type StorageKey string
type StorageValue []byte

type Storage struct {
	log  *PersistentLogger
	data map[StorageKey]StorageValue
	cfg  *StorageConfig
}

func NewStorage(log *PersistentLogger, cfg *StorageConfig) *Storage {
	st := &Storage{data: make(map[StorageKey]StorageValue), log: log, cfg: cfg}
	if log != nil {
		opChan := make(chan LogOperation, 10)
		fmt.Println("started recovering")
		go log.Replay(opChan)
		for op := range opChan {
			op.Recover(st)
		}
		fmt.Printf("ended recovering: recovered %v entries\n", len(st.data))
	}
	return st
}

func (s *Storage) Get(key StorageKey) (StorageValue, bool) {
	s.cfg.Update()
	value, exist := s.data[key]
	return value, exist
}

func (s *Storage) Put(key StorageKey, value StorageValue) {
	s.cfg.Update()
	s.data[key] = value
	if s.cfg.Persistent {
		op := PutOperation{K: key, V: value}
		s.log.Append(op)
	}
}

func (s *Storage) Delete(key StorageKey) {
	s.cfg.Update()
	_, exist := s.data[key]
	if !exist {
		return
	}
	delete(s.data, key)
	if s.cfg.Persistent {
		op := DeleteOperation{K: key}
		s.log.Append(op)
	}
}
