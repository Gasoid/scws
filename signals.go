package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type scwsServer interface {
	Shutdown(ctx context.Context) error
}

type scwsConfig interface {
	ParseEnv() error
}

// it will not work within docker, because docker entrypoint is bash
// TODO: Adjust Dockerfile in order to send signal to scws
func catchSignal(server scwsServer, config scwsConfig) {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		for {
			s := <-signalChanel
			switch s {
			case syscall.SIGHUP:
				fmt.Println("Signal hang up triggered.")
				config.ParseEnv()
			case syscall.SIGINT:
				fmt.Println("Signal interrupt triggered.")
				server.Shutdown(context.TODO())
			case syscall.SIGTERM:
				fmt.Println("Signal terminte triggered.")
				server.Shutdown(context.TODO())
			case syscall.SIGQUIT:
				fmt.Println("Signal quit triggered.")
				server.Shutdown(context.TODO())
			}
		}
	}()
}
