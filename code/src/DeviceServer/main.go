package main

import (
	"DeviceServer/Common"
	"DeviceServer/Config"
	"fmt"
	"log/syslog"
	"os"
	"os/signal"

	"gotcp"
	"syscall"
	"vislog"

	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
)

var Srv *gotcp.Server

//gatewayID,conn
var ConnInfo map[string]*gotcp.Conn = make(map[string]*gotcp.Conn)

var (
	version         = "1.1.3.1"
	versionTime     = "20171208"
	versionFunction = ""
)

func usage() bool {
	args := os.Args
	if len(args) == 2 && (args[1] == "--version" || args[1] == "-v" ||
		args[1] == "version") {
		fmt.Println("version:", version)
		fmt.Println("build time:", versionTime)
		fmt.Println("function:", versionFunction)
		return true
	} else if len(args) == 2 && (args[1] == "--help" || args[1] == "-h" ||
		args[1] == "help") {
		fmt.Println("1 --version -v version ")
		fmt.Println("2 kill -USR1 pid  open debug info")
		fmt.Println("3 kill -USR2 pid  close debug info")
		return true
	}

	return false
}

func main() {
	if usage() {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			log.Error("main:", err)
		}
	}()

	start()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	stop()
	log.Info("DeviceServer server quit")
}

func start() {
	Config.InitConfig()
	config := Config.GetConfig()
	initLog(config.LogFile, config.LogLevel, config.SysLogAddr)
	err := Common.InitCommon()
	if err != nil {
		log.Error("err:", err)
		return
	}

	log.Info("DeviceServer server is starting.....version:", version, ",port:", config.Addr)
	Srv = gotcp.NewServer(&CallBack{})
	go Srv.StartServer(config.Addr, "ControlServer")

	go httpInit(config.HTTPServerPORT)

}

func stop() {
}

func initLog(logfile string, loglevel string, syslogAddr string) {
	if logfile == "" {
		logfile = "DevStatusServer.log"
	}
	hook, err := vislog.NewVislogHook(logfile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.AddHook(hook)

	syshook, err := logrus_syslog.NewSyslogHook("udp", syslogAddr, syslog.LOG_DEBUG, os.Args[0])
	if err == nil {
		log.AddHook(syshook)
	}

	level, err := log.ParseLevel(loglevel)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&log.JSONFormatter{})

}
