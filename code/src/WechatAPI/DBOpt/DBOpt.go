package DBOpt

import (
	"WechatAPI/common"
	"sync"

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
	sqlString := "select 1 from t_user_info where appid=? and secret=?"
	rows, err := conn.Query(sqlString, appid, secret)
	if err != nil {
		log.Error("err:", err)
		return status, err
	}
	defer rows.Close()
	for rows.Next() {
		status = true
	}
	return status, err
}

//GetRoomInfo 通过设备ID获取房间信息
func (opt *DBOpt) GetRoomInfo(deviceID string) (roomnu, userAccount string, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return roomnu, userAccount, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select roomnu,user_account from t_device_bind_info A " +
		"inner join t_device_info B on A.device_id=B.device_id " +
		"inner join t_user_info C on B.user_id=C.id " +
		"where A.device_id='';"
	rows, err := conn.Query(sqlString, deviceID)
	if err != nil {
		return roomnu, userAccount, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&roomnu, &userAccount)
		if err != nil {
			log.Error("err:", err)
			return roomnu, userAccount, err
		}
	}
	return roomnu, userAccount, err
}

//CheckGatewayOnline 检查设备的网关是否在线
func (opt *DBOpt) CheckGatewayOnline(deviceID string) (gatewayID string, status bool, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return gatewayID, status, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select A.device_id,A.status from t_gateway_info A " +
		"inner join t_device_info B on A.id=B.gw_id " +
		"where B.device_id=?"
	rows, err := conn.Query(sqlString, deviceID)
	if err != nil {
		log.Error("err:", err)
		return gatewayID, status, err
	}
	defer rows.Close()
	for rows.Next() {
		var doorStatus int
		err = rows.Scan(&gatewayID, &doorStatus)
		if err != nil {
			log.Error("err:", err)
			return gatewayID, status, err
		}
		if doorStatus == 1 {
			status = true
		}
	}
	return gatewayID, status, err
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
