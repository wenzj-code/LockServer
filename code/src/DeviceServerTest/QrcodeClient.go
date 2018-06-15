package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"math/rand"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	QrcodeServerVersion uint32 = 106
)

type QrcodeClient struct {
	addr string
	conn net.Conn

	//客户端是否显示了二维码
	IsShowQrcode bool

	//二维码是否被扫描
	HasBeenScan bool

	//判断是否获取二维码
	HasQrcodeData bool

	QrcodeData []byte

	//status (0 未注册,1 注册成功,2注册失败[设备ID冲突])
	registerStatus int

	//心跳包计数
	connectTime int
}

var qrcodeClient QrcodeClient

func (opt *QrcodeClient) InitStatus() {
	opt.IsShowQrcode = false
	opt.HasBeenScan = false
}

//ConnectQrcodeServer .
func (opt *QrcodeClient) ConnectQrcodeServer(addr string) {
	opt.addr = addr
	var err error
	for {
		opt.conn, err = net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	log.Info("连接二维码服务成功")
	// 一连上qrserver就发送注册请求

	go opt.receiveData()
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				opt.connectTime++
				if opt.connectTime > 4 {
					opt.connectTime++
					opt.reconnect()
					log.Debug("超时")
				}
				opt.sendRegisterRequest(gConfigOpt.DeviceID)
			}
		}
	}()
}

//SetRegisterStatus .
func (opt *QrcodeClient) SetRegisterStatus(status int) {
	opt.registerStatus = status
}

//GetRegisterStatus .
func (opt *QrcodeClient) GetRegisterStatus() int {
	return opt.registerStatus
}

func (opt *QrcodeClient) reconnect() {
	// if !opt.isConnectQrcodeServer {
	// 	return
	// }
	log.Debug("重连二维码服务")
	if opt.conn != nil {
		opt.conn.Close()
	}
	opt.conn = nil
	var err error
	for {
		opt.conn, err = net.Dial("tcp", opt.addr)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	opt.sendRegisterRequest(gConfigOpt.DeviceID)
	go opt.receiveData()
}

func (opt *QrcodeClient) receiveData() {
	headLen := len("HTTP-JSON-BOCHIOT")
	head := make([]byte, headLen+5)
	for {

		if _, err := io.ReadFull(opt.conn, head[:headLen+1]); err != nil {
			log.Error("服务接收失败: " + err.Error())
			go opt.reconnect()
			return
		}
		opt.connectTime = 0

		//log.Debug("head:", head)

		if _, err := io.ReadFull(opt.conn, head[headLen+1:headLen+5]); err != nil {
			log.Error("服务接收失败: " + err.Error())
			go opt.reconnect()
			return
		}

		//log.Debug("head1:", head[headLen+1:headLen+5])

		len := binary.BigEndian.Uint32(head[headLen+1 : headLen+5])
		//log.Debug("len:", len)
		ibuf := make([]byte, len+5)

		if _, err := io.ReadFull(opt.conn, ibuf); err != nil {
			log.Error("服务接收失败: " + err.Error())
			go opt.reconnect()
			return
		}

		dataBuf := ibuf[:len]

		log.Debug("recv msg:", string(ibuf))

		dataMap := make(map[string]interface{})
		err := json.Unmarshal(dataBuf, &dataMap)
		if err != nil {
			log.Error("err:", err)
			continue
		}
		v, isExit := dataMap["cmd"]
		if !isExit {
			log.Error("cmd not exit:", string(dataBuf))
			continue
		}
		if v != "dev_ctrl" {
			continue
		}

		v, isExit = dataMap["device_info"]
		if !isExit {
			log.Error("device_info not exit:", string(dataBuf))
			continue
		}
		deviceMap := v.(map[string]interface{})
		v, isExit = deviceMap["device_mac"]
		if !isExit {
			log.Error("device_mac not exit:", string(dataBuf))
			continue
		}

		deviceID := v.(string)
		opt.sendCtrlRequest(gConfigOpt.DeviceID, deviceID)
		opt.sendBarrayRequest(gConfigOpt.DeviceID, deviceID)
	}
}

func (opt *QrcodeClient) sendRegisterRequest(deviceID string) error {
	//HTTP-JSON-BOCHIOT12345{"cmd": "gw_register","swm_gateway_info": {"gw_mac": "1aaa01000001","gw_ip_addr": "192.168.0.104"}}\n
	dataMap := make(map[string]interface{})
	gwMap := make(map[string]interface{})
	gwMap["gw_mac"] = deviceID
	gwMap["gw_ip_addr"] = "ipipip"
	dataMap["cmd"] = "gw_register"
	dataMap["swm_gateway_info"] = gwMap

	dataBuf, _ := json.Marshal(dataMap)

	return opt.sendDataTest(dataBuf)
}

func (opt *QrcodeClient) sendCtrlRequest(gwID, deviceID string) error {
	//HTTP-JSON-BOCHIOT12345{"cmd": "gw_register","swm_gateway_info": {"gw_mac": "1aaa01000001","gw_ip_addr": "192.168.0.104"}}\n
	dataMap := make(map[string]interface{})
	devMap := make(map[string]interface{})
	devMap["device_mac"] = deviceID
	devMap["switchStatus"] = "0"
	dataMap["cmd"] = "d2s_status"
	dataMap["gw_mac"] = gwID
	dataMap["device_info"] = devMap

	dataBuf, _ := json.Marshal(dataMap)
	return opt.sendDataTest(dataBuf)
}

func (opt *QrcodeClient) sendBarrayRequest(gwID, deviceID string) error {
	//HTTP-JSON-BOCHIOT12345{"cmd": "gw_register","swm_gateway_info": {"gw_mac": "1aaa01000001","gw_ip_addr": "192.168.0.104"}}\n
	dataMap := make(map[string]interface{})
	devMap := make(map[string]interface{})
	devMap["device_mac"] = deviceID
	devMap["battery"] = rand.Intn(100)
	dataMap["cmd"] = "d2s_battery"
	dataMap["gw_mac"] = gwID
	dataMap["device_info"] = devMap

	dataBuf, _ := json.Marshal(dataMap)
	return opt.sendDataTest(dataBuf)
}

func (opt *QrcodeClient) sendDataTest(msg []byte) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("二维码服务：", err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("HTTP-JSON-BOCHIOT12345")
	binary.Write(buf, binary.BigEndian, msg)
	buf.WriteString("\n")
	_, err := opt.conn.Write(buf.Bytes())
	if err != nil {
		log.Error("发送给二维码服务失败:", err)
		opt.reconnect()
		return err
	}
	return nil
}
