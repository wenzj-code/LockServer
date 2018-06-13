package DBOpt

import (
	"WechatAPI/common"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//GetUserPwd 获取用户的密码
func (opt *DBOpt) GetUserPwd(username string) (userInfo common.UserInfo, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return userInfo, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select id,user_pwd,appid,secret,user_type from t_user_info where user_account=?"
	rows, err := conn.Query(sqlString, username)
	if err != nil {
		log.Error("err:", err)
		return userInfo, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userInfo.UserID, &userInfo.UserPwd, &userInfo.AppID, &userInfo.Secret, &userInfo.UserType)
		if err != nil {
			log.Error("err:", err)
			return userInfo, err
		}
	}
	return userInfo, err
}

//CheckUserID 检查用户ID的合法性
func (opt *DBOpt) CheckUserID(userID int) (bool, error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return false, err
	}
	defer opt.releaseDB(conn)

	sqlString := "select 1 from t_user_info where id=?"
	rows, err := conn.Query(sqlString, userID)
	if err != nil {
		log.Error("err:", err)
		return false, opt.errOptException
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

//AddGatewayInfo 添加网关
func (opt *DBOpt) AddGatewayInfo(userID int, gatewayID, gatewayName string) error {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return opt.errDBConnect
	}
	defer opt.releaseDB(conn)

	sqlString := "insert into t_gateway_info(gateway_id,user_id,title) values(?,?,?) " +
		"on duplicate key update title=values(title)"
	err = opt.exec(conn, sqlString, gatewayID, userID, gatewayName)
	if err != nil {
		log.Error("err:", err)
		return opt.errOptException
	}
	return err
}

//CheckDeviceIDRoom 判断设备的唯一性与房间号＋设备号的唯一性
func (opt *DBOpt) CheckDeviceIDRoom(deviceID, roomNu string, userid int) (bool, error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return false, opt.errDBConnect
	}

	//该设备ID是否已被绑定，或者该用户下的房间号是否被绑定过
	sqlString := fmt.Sprintf("select ifnull(((select 1 from t_device_bind_info where device_id='%s') or "+
		"(select 1 from t_device_bind_info where user_id=%d and roomnu='%s')),0)", deviceID, userid, roomNu)

	log.Debug("sqlString:", sqlString)
	rows, err := conn.Query(sqlString)
	if err != nil {
		log.Error("err:", err)
		return false, opt.errOptException
	}
	defer rows.Close()
	if rows.Next() {
		var status int
		err = rows.Scan(&status)
		if err != nil {
			log.Error("err:", err)
			return false, opt.errOptException
		}
		if status != 0 {
			//该设备ID或者房间号已经被绑定了
			return false, nil
		}
	}

	return true, nil
}

//AddDeviceAndRoomBind 添加房间的绑定
func (opt *DBOpt) AddDeviceAndRoomBind(userID int, deviceID, roomNu string) error {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return opt.errDBConnect
	}
	defer opt.releaseDB(conn)

	sqlString := "insert into t_device_bind_info(device_id,roomnu,user_id) values(?,?,?)"
	err = opt.exec(conn, sqlString, deviceID, roomNu, userID)
	if err != nil {
		log.Error("err:", err)
		return opt.errOptException
	}
	return err
}
