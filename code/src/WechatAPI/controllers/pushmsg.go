package controllers

import (
	"WechatAPI/DBOpt"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

///////接收通知//////////
//DoorRecvReport 接收数据上报
func (c *MainController) DoorRecvReport() {
	deviceID := c.GetString("deviceid")
	barry := c.GetString("barry")
	status := c.GetString("status")

	pushConfig := DBOpt.GetDataOpt().GetDevicePushInfo(deviceID)
	if len(pushConfig.URL) < 10 {
		log.Error("还没配置推送地址，不推送:", deviceID)
		return
	}
	log.Debug("config:", pushConfig)

	roomNum, _, err := DBOpt.GetDataOpt().GetRoomInfo(deviceID)
	if err != nil {
		log.Error("err:", err)
		return
	}

	data := make(map[string]interface{})
	dataMap := make(map[string]interface{})

	data["Barry"] = barry
	data["Status"] = status
	dataMap["DeviceID"] = deviceID
	dataMap["RoomNum"] = roomNum
	dataMap["Data"] = data

	dataBuf, err := json.Marshal(dataMap)
	if err != nil {
		log.Error("err:", err)
		return
	}

	token, err := getToken(pushConfig.TokenURL, pushConfig.AppID, pushConfig.Secret)
	if err != nil {
		log.Error("err:", err)
		return
	}
	err = pushMsg(pushConfig.URL, token, dataBuf)
	if err != nil {
		log.Error("err:", err)
	} else {
		log.Info("推送成功:", deviceID)
	}
	return
}

func pushMsg(url, token string, msg []byte) error {
	var i int
	for i = 0; i < 4; i++ {
		// tr := &http.Transport{
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// }
		// client := &http.Client{Transport: tr}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		data, err1 := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err1 != nil {
			time.Sleep(100 * time.Millisecond)
			log.Error("err:", err)
			continue
		}
		dataInfo := make(map[string]interface{})
		err = json.Unmarshal(data, &dataInfo)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			log.Error("err:", err)
			continue
		}

		return nil
	}
	if i == 4 {
		return errors.New("推送失败")
	}
	return nil
}

func getToken(tokenURL, appid, secret string) (string, error) {
	dataInfo := make(map[string]interface{})
	var i int
	for i = 0; i < 4; i++ {
		reqURL := tokenURL + "?appid=" + appid + "&secret=" + secret
		log.Debug("获取token地址:", reqURL)

		// tr := &http.Transport{
		// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// }
		// client := &http.Client{Transport: tr}

		// resp, err := client.Get(reqURL)
		resp, err := http.Get(reqURL)
		if err != nil {
			log.Error("err:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		data, err1 := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err1 != nil {
			log.Error("err:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		log.Debug("请求token内容：", string(data))
		err = json.Unmarshal(data, &dataInfo)
		if err != nil {
			log.Error("err:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	if i == 4 {
		return "", errors.New("token err")
	}
	statusValue, ok := dataInfo["code"]
	if !ok {
		return "", errors.New("token请求没有状态这个字段code")
	}
	status := statusValue.(float64)
	if status != 0 {
		return "", errors.New(fmt.Sprint("错误代码:", status))
	}
	token, ok := dataInfo["token"]
	if !ok {
		return "", errors.New("token请求没有状态这个字段token")
	}
	return token.(string), nil
}
