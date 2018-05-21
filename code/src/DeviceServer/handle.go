package main

import (
	"encoding/json"
	"gotcp"

	log "github.com/Sirupsen/logrus"
)

type CallBack struct {
}

func (cb *CallBack) Close() {
}

func (cb *CallBack) HandleMsg(conn *gotcp.Conn, MsgBody []byte) error {
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
}

func (cb *CallBack) doorCtrlDeal(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	cb.push2RMQ(data)
}

func (cb *CallBack) reportInfoDeal(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	cb.ackMsg(conn, cmd, data)
	cb.push2RMQ(data)
}

func (cb *CallBack) push2RMQ(data map[string]interface{}) {
	dataBuf, _ := json.Marshal(data)
	err := ReportMsgRMQ.PublishTopic(dataBuf)
	if err != nil {
		log.Error("err:", err)
	}
}

func (cb *CallBack) ackMsg(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	data = make(map[string]interface{})
	data["Cmd"] = cmd
	data["code"] = 0

	dataBuf, _ := json.Marshal(data)
	cb.baseSendMsg(conn, dataBuf)
}

func (cb *CallBack) baseSendMsg(conn *gotcp.Conn, msg []byte) {

	conn.SendChan <- msg
}
