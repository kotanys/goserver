package main

import (
	"fmt"
	"io"
	"net/http"
)

type InternalStorageHTTPHandler struct {
	mux     *http.ServeMux
	storage *Storage
	cfg     *HTTPConfig
}

func (h *InternalStorageHTTPHandler) GetHandler(w http.ResponseWriter, req *http.Request) {
	params := getURLParameters(req.URL)
	key, exist := params["key"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, exist := h.storage.Get(StorageKey(key))
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%v\n", key)
		return
	}
	w.Write(value)
}

func (h *InternalStorageHTTPHandler) PutHandler(w http.ResponseWriter, req *http.Request) {
	params := getURLParameters(req.URL)
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
	h.storage.Put(StorageKey(key), body)
}

func (h *InternalStorageHTTPHandler) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	params := getURLParameters(req.URL)
	key, exist := params["key"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.storage.Delete(StorageKey(key))
}

func NewInternalStorageHTTPHandler(storage *Storage, cfg *HTTPConfig) *InternalStorageHTTPHandler {
	h := &InternalStorageHTTPHandler{}
	h.storage = storage
	h.cfg = cfg
	h.mux = &http.ServeMux{}
	h.mux.HandleFunc("/get_internal", h.GetHandler)
	h.mux.HandleFunc("/put_internal", h.PutHandler)
	h.mux.HandleFunc("/delete_internal", h.DeleteHandler)
	return h
}

func (h *InternalStorageHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.cfg.Update()
	h.mux.ServeHTTP(w, req)
}
