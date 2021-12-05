package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type scwsServer interface {
	Shutdown(ctx context.Context) error
}

type scwsSettings interface {
	Reload()
}

// it will not work within docker, because docker entrypoint is bash
// TODO: Adjust Dockerfile in order to send signal to scws
func catchSignal(srv scwsServer, setts scwsSettings) {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		for {
			s := <-signalChanel
			switch s {
			case syscall.SIGHUP:
				// Reload doesn't load new env variables or new values,
				// because child process has a copy of env vars from parent process
				// reload will work for config Maps
				log.Println("Got signal to reload settings.")
				setts.Reload()
			case syscall.SIGINT:
				log.Println("Signal interrupt triggered.")
				srv.Shutdown(context.TODO())
			case syscall.SIGTERM:
				log.Println("Signal terminte triggered.")
				srv.Shutdown(context.TODO())
			case syscall.SIGQUIT:
				log.Println("Signal quit triggered.")
				srv.Shutdown(context.TODO())
			}
		}
	}()
}
