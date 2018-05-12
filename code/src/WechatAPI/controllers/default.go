package controllers

import (
	"WechatAPI/DBOpt"
	"WechatAPI/common"
	"WechatAPI/config"
	"encoding/json"

	"crypto/sha256"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

//MainController .
type MainController struct {
	beego.Controller
}

//GetToken 通过APPID+Secrete生成token
func (c *MainController) GetToken() {
	appid := c.GetString("appid")
	secret := c.GetString("secret")
	log.Info("appid=", appid, ",secret=", secret)
	if appid == "" || secret == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	status, err := DBOpt.GetDataOpt().CheckAppIDSecret(appid, secret)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}
	//组织应答报文
	data := make(map[string]interface{})
	hs := sha256.New()
	io.WriteString(hs, secret+time.Now().String())
	token := fmt.Sprintf("%x", string(hs.Sum(nil)))

	common.RedisOpt.Set(token, 1, config.GetConfig().RedisTimeOut)

	data["code"] = 0
	data["token"] = token
	data["expired_in"] = config.GetConfig().RedisTimeOut
	c.Data["json"] = data
	c.ServeJSON()
}

//GetRoomInfo 通过设备ID获取房间信息
func (c *MainController) GetRoomInfo() {
	DeviceID := c.GetString("deviceid")
	Token := c.GetString("token")
	log.Info("DeviceID=", DeviceID, ",Token=", Token)
	if DeviceID == "" || Token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	TokenVal, status, err := common.RedisOpt.Get(Token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}
	if string(TokenVal) != Token {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	roomnu, userAccount, err := DBOpt.GetDataOpt().GetRoomInfo(DeviceID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	data := make(map[string]interface{})
	data["roomnu"] = roomnu
	data["userid"] = userAccount
	data["code"] = 0
	c.Data["json"] = data
	c.ServeJSON()
	return
}

//DoorCtrlOpen 开门
func (c *MainController) DoorCtrlOpen() {
	DeviceID := c.GetString("deviceid")
	Token := c.GetString("token")
	log.Info("DeviceID=", DeviceID, ",Token=", Token)
	if DeviceID == "" || Token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	TokenVal, status, err := common.RedisOpt.Get(Token)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10007)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}
	if string(TokenVal) != Token {
		log.Info("Token数据不存在")
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	gatewayID, onlineStatus, err := DBOpt.GetDataOpt().CheckGatewayOnline(DeviceID)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	dataCtrlMap := make(map[string]interface{})
	dataCtrlMap["GatewayID"] = gatewayID
	dataCtrlMap["DeviceID"] = DeviceID
	dataCtrlBuffer, _ := json.Marshal(&dataCtrlMap)

	common.RMQOpt.PublishTopic(dataCtrlBuffer)

	data := make(map[string]interface{})
	if onlineStatus {
		data["code"] = common.GetErrCodeJSON(0)
	} else {
		data["code"] = common.GetErrCodeJSON(10008)
	}
	c.Data["json"] = data
	c.ServeJSON()
	return
}
