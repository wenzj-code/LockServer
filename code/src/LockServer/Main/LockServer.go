package main

import (
	"LockServer/Config"
	"LockServer/TcpServer"
	"gotcp"
	"vislog"

	log "github.com/Sirupsen/logrus"
)

type lockServer struct {
	config    *Config.Option //配置文件
	deviceSrv *gotcp.Server
}

func newLockServer() *lockServer {
	l := new(lockServer)
	return l
}

func (opt *lockServer) initLog(logfile string, loglevel string, syslogAddr string) {
	if logfile == "" {
		logfile = "LockServer.log"
	}

	hook, err := vislog.NewVislogHook(logfile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.AddHook(hook)

	level, err := log.ParseLevel(loglevel)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&log.JSONFormatter{})
}

/*
InitService ...: 初始化所有服务
*/
func (opt *lockServer) InitService() error {
	Config.InitConfig()
	opt.config = Config.GetConfig()

	opt.initLog(opt.config.LogFile, opt.config.LogLevel, opt.config.SysLogAddr)
	log.Info("InitService success:", opt.config)
	return nil
}

/*
StartService ...: 开始服务
*/
func (opt *lockServer) StartService() {
	config := opt.config

	opt.deviceSrv = gotcp.NewServer(&TcpServer.DeviceServer{})
	go opt.deviceSrv.StartServer(config.TCPAddr, "DeviceServer")
	log.Info("StartService:", config.TCPAddr)
}

/*
StopService ...
*/
func (opt *lockServer) StopService() {
	opt.deviceSrv.StopServer()
}
