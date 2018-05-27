package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

const (
	MsgTypeScanInfo = 1
)

func HandlerMsg(MsgBody []byte, messageID string, ack func(string, string, error) error) (err error) {
	err = handlerMsgDeal(MsgBody)
	msgMD5 := fmt.Sprintf("%x", md5.Sum(MsgBody))
	return ack(messageID, msgMD5, err)
}

//HandleMqMsg .
func handlerMsgDeal(MsgBody []byte) error {
	//{"GatewayID:"asfafafaf","DeviceID:"13124124124"}
	if len(MsgBody) < 10 {
		return nil
	}
	dataMap := make(map[string]interface{})
	err := json.Unmarshal(MsgBody, &dataMap)
	if err != nil {
		log.Error("err: ", err)
		return nil
	}

	val, exist := dataMap["GatewayID"]
	if !exist {
		log.Error("GatewayID not exist")
		return nil
	}
	GatewayID := val.(string)

	val, exist = dataMap["DeviceID"]
	if !exist {
		log.Error("DeviceID not exist")
		return nil
	}
	DeviceID := val.(string)

	log.Debug("DeviceID:", DeviceID, ",GatewayID:", GatewayID)
	ReportMsgRMQ.PublishTopic([]byte(DeviceID))
	return nil
}
