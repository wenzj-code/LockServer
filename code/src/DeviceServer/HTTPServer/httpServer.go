package HTTPServer

import (
	"DeviceServer/Handle"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

/*
HttpInit Http服务的初始化
*/
func HTTPInit(HTTPAddrPort string) error {
	log.Info("httpserver start:", HTTPAddrPort)
	http.HandleFunc("/dev-ctrl", httpServerFunc)
	err := http.ListenAndServe(HTTPAddrPort, nil)
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

	log.Debug("value:", req.Form)
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

	conn, isExist := Handle.ConnInfo[gwid[0]]
	if !isExist {
		log.Error("该网关不在线:", gwid)
		return
	}

	Handle.DevCtrl(conn, gwid[0], deviceid[0])
}
