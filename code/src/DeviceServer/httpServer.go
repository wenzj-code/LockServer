package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

func httpInit(HTTPAddr string) error {
	log.Info("httpserver start:", HTTPAddr)

	http.HandleFunc("/dev-ctrl", httpServerFunc)
	err := http.ListenAndServe(HTTPAddr, nil)
	if err != nil {
		log.Error("err:", err)
		os.Exit(0)
	}
	return err
}

func httpServerFunc(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Error("err:", err)
		return
	}

	//log.Debug("value:", req.Form)
	io.WriteString(w, "recv ok")

	gwid, isExist := req.Form["gwid"]
	if !isExist {
		log.Error("gwid 字段不存在:", req.Form)
		return
	}

	deviceid, isExist := req.Form["deviceid"]
	if !isExist {
		log.Error("deviceid 字段不存在:", req.Form)
		return
	}

	conn, isExist := ConnInfo[gwid[0]]
	if !isExist {
		log.Error("该网关不在线:", gwid)
		return
	}
	data := make(map[string]interface{})
	data["DeviceID"] = deviceid[0]
	data["GatewayID"] = gwid[0]

	dataBuf, _ := json.Marshal(data)
	log.Debug("dataBuf:", string(dataBuf))
	BaseSendMsg(conn, dataBuf)
}
