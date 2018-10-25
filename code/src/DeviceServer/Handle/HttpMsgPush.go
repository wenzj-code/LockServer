package Handle

import (
	"DeviceServer/Config"
	"DeviceServer/DBOpt"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

/*
	模块说明：用于信息推送到WechatAPI
*/

//推送开门，电量信息给WechatAPI
func pushMsgDevCtrl(deviceID, requestid string, barray float64, status int) {
	config := Config.GetConfig()
	httpServerIP := fmt.Sprintf("http://%s/report/door-ctrl-rsp?deviceid=%s&barry=%f&status=%d&requestid=%s",
		config.ReportHTTPAddr, deviceID, barray, status, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		return
	}
	defer resp.Body.Close()
	rspData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		return
	}
	log.Info("上报成功:", deviceID, ",msg:", string(rspData))
}

//推送消息发卡/密码的响应给WechatAPI, status = dna
func pushMsgSettingPassword(deviceID, keyVal, requestid string, keyType int, status int) {
	config := Config.GetConfig()
	httpServerIP := fmt.Sprintf("http://%s/report/dev-setting-password-status?deviceid=%s&keyvalue=%s&keytype=%d&status=%d&requestid=%s",
		config.ReportHTTPAddr, deviceID, keyVal, keyType, status, requestid)
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
	log.Info("卡/密码设置上报成功:", deviceID)
}

//取消卡/密码开门的响应给WechatAPI, status = dna
func pushMsgCancelPassword(deviceID, keyVal, requestid string, keyType int, status int) {
	config := Config.GetConfig()
	httpServerIP := fmt.Sprintf("http://%s/report/dev-cancel-password-status?deviceid=%s&keyvalue=%s&keytype=%d&status=%d&requestid=%s",
		config.ReportHTTPAddr, deviceID, keyVal, keyType, status, requestid)
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
	log.Info("卡/密码取消上报成功:", deviceID)
}

//取消卡/密码开门的响应给WechatAPI
func pushMsgCardOpenLockRsp(deviceID, keyVal, openTime, requestid string, keyType int) {
	config := Config.GetConfig()
	httpServerIP := fmt.Sprintf("http://%s/report/card-openlock-record?deviceid=%s&keyvalue=%s&keytype=%d&opentime=%s&requestid=%s",
		config.ReportHTTPAddr, deviceID, keyVal, keyType, openTime, requestid)
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
	log.Info("刷卡上报成功:", deviceID)
}

//@cmt 发送 *清除节点卡号密码信息*的结果给 WechatAPI
func pushMsgResetDev(deviceID, requestid string, resetStatus int) {
	config := Config.GetConfig()
	httpServerIP := fmt.Sprintf("http://%s/report/reset-dev?deviceid=%s&resetStatus=%d&requestid=%s",
		config.ReportHTTPAddr, deviceID, resetStatus, requestid)
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
	log.Info("清除节点卡号密码信息 上报成功:", deviceID)
}

//@cmt 发送 *设备常开常闭*的结果给 WechatAPI
func pushMsgDevNonc(deviceID, requestid string, status, setStatus int) {
	config := Config.GetConfig()
	httpServerIP := fmt.Sprintf("http://%s/report/dev-nonc-set?deviceid=%s&setStatus=%d&status=%d&requestid=%s",
		config.ReportHTTPAddr, deviceID, setStatus, status, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		return
	}
	defer resp.Body.Close()

	ackMsg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		return
	}
	log.Info("设备常开常闭状态 上报成功:", deviceID, ",ackmsg:", string(ackMsg))
}

//网关掉线通知
func PushMsgGatewayOffLineRsp(gatewayID string) error {
	err := DBOpt.GetDataOpt().SetGatwayOffline(gatewayID)
	if err != nil {
		log.Error("err:", err)
	}
	config := Config.GetConfig()
	httpServerIP := fmt.Sprintf("http://%s/report/gateway-offline-report?gatewayid=%s",
		config.ReportHTTPAddr, gatewayID)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		return err
	}
	log.Info("掉线上报成功:", gatewayID)
	return nil
}
