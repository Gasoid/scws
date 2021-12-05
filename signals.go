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
				fmt.Println("Signal hang up triggered.")
				setts.Reload()
			case syscall.SIGINT:
				fmt.Println("Signal interrupt triggered.")
				srv.Shutdown(context.TODO())
			case syscall.SIGTERM:
				fmt.Println("Signal terminte triggered.")
				srv.Shutdown(context.TODO())
			case syscall.SIGQUIT:
				fmt.Println("Signal quit triggered.")
				srv.Shutdown(context.TODO())
			}
		}
	}()
}
