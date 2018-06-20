package Handle

import (
	"DeviceServer/Common"
	"DeviceServer/Config"
	"DeviceServer/DBOpt"
	"fmt"
	"gotcp"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

//网关注册信息
func gatewayRegister(conn *gotcp.Conn, cmd string, dataMap map[string]interface{}) {
	val, isExist := dataMap["swm_gateway_info"]
	if !isExist {
		log.Error("swm_gateway_info 字段不存在:", dataMap)
		return
	}
	gwInfo := val.(map[string]interface{})
	val, isExist = gwInfo["gw_mac"]
	if !isExist {
		log.Error("gw_mac 字段不存在:", dataMap)
		return
	}
	gatewayID := val.(string)

	ConnInfo[gatewayID] = conn

	//网关注册的时候，保存网关所注册的服务器地址到Redis
	err := Common.RedisServerOpt.Set(gatewayID, Config.GetConfig().HTTPServer, Config.GetConfig().RedisTimeOut)
	if err != nil {
		log.Error("err:", err)
		return
	}

	dataMap = make(map[string]interface{})
	dataMap["cmd"] = cmd
	dataMap["systemTime"] = time.Now().Format("2006-01-02 15:04:05")
	dataMap["statuscode"] = 0
	ackGateway(conn, dataMap)
}

//开门状态返回
func doorCtrlDeal(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	val, isExist := data["device_info"]
	if !isExist {
		log.Error("device_info 字段不存在:", data)
		return
	}
	deviceInfo := val.(map[string]interface{})
	val, isExist = deviceInfo["device_mac"]
	if !isExist {
		log.Error("device_mac 字段不存在:", data)
		return
	}
	deviceID := val.(string)

	pushMsg(deviceID, -1, 1)
}

//电量信息上报
func doorReportBarry(conn *gotcp.Conn, cmd string, data map[string]interface{}) {
	val, isExist := data["device_info"]
	if !isExist {
		log.Error("device_info 字段不存在:", data)
		return
	}
	deviceInfo := val.(map[string]interface{})
	val, isExist = deviceInfo["device_mac"]
	if !isExist {
		log.Error("device_mac 字段不存在:", data)
		return
	}
	deviceID := val.(string)

	val, isExist = deviceInfo["battery"]
	if !isExist {
		log.Error("battery 字段不存在:", data)
		return
	}
	battery := val.(float64)

	pushMsg(deviceID, battery, 1)
}

//获取设备列表
func requestDeviceList(conn *gotcp.Conn, cmd string, data map[string]interface{}) {

	val, isExist := data["swm_gateway_info"]
	if !isExist {
		log.Error("swm_gateway_info 字段不存在:", data)
		return
	}

	gatewayInfo := val.(map[string]interface{})
	val, isExist = gatewayInfo["gw_mac"]
	if !isExist {
		log.Error("gw_mac 字段不存在:", data)
		return
	}
	gatewayID := val.(string)
	//通过网关ID查询数据库,获取网关下的所有设备
	deviceList, err := DBOpt.GetDataOpt().GetDeviceIDList(gatewayID)
	if err != nil {
		log.Error("err:", err)
		return
	}
	log.Debug("deviceList:", deviceList)

	gwMap := make(map[string]interface{})
	deviceInfoArray := make([]Common.DeviceInfo, 0)
	gwMap["gw_mac"] = gatewayID

	count := 0
	//设备列表过大，分包处理
	lenMap := len(deviceList)
	countDeviceList := 0
	for k := range deviceList {
		countDeviceList++
		deviceInfo := new(Common.DeviceInfo)
		deviceInfo.DeviceID = k
		deviceInfo.RegStatus = 1
		deviceInfoArray = append(deviceInfoArray, *deviceInfo)
		//50个设备分包，或者最后一包
		if count == 50 || countDeviceList == lenMap {
			dataMap := make(map[string]interface{})
			dataMap["cmd"] = "d2s_request_devices"
			dataMap["swm_gateway_info"] = gwMap
			dataMap["device_info"] = deviceInfoArray
			dataMap["statuscode"] = 0
			ackGateway(conn, dataMap)

			count = 0
			deviceInfoArray = make([]Common.DeviceInfo, 0)
		}
	}
}

//DevCtrl 控制开门
func DevCtrl(conn *gotcp.Conn, gatewayID, deviceID string) {
	dataMap := make(map[string]interface{})
	deviceInfo := make(map[string]interface{})
	deviceInfo["device_mac"] = deviceID
	deviceInfo["switchStatus"] = 1
	dataMap["cmd"] = "dev_ctrl"
	dataMap["device_info"] = deviceInfo
	dataMap["statuscode"] = 0

	ackGateway(conn, dataMap)
}

//推送消息给WechatAPI
func pushMsg(deviceID string, barray float64, status int) {
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
