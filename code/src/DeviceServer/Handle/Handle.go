package Handle

import (
	"DeviceServer/Common"
	"bytes"
	"encoding/json"
	"errors"
	"gotcp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

//gatewayID,conn
var ConnInfo map[string]*gotcp.Conn = make(map[string]*gotcp.Conn)

type CallBack struct {
}

func (cb *CallBack) Close() {
}

func (cb *CallBack) HandleMsg(conn *gotcp.Conn, MsgBody []byte) error {
	if len(MsgBody) < 10 {
		log.Debug("msg:", string(MsgBody))

		return nil
	}
	log.Debug("msg:", string(MsgBody))
	if len(MsgBody) < len(Common.DefaultHead)+10 {
		log.Debug("wrong pack")
		return errors.New("wrong pack")
	}
	if !strings.Contains(string(MsgBody), Common.DefaultHead) {
		log.Debug("head err")
		return errors.New("head err")
	}
	jsonData := MsgBody[len(Common.DefaultHead)+5:]
	data := make(map[string]interface{})
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Error("err: ", err)
		return nil
	}
	val, exist := data["cmd"]
	if !exist {
		log.Error("cmd not exist:", string(MsgBody))
		return nil
	}

	log.Info("data:", data)
	cmd := val.(string)
	switch cmd {
	case "gw_register": //网关注册
		gatewayRegister(conn, cmd, data)
	case "d2s_status": //开门返回来的状态
		doorCtrlDeal(conn, cmd, data)
	case "d2s_request_devices": //网关请求所有节点信息
		requestDeviceList(conn, cmd, data)
	case "d2s_battery": //上报电量
		doorReportBarry(conn, cmd, data)
	default:
		log.Error("cmd invalid:", cmd)
	}

	return nil
}

func ackGateway(conn *gotcp.Conn, dataMap map[string]interface{}) {
	dataBuf, err := json.Marshal(dataMap)
	if err != nil {
		log.Error("err:", err)
		return
	}

	protoclBuf := getPackage(dataBuf)
	baseSendMsg(conn, protoclBuf)
}

func baseSendMsg(conn *gotcp.Conn, msg []byte) {
	conn.SendChan <- getPackage(msg)
}

/*
获取网关通信协议的组装包格式
*/
func getPackage(msg []byte) []byte {
	var crc int
	len := len(msg)
	var dataBuff bytes.Buffer
	//包头
	dataBuff.WriteString(Common.DefaultHead)
	//状态
	dataBuff.WriteByte(0x23)
	//包体长度
	dataBuff.WriteByte(byte(len >> 24))
	dataBuff.WriteByte(byte(len >> 16))
	dataBuff.WriteByte(byte(len >> 8))
	dataBuff.WriteByte(byte(len))
	//包内容
	dataBuff.Write(msg)
	//crc，目前不需要用到
	dataBuff.WriteByte(byte(crc >> 24))
	dataBuff.WriteByte(byte(crc >> 16))
	dataBuff.WriteByte(byte(crc >> 8))
	dataBuff.WriteByte(byte(crc))

	return dataBuff.Bytes()
}
