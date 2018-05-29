package main

import (
	"DeviceServer/Common"
	"DeviceServer/Config"
	"encoding/json"
	"fmt"
	"gotcp"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type CallBack struct {
}

func (cb *CallBack) Close() {
}

func (cb *CallBack) HandleMsg(conn *gotcp.Conn, MsgBody []byte) error {
	log.Debug("msg:", string(MsgBody))
	data := make(map[string]interface{})
	err := json.Unmarshal(MsgBody, &data)
	if err != nil {
		log.Error("err: ", err)
		return nil
	}
	val, exist := data["Cmd"]
	if !exist {
		log.Error("cmd not exist:", string(MsgBody))
		return nil
	}

	log.Info("data:", data)
	cmd := val.(string)
	switch cmd {
	case "HB":
		cb.heartbeatDeal(conn, cmd, data)
	case "CTRL":
		cb.doorCtrlDeal(conn, cmd, data)
	case "REPORT":
		cb.reportInfoDeal(conn, cmd, data)
	default:
		log.Error("cmd invalid:", cmd)
	}

	return nil
}

func (cb *CallBack) heartbeatDeal(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	cb.ackMsg(conn, cmd, data)
	val, isExist := data["GatewayID"]
	if !isExist {
		log.Error("GatewayID 字段不存在:", data)
		return
	}
	DeviceID := val.(string)
	ConnInfo[DeviceID] = conn

	IPaddr := fmt.Sprintf("%s:%d", Common.GetLocalIP(), Config.GetConfig().HTTPServerPORT)

	err := Common.RedisServerOpt.Set(DeviceID, IPaddr, Config.GetConfig().RedisTimeOut)
	if err != nil {
		log.Error("err:", err)
		return
	}
}

func (cb *CallBack) doorCtrlDeal(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	cb.ackMsg(conn, cmd, data)

	DeviceID, isExist := data["DeviceID"]
	if !isExist {
		log.Error("DeviceID 字段不存在:", data)
		return
	}
	Barry, isExist := data["Barry"]
	if !isExist {
		log.Error("Barry 字段不存在:", data)
		return
	}
	Status, isExist := data["Status"]
	if !isExist {
		log.Error("Status 字段不存在:", data)
		return
	}

	cb.pushMsg(DeviceID.(string), Barry.(float64), int(Status.(float64)))
}

func (cb *CallBack) reportInfoDeal(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	cb.ackMsg(conn, cmd, data)

	DeviceID, isExist := data["DeviceID"]
	if !isExist {
		log.Error("DeviceID 字段不存在:", data)
		return
	}

	// GatewayID, isExist := data["GatewayID"]
	// if !isExist {
	// 	log.Error("GatewayID 字段不存在:", data)
	// 	return
	// }
	Barry, isExist := data["Barry"]
	if !isExist {
		log.Error("Barry 字段不存在:", data)
		return
	}

	cb.pushMsg(DeviceID.(string), Barry.(float64), 0)
}

func (cb *CallBack) pushMsg(deviceID string, barray float64, status int) {
	config := Config.GetConfig()
	httpServerIP := fmt.Sprintf("http://%s/report/dev-status?deviceid=%s&barry=%f&status=%d", config.ReportHTTPAddr, deviceID, barray, status)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		return
	}
	log.Info("上报成功:", deviceID)
}

func (cb *CallBack) ackMsg(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	data = make(map[string]interface{})
	data["Cmd"] = cmd
	data["code"] = 0

	dataBuf, _ := json.Marshal(data)
	BaseSendMsg(conn, dataBuf)
}

func BaseSendMsg(conn *gotcp.Conn, msg []byte) {
	conn.SendChan <- msg
}
