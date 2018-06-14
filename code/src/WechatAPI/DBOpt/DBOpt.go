package DBOpt

import (
	"WechatAPI/common"
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type DBOpt struct {
	BaseDB
}

var dataOpt *DBOpt
var onceDataOpt sync.Once

//GetDataOpt .获取数据平台对象
func GetDataOpt() *DBOpt {
	onceDataOpt.Do(func() {
		dataOpt = new(DBOpt)
	})
	return dataOpt
}

//CheckAppIDSecret 校验
func (opt *DBOpt) CheckAppIDSecret(appid, secret string) (status bool, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return status, err
	}
	defer opt.releaseDB(conn)
	sqlString := fmt.Sprintf("select 1 from t_user_info where appid='%s' and secret='%s'", appid, secret)
	rows, err := conn.Query(sqlString)
	if err != nil {
		log.Error("err:", err)
		return status, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return status, err
}

//GetRoomInfo 通过设备ID获取房间信息
func (opt *DBOpt) GetRoomInfo(deviceID string) (roomnu, appid string, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return roomnu, appid, err
	}
	defer opt.releaseDB(conn)
	sqlString := fmt.Sprintf("select roomnu,appid from t_device_bind_info A "+
		"inner join t_device_info B on A.device_id=B.device_id "+
		"inner join t_user_info C on B.user_id=C.id "+
		"where A.device_id='%s';", deviceID)
	rows, err := conn.Query(sqlString)
	if err != nil {
		log.Error("err:", err)
		return roomnu, appid, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&roomnu, &appid)
		if err != nil {
			log.Error("err:", err)
			return roomnu, appid, err
		}
	}
	return roomnu, appid, err
}

//GetDeviceID 通过房间号与用户ID获取设备ＩＤ
func (opt *DBOpt) GetDeviceID(roomnu string, appid string) (deviceID string, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return deviceID, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select device_id from t_device_bind_info a,t_user_info b where roomnu=? and b.id=a.user_id and b.appid=?"
	rows, err := conn.Query(sqlString, roomnu, appid)
	if err != nil {
		log.Error("err:", err)
		return deviceID, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&deviceID)
		if err != nil {
			log.Error("err:", err)
		}
	}
	return deviceID, err
}

//CheckGatewayOnline 检查设备的网关是否在线
func (opt *DBOpt) CheckGatewayOnline(deviceID string) (gatewayID string, gwStatus, devStatus bool, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return gatewayID, gwStatus, devStatus, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select A.gateway_id,A.status,B.status from t_gateway_info A " +
		"inner join t_device_info B on A.id=B.gw_id " +
		"where B.device_id=?"
	rows, err := conn.Query(sqlString, deviceID)
	if err != nil {
		log.Error("err:", err)
		return gatewayID, gwStatus, devStatus, err
	}
	defer rows.Close()
	for rows.Next() {
		var gwOnline, devOnline int
		err = rows.Scan(&gatewayID, &gwOnline, &devOnline)
		if err != nil {
			log.Error("err:", err)
			return gatewayID, gwStatus, devStatus, err
		}
		if gwOnline == 1 {
			gwStatus = true
		}
		if devOnline == 1 {
			devStatus = true
		}
	}
	return gatewayID, gwStatus, devStatus, err
}

//GetDevicePushInfo 获取推送的配置
func (opt *DBOpt) GetDevicePushInfo(deviceID string) (config common.PushConfig) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return config
	}
	defer opt.releaseDB(conn)
	sqlString := "select A.url,A.token_url,A.appid,A.secret from t_manger_pushsetting_info A " +
		"inner join t_user_info B on A.user_id=B.id " +
		"inner join t_device_info C on C.user_id=B.id " +
		"where C.device_id=?"
	rows, err := conn.Query(sqlString, deviceID)
	if err != nil {
		log.Error("err:", err)
		return config
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&config.URL, &config.TokenURL, &config.AppID, &config.Secret)
		if err != nil {
			log.Error("err:", err)
			return config
		}
	}
	return config
}

//WechatOpenMethod 微信开门方式
func (opt *DBOpt) WechatOpenMethod(deviceID string) error {
	return opt.addDoorOpenHistory(1, deviceID)
}

//CardMethod 滴卡开门方式
func (opt *DBOpt) CardMethod(deviceID string) error {
	return opt.addDoorOpenHistory(2, deviceID)
}

//KeyMethod 钥匙开门方式
func (opt *DBOpt) KeyMethod(deviceID string) error {
	return opt.addDoorOpenHistory(3, deviceID)
}

//PasswordMethod 密码开门方式
func (opt *DBOpt) PasswordMethod(deviceID string) error {
	return opt.addDoorOpenHistory(4, deviceID)
}

func (opt *DBOpt) addDoorOpenHistory(openMethod int, deviceID string) error {
	sqlString := "insert into t_device_open_info(device_id,method_id,open_time) values(?,?,?)"
	err := opt.exec(nil, sqlString, deviceID, openMethod, time.Now().Unix())
	if err != nil {
		log.Error("err:", err)
	}
	return err
}
