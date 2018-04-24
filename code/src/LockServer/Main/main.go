package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func main() {

	var lockServer lockServer
	lockServer.InitService()
	lockServer.StartService()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	lockServer.StopService()
	log.Info("task interface server quit")
}
