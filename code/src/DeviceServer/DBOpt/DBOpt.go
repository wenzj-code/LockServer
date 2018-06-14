package DBOpt

import (
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

//GetDeviceIDList 通过房间号与用户ID获取设备ＩＤ
func (opt *DBOpt) GetDeviceIDList(gatewayID string) (devListMap map[string]bool, err error) {
	conn, err := opt.connectDB()
	if err != nil {
		log.Error("err:", err)
		return devListMap, err
	}
	defer opt.releaseDB(conn)
	sqlString := "select device_id from t_device_info a,t_gateway_info b where a.gw_id=b.id and b.gateway_id=?"
	rows, err := conn.Query(sqlString, gatewayID)
	if err != nil {
		log.Error("err:", err)
		return devListMap, err
	}
	defer rows.Close()

	devListMap = make(map[string]bool)
	var deviceID string
	for rows.Next() {
		err = rows.Scan(&deviceID)
		if err != nil {
			log.Error("err:", err)
			return devListMap, err
		}
		devListMap[deviceID] = true
	}
	return devListMap, err
}
