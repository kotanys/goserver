package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type PersistentLogger struct {
	file *os.File
}

type LogOperation interface {
	GetBytes() ([]byte, error)
	Recover(s *Storage) error
}

type PutOperation struct {
	K StorageKey   `json:"key"`
	V StorageValue `json:"value"`
}

func (op PutOperation) GetBytes() ([]byte, error) {
	jsonBytes, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, 0, len("PUT ")+len(jsonBytes)+1)
	bytes = append(bytes, []byte("PUT ")...)
	bytes = append(bytes, jsonBytes...)
	bytes = append(bytes, []byte("\n")...)
	return bytes, nil
}

func (op PutOperation) Recover(s *Storage) error {
	s.Put(op.K, op.V)
	return nil
}

type DeleteOperation struct {
	K StorageKey `json:"key"`
}

func (op DeleteOperation) GetBytes() ([]byte, error) {
	jsonBytes, err := json.Marshal(op)
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, 0, len("DELETE ")+len(jsonBytes)+1)
	bytes = append(bytes, []byte("DELETE ")...)
	bytes = append(bytes, jsonBytes...)
	bytes = append(bytes, []byte("\n")...)
	return bytes, nil
}

func (op DeleteOperation) Recover(s *Storage) error {
	s.Delete(op.K)
	return nil
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

func (log *PersistentLogger) Append(op LogOperation) {
	bytes, err := op.GetBytes()
	if err != nil {
		panic(err)
	}
	if _, err = log.file.Write(bytes); err != nil {
		panic(err)
	}
}

func getOperation(line string) (LogOperation, error) {
	tag, data, found := strings.Cut(line, " ")
	if !found {
		return nil, fmt.Errorf("invalid log line format: need ' '")
	}
	var op LogOperation
	switch tag {
	case "PUT":
		putOp := PutOperation{}
		if err := json.Unmarshal([]byte(data), &putOp); err != nil {
			return nil, err
		}
		op = putOp
	case "DELETE":
		delOp := DeleteOperation{}
		if err := json.Unmarshal([]byte(data), &delOp); err != nil {
			return nil, err
		}
		op = delOp
	default:
		return nil, fmt.Errorf("unknown log line type: %s", tag)
	}
	return op, nil
}

func (log *PersistentLogger) Replay(ch chan LogOperation) error {
	log.file.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(log.file)
	defer close(ch)
	for scanner.Scan() {
		bytes := scanner.Text()
		op, err := getOperation(bytes)
		if err != nil {
			return err
		}
		ch <- op
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
