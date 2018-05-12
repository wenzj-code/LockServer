package common

import (
	"RMQ"
	Redis "RedisOpt"
)

var errCodeMap map[int]string

//RedisOpt 操作
var RedisOpt *Redis.RedisOpt

//RMQOpt 消息列表操作
var RMQOpt *RMQ.RMQOpt

func init() {
	RedisOpt = &Redis.RedisOpt{}
	RMQOpt = &RMQ.RMQOpt{}

	errCodeMap = make(map[int]string)
	errCodeMap[0] = "成功"
	errCodeMap[10001] = "secret 不正确"
	errCodeMap[10002] = "接口凭证过期"
	errCodeMap[10003] = "参数出错"
	errCodeMap[10004] = "房间号不存在"
	errCodeMap[10005] = "设备ID不存在"
	errCodeMap[10006] = "数据库服务异常"
	errCodeMap[10007] = "Redis服务异常"
	errCodeMap[10008] = "Redis服务异常"
}

//GetErrCodeJSON 获取错误信息
func GetErrCodeJSON(code int) map[string]interface{} {
	data := make(map[string]interface{})
	data["code"] = code
	data["errmsg"] = errCodeMap[code]
	return data
}
