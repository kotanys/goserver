package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PersistentLogger struct {
	file *os.File
}

type PutOperation struct {
	K StorageKey   `json:"key"`
	V StorageValue `json:"value"`
}
type DeleteOperation struct {
	K StorageKey `json:"key"`
}

func NewPeristentLogger(fileName string) (*PersistentLogger, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	log := &PersistentLogger{file: file}
	return log, nil
}

func (log *PersistentLogger) Close() {
	if err := log.file.Close(); err != nil {
		panic(err)
	}
}

func (log *PersistentLogger) AppendPut(op PutOperation) {
	bytes, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}
	if _, err = fmt.Fprint(log.file, "PUT "); err != nil {
		panic(err)
	}
	if _, err = log.file.Write(bytes); err != nil {
		panic(err)
	}
	if _, err = fmt.Fprintln(log.file); err != nil {
		panic(err)
	}
}
func (log *PersistentLogger) AppendDelete(op DeleteOperation) {
	bytes, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}
	if _, err = fmt.Fprint(log.file, "DELETE "); err != nil {
		panic(err)
	}
	if _, err = log.file.Write(bytes); err != nil {
		panic(err)
	}
	if _, err = fmt.Fprintln(log.file); err != nil {
		panic(err)
	}
}
