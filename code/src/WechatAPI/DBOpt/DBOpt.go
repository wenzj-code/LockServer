package DBOpt

import (
	"WechatAPI/common"
	"fmt"
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
	sqlString := fmt.Sprintf("select 1 from hotel_base_info where appid='%s' and secret='%s'", appid, secret)
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
	sqlString := fmt.Sprintf("select roomnu,appid from hotel_room_info A "+
		"inner join hotel_base_info B on A.hotel_id=B.id "+
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
	sqlString := "select device_id from hotel_room_info a,hotel_base_info b where roomnu=? and b.id=a.hotel_id and b.appid=?"
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
		"inner join hotel_room_info B on A.id=B.gw_id " +
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
	sqlString := "select A.app_domain,A.app_domain,A.appid,A.secret from hotel_base_info A " +
		"inner join hotel_room_info B on B.hotel_id=B.id " +
		"where B.device_id=?"
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
		config.URL += "/api/device/v1/notify"
		config.TokenURL += "/api/device/v1/token"
	}
	return config
}

//GetDoorCardInfo 获取门禁的信息
func (opt *DBOpt) GetDoorCardInfo(roomnu, appid string) (gatewayID string, deviceID string, status bool, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return gatewayID, deviceID, status, err
	}
	defer opt.releaseDB(conn)
	var ret int
	sqlString := "select b.gateway_id,a.device_id,a.status  from hotel_public_room a " +
		"inner join t_gateway_info b on b.id=a.gw_id" +
		"inner join hotel_base_info c on c.id=a.hotel_id and c.appid=? " +
		"where a.room=? "
	rows, err := conn.Query(sqlString, appid, roomnu)
	if err != nil {
		log.Error("err:", err)
		return gatewayID, deviceID, status, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&gatewayID, &deviceID, &ret)
		if err != nil {
			log.Error("err:", err)
			return gatewayID, deviceID, status, err
		}
	}
	if ret == 1 {
		status = true
	}
	return gatewayID, deviceID, status, err
}
