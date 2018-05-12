package main

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

const (
	MsgTypeScanInfo = 1
)

//HandleMqMsg .
func HandleMqMsg(MsgInfo []byte) error {
	//{"ScanID:"asfafafaf","DeviceID:"13124124124","Status:"SUCCESS/FAILED/Progress","Percent":50}
	defer func() {
		if err := recover(); err != nil {
			log.Error("HandleMqMsg:", err) // 这里的err其实就是panic传入的内容，55
		}
	}()
	if len(MsgInfo) == 0 {
		log.Error("err rmq msg len is 0")
		return nil
	}
	var msg map[string]interface{}

	err := json.Unmarshal(MsgInfo, &msg)
	if err != nil {
		log.Error("fail to unmarshal msg from mq: ", err, ",msg:", msg)
		return nil
	}

	return handleScanInfoMsg(msg)
}

func handleScanInfoMsg(msg map[string]interface{}) error {
	//{"ScanID:"asfafafaf","DeviceID:"13124124124","Action":"Normal/BodyEval","Status:"SUCCESS/FAILED/Progress","Percent":50}

	defer func() {
		if err := recover(); err != nil {
			log.Error("异常:", err)
		}
	}()
	log.Debug("handleScanInfoMsg:", msg)

	DeviceID, ok := msg["DeviceID"].(string)
	if !ok {
		log.Error("fail to get DeviceID")
		return nil
	}
	conn, ok := ConnInfo[DeviceID]
	if !ok {
		log.Error("fail to find device connection")
		return nil
	}
	log.Debug("DeviceID:", DeviceID, ",coon:", conn.GetRemoteAddr())

	ScanID, ok := msg["ScanID"].(string)
	if !ok {
		log.Error("fail to find scan id")
		return nil
	}

	Status, ok := msg["Status"].(string)
	if !ok {
		log.Error("fail to find Status")
		return nil
	}

	Action, ok := msg["Action"].(string)
	if !ok {
		log.Error("fail to find Action")
		return nil
	}

	actionType := 0
	if Action == "Normal" {
		actionType = 1
	} else if Action == "BodyEval" {
		actionType = 2
	} else {
		log.Error("Error Action value:", Action)
		return nil
	}

	if Status == "SUCCESS" {
		handleNotifyStatusResponse(conn, true, ScanID, actionType)
	} else if Status == "FAILED" {
		handleNotifyStatusResponse(conn, false, ScanID, actionType)
	} else {
		Percent, ok := msg["Percent"].(float64)
		if !ok {
			log.Error("fail to find Percent")
			return nil
		}
		handleReportPercentResponse(conn, uint32(Percent), ScanID)
	}
	return nil
}
