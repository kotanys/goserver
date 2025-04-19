package main

import "fmt"

type StorageKey string
type StorageValue []byte

type Storage struct {
	persistent bool
	log        *PersistentLogger
	data       map[StorageKey]StorageValue
}

func NewStorage(log *PersistentLogger, persistent bool) *Storage {
	st := &Storage{data: make(map[StorageKey]StorageValue), log: log, persistent: false}
	if persistent {
		opChan := make(chan LogOperation, 10)
		fmt.Println("started recovering")
		go log.Replay(opChan)
		for op := range opChan {
			op.Recover(st)
		}
		fmt.Printf("ended recovering: recovered %v entries\n", len(st.data))
	}
	st.persistent = true
	return st
}

func (s *Storage) Get(key StorageKey) (StorageValue, bool) {
	value, exist := s.data[key]
	return value, exist
}

func (s *Storage) Put(key StorageKey, value StorageValue) {
	s.data[key] = value
	if s.persistent {
		op := PutOperation{K: key, V: value}
		s.log.Append(op)
	}
}

func (s *Storage) Delete(key StorageKey) {
	_, exist := s.data[key]
	if !exist {
		return
	}
	delete(s.data, key)
	if s.persistent {
		op := DeleteOperation{K: key}
		s.log.Append(op)
	}
}
