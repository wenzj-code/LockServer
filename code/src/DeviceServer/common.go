package main

import (
	Redis "RedisOpt"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

//RedisServerOpt 服务列表
var RedisServerOpt *Redis.RedisOpt

func InitCommon() error {
	InitConfig()
	config := GetConfig()

	RedisServerOpt = &Redis.RedisOpt{}
	err := RedisServerOpt.InitSingle(config.RedisAddr, config.RedisPwd, config.RedisServerNum)
	if err != nil {
		log.Error("err:", err)
	}
	return err
}

//GetLocalIP 获取本地IP
func GetLocalIP() (string, error) {
	ip, err := execSh("ifconfig | grep ^e -A2 | grep 'inet addr' | awk '{print $2}' | awk -F: '{print $2}'")
	return string(ip), err
}

func execSh(cmdStr string) ([]byte, error) {
	args := strings.Split(cmdStr, " ")
	cmd := exec.Command(args[0], args[1:]...)
	data, err := cmd.Output()
	if err != nil {
		return data, err
	}
	return data, err
}
