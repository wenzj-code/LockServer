package main

import (
	"fmt"
	"log/syslog"
	"os"
	"os/signal"

	RMQ "RMQ"
	"gotcp"
	"syscall"
	"vislog"

	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/syslog"
)

var RecvMsgRMQ RMQ.RMQOpt
var ReportMsgRMQ RMQ.RMQOpt

var gConfigOpt Option
var Srv *gotcp.Server
var ConnInfo map[string]*gotcp.Conn = make(map[string]*gotcp.Conn)

//key:deviceID, val:bool(是否收到回应)
var QrcodeACKRecv map[string]bool

func init() {
	QrcodeACKRecv = make(map[string]bool)
}

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
	log.Info("Qrcode server quit")
}

func start() {
	gConfigOpt = loadConfig()
	initLog(gConfigOpt.LogFile, gConfigOpt.LogLevel, gConfigOpt.SysLogAddr)

	log.Info("Qrcode server is starting.....version:", version)
	RecvMsgRMQ.InitMQTopic(gConfigOpt.RecvAmqpURI, gConfigOpt.RecvExchangeName, gConfigOpt.RecvChanReadQName,
		"", gConfigOpt.RecvRoutKey, HandlerMsg)

	ReportMsgRMQ.InitMQTopic(gConfigOpt.ReportAmqpURI, gConfigOpt.ReportExchnageName, "", "", gConfigOpt.ReportRoutKey, nil)

	Srv = gotcp.NewServer(&CallBack{})
	go Srv.StartServer(gConfigOpt.Addr, "ControlServer")
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
