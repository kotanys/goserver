package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

type StorageHTTPHandler struct {
	mux     *http.ServeMux
	storage *Storage
	cfg     *HTTPConfig
}

func getURLParameters(u *url.URL) map[string]string {
	urlParams, _ := url.ParseQuery(u.RawQuery)
	mapParams := make(map[string]string)
	for k, v := range urlParams {
		mapParams[k] = v[0]
	}
	return mapParams
}

func (h *StorageHTTPHandler) isMethodAllowed(method string) bool {
	return slices.ContainsFunc(h.cfg.Methods, func(s string) bool { return strings.EqualFold(s, "all") }) || slices.Contains(h.cfg.Methods, method)
}

func (h *StorageHTTPHandler) GetHandler(w http.ResponseWriter, req *http.Request) {
	if !h.isMethodAllowed(http.MethodGet) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
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

func (h *StorageHTTPHandler) PutHandler(w http.ResponseWriter, req *http.Request) {
	if !h.isMethodAllowed(http.MethodPut) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
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

func (h *StorageHTTPHandler) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	if !h.isMethodAllowed(http.MethodDelete) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	params := getURLParameters(req.URL)
	key, exist := params["key"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.storage.Delete(StorageKey(key))
}

func NewStorageHTTPHandler(storage *Storage, cfg *HTTPConfig) *StorageHTTPHandler {
	h := &StorageHTTPHandler{}
	h.storage = storage
	h.cfg = cfg
	h.mux = &http.ServeMux{}
	h.mux.HandleFunc("/get", h.GetHandler)
	h.mux.HandleFunc("/put", h.PutHandler)
	h.mux.HandleFunc("/delete", h.DeleteHandler)
	return h
}

func (h *StorageHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.mux.ServeHTTP(w, req)
}
