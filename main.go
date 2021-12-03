package main

import (
	"os"
	"os/signal"
	"syscall"
)

func catchSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}

func main() {
	catchSignal()
	Run()
}
