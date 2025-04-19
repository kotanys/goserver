package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func startServer(server *http.Server) {
	fmt.Println("starting server on port", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("server closed with status:", err)
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("usage: %v config_file\n", os.Args[0])
		fmt.Println("provide a configuration file")
		return
	}
	config, err := ReadConfig(os.Args[len(os.Args)-1])
	if err != nil {
		panic(err)
	}

	var logger *PersistentLogger = nil
	if config.LogFile != "" {
		logger, err = NewPeristentLogger(config.LogFile)
		if err != nil {
			fmt.Println("! error creating the logger:", err.Error())
		}
		defer logger.Close()
	}
	storage := NewStorage(logger, logger != nil)
	cfgHttp := MakeHTTPConfig(config)
	handler := NewStorageHTTPHandler(storage, cfgHttp)
	server := &http.Server{Addr: fmt.Sprintf(":%v", config.Port), Handler: handler}
	go startServer(server)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal

	fmt.Println("SIGINT recieved. Shutting down the server.")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}
