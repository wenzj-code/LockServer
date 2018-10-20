package controllers

/*
	该模块主要用来接收微信端的请求设置
*/

import (
	"WechatAPI/DBOpt"
	"WechatAPI/common"
	"WechatAPI/config"
	"io/ioutil"
	"net/http"

	"crypto/sha256"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

//WechatController .
type WechatController struct {
	beego.Controller
}

//GetToken 通过APPID+Secrete生成token
func (c *WechatController) GetToken() {
	appid := c.GetString("appid")
	secret := c.GetString("secret")
	log.Info("appid=", appid, ",secret=", secret)
	if appid == "" || secret == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	//判断appid,secret是否存在,权限判断
	status, err := DBOpt.GetDataOpt().CheckAppIDSecret(appid, secret)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		log.Debug("appid,secrete不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}
	//组织应答报文,生成token的方法
	data := make(map[string]interface{})
	hs := sha256.New()
	io.WriteString(hs, secret+time.Now().String())
	token := fmt.Sprintf("%x", string(hs.Sum(nil)))

	//将token保存到Redis缓存
	common.RedisTokenOpt.Set(token, 1, config.GetConfig().RedisTokenTimeOut)

	data["code"] = 0
	data["token"] = token
	data["expired_in"] = config.GetConfig().RedisTokenTimeOut
	c.Data["json"] = data
	c.ServeJSON()
}

//DoorCtrlOpen 开门
func (c *WechatController) DoorCtrlOpen() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	method, err := c.GetInt("method")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	token := c.GetString("token")
	requestid := c.GetString("requestid")
	log.Info("DoorCtrlOpen DeviceID=", roomnu, ",Token=", token, ",appid:", appid)
	if roomnu == "" || appid == "" || token == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	//门禁开门
	if method == 3 {

	}

	serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, method)
	if !status {
		return
	}

	//向设备服务请求开门
	httpServerIP := fmt.Sprintf("http://%s/dev-ctrl?gwid=%s&deviceid=%s&requestid=%s", serverIP, gatewayID, DeviceID, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

	return
}

//SettingCardPassword 发卡、密码
func (c *WechatController) SettingCardPassword() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	keyvalue := c.GetString("keyvalue")
	keytype, err := c.GetInt("keytype")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	expireDate, err := c.GetInt64("expire-date")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	token := c.GetString("token")
	requestid := c.GetString("requestid")
	log.Info("SettingCardPassword DeviceID=", roomnu, ",Token=", token, ",appid:", appid)
	if roomnu == "" || appid == "" || token == "" ||
		keyvalue == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, 0)
	if !status {
		return
	}

	//向设备服务请求发卡
	httpServerIP := fmt.Sprintf("http://%s/setting-card-password?gwid=%s&deviceid=%s&keyvalue=%s&keytype=%d&expire-date=%d&requestid=%s",
		serverIP, gatewayID, DeviceID, keyvalue, keytype, expireDate, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

}

//CancleCardPassword 发卡、密码
func (c *WechatController) CancleCardPassword() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	keyvalue := c.GetString("keyvalue")
	keytype, err := c.GetInt("keytype")
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	token := c.GetString("token")
	requestid := c.GetString("requestid")

	log.Info("DoorCtrlOpen DeviceID=", roomnu, ",Token=", token, ",appid:", appid)
	if roomnu == "" || appid == "" || token == "" ||
		keyvalue == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, 0)
	if !status {
		return
	}

	//向设备服务请求发卡
	httpServerIP := fmt.Sprintf("http://%s/cancel-card-password?gwid=%s&deviceid=%s&keyvalue=%s&keytype=%d&requestid=%s",
		serverIP, gatewayID, DeviceID, keyvalue, keytype, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

}

func (c *WechatController) checkAppidUser(roomnu, appid, token string, method int) (serverIP, gatewayID, DeviceID string, gwOnline bool) {
	var devOnline bool

	//从Redis里判断该token是否存在，不存在，则没有权限访问
	_, status, err := common.RedisTokenOpt.Get(token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return "", "", "", false
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return "", "", "", false
	}

	if method == 3 {
		gatewayID, DeviceID, gwOnline, err = DBOpt.GetDataOpt().GetDoorCardInfo(roomnu, appid)
		if err != nil {
			log.Error("err:", err)
			c.Data["json"] = common.GetErrCodeJSON(10006)
			c.ServeJSON()
			return "", "", "", false
		}
	} else {
		//通过房间号与酒店appid获取设备id信息
		DeviceID, err = DBOpt.GetDataOpt().GetDeviceID(roomnu, appid)
		if err != nil {
			log.Error("err:", err)
			c.Data["json"] = common.GetErrCodeJSON(10006)
			c.ServeJSON()
			return "", "", "", false
		}
		if len(DeviceID) == 0 {
			log.Error("房间数据不存在:", roomnu, ",userid:", appid)
			c.Data["json"] = common.GetErrCodeJSON(10004)
			c.ServeJSON()
			return "", "", "", false
		}
		log.Debug("DeviceID:", DeviceID)

		//通过设备ID获取网关ID与在线状态
		gatewayID, gwOnline, devOnline, err = DBOpt.GetDataOpt().CheckGatewayOnline(DeviceID)
		if err != nil {
			log.Error("err:", err)
			c.Data["json"] = common.GetErrCodeJSON(10006)
			c.ServeJSON()
			return "", "", "", false
		}
	}

	//用Redis获取该网关连接到哪台服务器，并且或者所在连接的服务器地址
	dataBuf, isExist, err := common.RedisServerListOpt.Get(gatewayID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return "", "", "", false
	}
	if !isExist {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10008)
		c.ServeJSON()
		return "", "", "", false
	}
	serverIP = string(dataBuf)

	log.Debug("gatewayID:", gatewayID, ",gwOnline:", gwOnline)

	//目前网关心跳只是一人空包，没有网关ＩＤ，无法做到网关是否线
	devOnline = true
	var errcode int
	if gwOnline {
		//网关在线
		errcode = 0
		if devOnline {
			errcode = 0
		} else {
			//设备不在线
			errcode = 10009
		}
	} else {
		//网关不在线
		errcode = 10008
	}
	c.Data["json"] = common.GetErrCodeJSON(errcode)
	c.ServeJSON()
	if errcode != 0 {
		log.Info("网关或者设备不在线：gw=", gatewayID, ",deviceID=", DeviceID)
		return "", "", "", false
	}
	return serverIP, gatewayID, DeviceID, true
}

