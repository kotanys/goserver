package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func getUrlParameters(u *url.URL) map[string]string {
	urlParams, _ := url.ParseQuery(u.RawQuery)
	mapParams := make(map[string]string)
	for k, v := range urlParams {
		mapParams[k] = v[0]
	}
	return mapParams
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("Getting %q\n", req.URL)
	params := getUrlParameters(req.URL)
	fmt.Fprintln(w, params)
}

func putHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("Putting %q\n", req.URL)
	params := getUrlParameters(req.URL)
	fmt.Fprintln(w, params)
}

func startServer(server *http.Server) {
	fmt.Println("Server started on", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

const configFile string = "config.json"

func main() {
	config, err := ReadConfig(configFile)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/put", putHandler)
	server := &http.Server{Addr: fmt.Sprintf("localhost:%v", config.Port), Handler: mux}
	go startServer(server)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal

	fmt.Println("SIGINT recieved. Shutting down the server.")
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(context); err != nil {
		fmt.Println(err)
	}
}
