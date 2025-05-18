package main

import (
	"bytes"
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
	client  *http.Client
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

func (h *StorageHTTPHandler) resendPutToSlaves(key StorageKey, value StorageValue) error {
	for _, slave := range h.cfg.Slaves {
		req, _ := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("http://localhost:%v/put_internal?key=%v", slave, key),
			bytes.NewBuffer(value),
		)
		resp, err := h.client.Do(req)
		if err != nil {
			fmt.Printf("Error resending PUT to %v\n", slave)
			return err
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Slave %v returned %v on resending PUT. Not good, desyncs are possible!\n", slave, resp.StatusCode)
			return err
		}
	}
	return nil
}

func (h *StorageHTTPHandler) resendDeleteToSlaves(key StorageKey) error {
	for _, slave := range h.cfg.Slaves {
		req, _ := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("http://localhost:%v/delete_internal?key=%v", slave, key),
			nil,
		)
		resp, err := h.client.Do(req)
		if err != nil {
			fmt.Printf("Error resending DELETE to %v\n", slave)
			return err
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Slave %v returned %v on resending DELETE. Not good, desyncs are possible!\n", slave, resp.StatusCode)
			return err
		}
	}
	return nil
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
	if err := h.resendPutToSlaves(StorageKey(key), body); err != nil {
		return
	}
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
	if err := h.resendDeleteToSlaves(StorageKey(key)); err != nil {
		return
	}
}

func NewStorageHTTPHandler(storage *Storage, cfg *HTTPConfig) *StorageHTTPHandler {
	h := &StorageHTTPHandler{}
	h.storage = storage
	h.cfg = cfg
	h.client = &http.Client{}
	h.mux = &http.ServeMux{}
	h.mux.HandleFunc("/get", h.GetHandler)
	h.mux.HandleFunc("/put", h.PutHandler)
	h.mux.HandleFunc("/delete", h.DeleteHandler)
	return h
}

func (h *StorageHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.mux.ServeHTTP(w, req)
}
