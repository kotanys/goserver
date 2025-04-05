package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const configFile string = "config.json"

func startServer(server *http.Server) {
	fmt.Println("starting server on port", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("server closed with status:", err)
	}
}

func main() {
	config, err := ReadConfig(configFile)
	if err != nil {
		panic(err)
	}

	storage := NewStorage()
	handler := NewStorageHttpHandler(storage)
	server := &http.Server{Addr: fmt.Sprintf(":%v", config.Port), Handler: handler}
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
