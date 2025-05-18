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

func watchConfigChange(cfg *Config, filePath string, ctx context.Context) {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Can't stat %v for watching\n", filePath)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			stat, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("Can't stat %v for watching\n", filePath)
				return
			}

			if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
				cfg.Update(filePath)
				initialStat = stat
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("usage: %v config_file\n", os.Args[0])
		fmt.Println("provide a configuration file")
		return
	}
	fileName := os.Args[len(os.Args)-1]
	config, err := ReadConfig(fileName)
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

	cfgStorage := MakeStorageConfig(config)
	storage := NewStorage(logger, cfgStorage)

	servers := []*http.Server{}
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

	watchCtx, watchCancel := context.WithCancel(context.Background())
	go watchConfigChange(config, fileName, watchCtx)

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<-stopSignal

	watchCancel()
	fmt.Println("SIGINT recieved. Shutting down the server.")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	for _, server := range servers {
		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println(err)
		}
	}
}
