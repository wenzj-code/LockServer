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

	sqlString := "insert ignore into t_gateway_info(gateway_id,user_id,title) values(?,?,?) " +
		"on duplicate key update title=values(title)"
	err = opt.exec(conn, sqlString, gatewayID, userID, gatewayName)
	if err != nil {
		log.Error("err:", err)
		return opt.errOptException
	}
	return err
}

//CheckDeviceBeenBind 判断该该设备是否被绑定
func (opt *DBOpt) CheckDeviceBeenBind(deviceID string) (bool, error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return false, opt.errDBConnect
	}

	//该设备ID是否已被绑定，或者该用户下的房间号是否被绑定过
	sqlString := fmt.Sprintf("select 1 from t_device_info where device_id='%s'", deviceID)

	//log.Debug("sqlString:", sqlString)
	rows, err := conn.Query(sqlString)
	if err != nil {
		log.Error("err:", err)
		return false, opt.errOptException
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}

	return false, nil
}

//CheckRoomBeenBind 判断该用户的房间号是否被添加
func (opt *DBOpt) CheckRoomBeenBind(roomNu string, userid int) (bool, error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return false, opt.errDBConnect
	}

	//该设备ID是否已被绑定，或者该用户下的房间号是否被绑定过
	sqlString := fmt.Sprintf("select 1 from t_device_bind_info where user_id=%d and roomnu='%s'", userid, roomNu)

	//log.Debug("sqlString:", sqlString)
	rows, err := conn.Query(sqlString)
	if err != nil {
		log.Error("err:", err)
		return false, opt.errOptException
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}

	return false, nil
}

//GetAgentID 获取该设备所在的代理商ID
func (opt *DBOpt) GetAgentID(userid int) (agentid int, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return agentid, opt.errDBConnect
	}
	sqlString := fmt.Sprintf("select agent_id from t_user_info where id=(select parent_id from t_user_info where id=%d)", userid)
	//log.Debug("sqlString:", sqlString)
	rows, err := conn.Query(sqlString)
	if err != nil {
		log.Error("err:", err)
		return agentid, opt.errOptException
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.Scan(&agentid); err != nil {
			log.Error("err:", err)
			return agentid, err
		}
	}

	return agentid, nil
}

//CheckGatewayExist 检查对应用户下是否有该网关
func (opt *DBOpt) CheckGatewayExist(gatewayid string, userid int) (gwid int, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return gwid, opt.errDBConnect
	}
	sqlString := fmt.Sprintf("select id from t_gateway_info where gateway_id='%s' and user_id=%d", gatewayid, userid)
	//log.Debug("sqlString:", sqlString)
	rows, err := conn.Query(sqlString)
	if err != nil {
		log.Error("err:", err)
		return gwid, opt.errOptException
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.Scan(&gwid); err != nil {
			log.Error("err:", err)
			return gwid, err
		}
	}

	return gwid, nil
}

//AddDeviceAndRoomBind 添加房间的绑定
func (opt *DBOpt) AddDeviceAndRoomBind(userID, gwid int, deviceID, roomNu string) error {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return opt.errDBConnect
	}
	defer opt.releaseDB(conn)

	sqlString := "insert into t_device_info(device_name,device_id,user_id,gw_id) values(?,?,?,?)"
	err = opt.exec(conn, sqlString, "设备", deviceID, userID, gwid)
	if err != nil {
		log.Error("err:", err)
		return opt.errOptException
	}

	sqlString = "insert into t_device_bind_info(device_id,roomnu,user_id) values(?,?,?)"
	err = opt.exec(conn, sqlString, deviceID, roomNu, userID)
	if err != nil {
		log.Error("err:", err)
		return opt.errOptException
	}
	return err
}
