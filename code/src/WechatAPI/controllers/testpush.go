package controllers

import (
	"WechatAPI/common"
	"WechatAPI/config"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"
)

func (c *MainController) TestToken() {
	appid := c.GetString("appid")
	secret := c.GetString("secret")
	log.Info("appid=", appid, ",secret=", secret)
	if appid == "" || secret == "" {
		c.Data["json"] = common.GetErrCodeJSON(10003)
		c.ServeJSON()
		return
	}

	if appid != "1111" || secret != "2222" {
		c.Data["json"] = common.GetErrCodeJSON(10001)
		c.ServeJSON()
		return
	}

	data := make(map[string]interface{})
	hs := sha256.New()
	io.WriteString(hs, secret+time.Now().String())
	token := fmt.Sprintf("%x", string(hs.Sum(nil)))

	data["code"] = 0
	data["token"] = token
	data["expired_in"] = config.GetConfig().RedisTokenTimeOut
	c.Data["json"] = data
	c.ServeJSON()
}

func (c *MainController) TestPush() {
	buf, _ := ioutil.ReadAll(c.Ctx.Request.Body)
	log.Debug("push msg:\n", string(buf))
	c.Data["json"] = common.GetErrCodeJSON(0)
	c.ServeJSON()
}

/////////////////////////
