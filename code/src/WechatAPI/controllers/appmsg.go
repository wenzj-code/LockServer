package controllers

import (
	"WechatAPI/DBOpt"
	"WechatAPI/common"
	"crypto/md5"
	"encoding/hex"

	log "github.com/Sirupsen/logrus"
)

//AppLogin APP登陆
func (c *MainController) AppLogin() {
	username := c.GetString("username")
	pwd := c.GetString("pwd")

	log.Debug("username:", username)
	log.Debug("pwd:", pwd)
	userInfo, err := DBOpt.GetDataOpt().GetUserPwd(username)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if userInfo.UserType != 3 {
		log.Error("只有管理用户有权限登陆")
		c.Data["json"] = common.GetErrCodeJSON(10000)
		c.ServeJSON()
		return
	}

	m := md5.New()
	m.Write([]byte(pwd))
	pwdMd5 := hex.EncodeToString(m.Sum(nil))

	if userInfo.UserPwd != pwdMd5 {
		log.Error("登陆密码不匹配:name=", username, ",pwd=", pwd, ",pwdMd5=", pwdMd5)
		c.Data["json"] = common.GetErrCodeJSON(10010)
		c.ServeJSON()
		return
	}

	dataMap := make(map[string]interface{})
	dataMap["appid"] = userInfo.AppID
	dataMap["secret"] = userInfo.Secret
	dataMap["userid"] = userInfo.UserID
	dataMap["code"] = 0
	c.Data["json"] = dataMap
	c.ServeJSON()
}

//AddGateway 添加网关
func (c *MainController) AddGateway() {
	gwid := c.GetString("gwid")
	gwname := c.GetString("gwname")
	token := c.GetString("token")
	userid, err := c.GetInt("userid")
	if err != nil {
		log.Error("err:", err)
	}

	if gwid == "" || gwname == "" || userid == 0 || token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	_, status, err := common.RedisTokenOpt.Get(token)
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

	status, err = DBOpt.GetDataOpt().CheckUserID(userid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("用户ID不存在:", userid)
		c.Data["json"] = common.GetErrCodeJSON(10012)
		c.ServeJSON()
		return
	}

	err = DBOpt.GetDataOpt().AddGatewayInfo(userid, gwid, gwname)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	c.Data["json"] = common.GetErrCodeJSON(0)
	c.ServeJSON()
}

//BindDeviceRoom 绑定房间与设备
func (c *MainController) BindDeviceRoom() {
	deviceid := c.GetString("deviceid")
	roomnu := c.GetString("roomnu")
	token := c.GetString("token")
	userid, err := c.GetInt("userid")
	if err != nil {
		log.Error("err:", err)
	}

	if deviceid == "" || roomnu == "" || userid == 0 || token == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	_, status, err := common.RedisTokenOpt.Get(token)
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

	status, err = DBOpt.GetDataOpt().CheckUserID(userid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("用户ID不存在:", userid)
		c.Data["json"] = common.GetErrCodeJSON(10012)
		c.ServeJSON()
		return
	}

	status, err = DBOpt.GetDataOpt().CheckDeviceIDRoom(deviceid, roomnu, userid)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}
	if !status {
		log.Info("房间号或者设备ID已经被绑定过了:", deviceid, ",", roomnu)
		c.Data["json"] = common.GetErrCodeJSON(10011)
		c.ServeJSON()
		return
	}

	err = DBOpt.GetDataOpt().AddDeviceAndRoomBind(userid, deviceid, roomnu)
	if err != nil {
		log.Error("err:", err)
		c.Data["json"] = common.GetErrCodeJSON(10006)
		c.ServeJSON()
		return
	}

	c.Data["json"] = common.GetErrCodeJSON(0)
	c.ServeJSON()
	return
}
