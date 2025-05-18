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
	servers := []*http.Server{}
	storage := NewStorage(logger, config.Persistent)
	cfgHttp := MakeHTTPConfig(config, false)
	handler := NewStorageHTTPHandler(storage, cfgHttp)
	server := &http.Server{Addr: fmt.Sprintf(":%v", cfgHttp.Port), Handler: handler}
	servers = append(servers, server)
	go startServer(server)

	if config.InternalPort != 0 {
		internalCfgHttp := MakeHTTPConfig(config, true)
		internalHandler := NewInternalStorageHTTPHandler(storage, internalCfgHttp)
		internalServer := &http.Server{Addr: fmt.Sprintf(":%v", internalCfgHttp.Port), Handler: internalHandler}
		servers = append(servers, internalServer)
		go startServer(internalServer)
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal

	fmt.Println("SIGINT recieved. Shutting down the server.")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, server := range servers {
		if err := server.Shutdown(ctx); err != nil {
			fmt.Println(err)
		}
	}
}
