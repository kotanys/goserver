package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Loggable[T any] interface {
	FromString(string) T
	ToString(T) string
}

type Logger[T Loggable[T]] interface {
	Log(T)
	Retrieve() []T
}

type StorageHttpHandler struct {
	mux     *http.ServeMux
	storage *Storage
}

func getUrlParameters(u *url.URL) map[string]string {
	urlParams, _ := url.ParseQuery(u.RawQuery)
	mapParams := make(map[string]string)
	for k, v := range urlParams {
		mapParams[k] = v[0]
	}
	return mapParams
}

func (handler *StorageHttpHandler) GetHandler(w http.ResponseWriter, req *http.Request) {
	params := getUrlParameters(req.URL)
	key, exist := params["key"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, exist := handler.storage.Get(StorageKey(key))
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%v\n", key)
		return
	}
	w.Write(value)
}

func (handler *StorageHttpHandler) PutHandler(w http.ResponseWriter, req *http.Request) {
	params := getUrlParameters(req.URL)
	key, exist := params["key"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v\n", err.Error())
		return
	}
	handler.storage.Put(StorageKey(key), body)
}

func NewStorageHttpHandler(storage *Storage) *StorageHttpHandler {
	handler := &StorageHttpHandler{}
	handler.storage = storage
	handler.mux = &http.ServeMux{}
	handler.mux.HandleFunc("/get", handler.GetHandler)
	handler.mux.HandleFunc("/put", handler.PutHandler)
	return handler
}

func (handler *StorageHttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler.mux.ServeHTTP(w, req)
}