//@cmt clear node dev's ekey info in flash   （WechatAPI-->DeviceServer）
func (c *WechatController) ResetDev() {

	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	token := c.GetString("token")
	requestid := c.GetString("requestid")

	log.Info("DevReset: Roomnu=", roomnu, ",Token=", token, ",appid:", appid)
	if roomnu == "" || appid == "" || token == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, 0)
	if !status {
		return
	}

	//通过http发送给DeviceServer....
	httpServerIP := fmt.Sprintf("http://%s/dev-reset?gwid=%s&deviceid=%s&requestid=%s",
		serverIP, gatewayID, DeviceID, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

}

//@cmt 节点设备常开常闭
func (c *WechatController) NoncDev() {
	roomnu := c.GetString("roomnu")
	appid := c.GetString("appid")
	token := c.GetString("token")
	requestid := c.GetString("requestid")

	actionType, err := c.GetInt("actiontype")
	if err != nil {
		log.Error("actionType err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	devType, err := c.GetInt("devtype")
	if err != nil {
		log.Error("devType err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	log.Info("DevNonc: Roomnu=", roomnu, ",Token=", token, ",appid:", appid)
	if roomnu == "" || appid == "" || token == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, 0)
	if !status {
		return
	}

	//通过http发送给DeviceServer....
	httpServerIP := fmt.Sprintf("http://%s/dev-nonc-set?gwid=%s&deviceid=%s&requestid=%s&actiontype=%d&devtype=%d",
		serverIP, gatewayID, DeviceID, requestid, actionType, devType)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
}

//@cmt 设备*测试模式*
func (c *WechatController) SetTestModeDev() {
	gwid := c.GetString("gwid")
	device_mac := c.GetString("device_mac")
	requestid := c.GetString("requestid")
	tx_rate, err := c.GetInt("tx_rate")
	if err != nil {
		log.Error("tx_rate err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}
	tx_wait, err := c.GetInt("tx_wait")
	if err != nil {
		log.Error("tx_wait err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}
	log.Info("SetTestModeDev: gwid=", gwid, ",device_mac:", device_mac, ",tx_rate:", tx_rate, ",tx_wait:", tx_wait, ",requestid:", requestid)
	if gwid == "" || device_mac == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	// serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, 0)
	// if !status {
	// 	return
	// }
	//@cmt 用Redis获取该网关连接到哪台服务器，并且或者所在连接的服务器地址
	dataBuf, isExist, err := common.RedisServerListOpt.Get(gwid)
	if err != nil {
		log.Error("err:", err)
		return
	}
	if !isExist {
		log.Error("err:", err)
		return
	}
	serverIP := string(dataBuf) //get http server IP
	//通过http发送给DeviceServer....
	httpServerIP := fmt.Sprintf("http://%s/set-test-mode?gwid=%s&deviceid=%s&tx_rate=%d&tx_wait=%d&requestid=%s",
		serverIP, gwid, device_mac, tx_rate, tx_wait, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

}

//@cmt set device work mode
func (c *WechatController) SetWorkModeDev() {
	gwid := c.GetString("gwid")
	device_mac := c.GetString("device_mac")
	requestid := c.GetString("requestid")

	log.Info("SetWorkModeDev: gwid=", gwid, ",device_mac:", device_mac, ",requestid:", requestid)
	if gwid == "" || device_mac == "" || requestid == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	// serverIP, gatewayID, DeviceID, status := c.checkAppidUser(roomnu, appid, token, 0)
	// if !status {
	// 	return
	// }
	//@cmt 用Redis获取该网关连接到哪台服务器，并且或者所在连接的服务器地址
	dataBuf, isExist, err := common.RedisServerListOpt.Get(gwid)
	if err != nil {
		log.Error("err:", err)
		return
	}
	if !isExist {
		log.Error("err:", err)
		return
	}
	serverIP := string(dataBuf) //get http server IP
	//通过http发送给DeviceServer....
	httpServerIP := fmt.Sprintf("http://%s/set-work-mode?gwid=%s&deviceid=%s&requestid=%s",
		serverIP, gwid, device_mac, requestid)
	log.Debug("httpServerIP:", httpServerIP)
	resp, err := http.Get(httpServerIP)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

}
