package HTTPServer

import (
	"DeviceServer/Handle"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

/*
HttpInit Http服务的初始化
*/
func HTTPInit(HTTPAddrPort string) error {
	log.Info("httpserver start:", HTTPAddrPort)
	http.HandleFunc("/dev-ctrl", httpServerFuncDevCtrl)
	http.HandleFunc("/cancel-card-password", httpServerFuncCancelCard)
	http.HandleFunc("/setting-card-password", httpServerFuncSettingCard)
	http.HandleFunc("/dev-reset", httpServerFuncResetDev) //@cmt 
	http.HandleFunc("/dev-nonc-set", httpServerFuncNoncDev) //@cmt
	http.HandleFunc("/set-mode", httpServerFuncSetTestMode) //@cmt
	err := http.ListenAndServe(HTTPAddrPort, nil)
	if err != nil {
		log.Error("err:", err)
		os.Exit(0)
	}
	return err
}

//http路由回调函数
func httpServerFuncDevCtrl(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Error("err:", err)
		return
	}

	log.Debug("value:", req.Form)
	io.WriteString(w, "recv ok")

	gwid, isExist := req.Form["gwid"]
	if !isExist {
		log.Error("gwid 字段不存在:", req.Form)
		return
	}

	deviceid, isExist := req.Form["deviceid"]
	if !isExist {
		log.Error("deviceid 字段不存在:", req.Form)
		return
	}

	requestid, isExist := req.Form["requestid"]
	if !isExist {
		log.Error("requestid 字段不存在:", req.Form)
		return
	}

	conn, isExist := Handle.ConnInfo[gwid[0]]
	if !isExist {
		log.Error("该网关不在线:", gwid)
		return
	}

	//开门控制,转发到对应的网关
	Handle.DevCtrl(conn, gwid[0], deviceid[0], requestid[0])
}

//取消发卡
func httpServerFuncCancelCard(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Error("err:", err)
		return
	}

	log.Debug("value:", req.Form)
	io.WriteString(w, "recv ok")

	gwid, isExist := req.Form["gwid"]
	if !isExist {
		log.Error("gwid 字段不存在:", req.Form)
		return
	}

	requestid, isExist := req.Form["requestid"]
	if !isExist {
		log.Error("requestid 字段不存在:", req.Form)
		return
	}

	deviceid, isExist := req.Form["deviceid"]
	if !isExist {
		log.Error("deviceid 字段不存在:", req.Form)
		return
	}
	keyvalue, isExist := req.Form["keyvalue"]
	if !isExist {
		log.Error("keyvalue 字段不存在:", req.Form)
		return
	}
	keytype, isExist := req.Form["keytype"]
	if !isExist {
		log.Error("keytype 字段不存在:", req.Form)
		return
	}
	keytypeFloat, err := strconv.ParseFloat(keytype[0], 32)
	if err != nil {
		log.Error("err:", err)
		return
	}

	conn, isExist := Handle.ConnInfo[gwid[0]]
	if !isExist {
		log.Error("该网关不在线:", gwid)
		return
	}

	//开门控制,转发到对应的网关
	Handle.DevCancelPassword(conn, deviceid[0], keyvalue[0], requestid[0], int(keytypeFloat))
}

//发卡
func httpServerFuncSettingCard(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Error("err:", err)
		return
	}

	log.Debug("value:", req.Form)
	io.WriteString(w, "recv ok")

	gwid, isExist := req.Form["gwid"]
	if !isExist {
		log.Error("gwid 字段不存在:", req.Form)
		return
	}
	requestid, isExist := req.Form["requestid"]
	if !isExist {
		log.Error("requestid 字段不存在:", req.Form)
		return
	}
	deviceid, isExist := req.Form["deviceid"]
	if !isExist {
		log.Error("deviceid 字段不存在:", req.Form)
		return
	}

	keyvalue, isExist := req.Form["keyvalue"]
	if !isExist {
		log.Error("keyvalue 字段不存在:", req.Form)
		return
	}
	keytype, isExist := req.Form["keytype"]
	if !isExist {
		log.Error("keytype 字段不存在:", req.Form)
		return
	}
	keytypeFloat, err := strconv.ParseFloat(keytype[0], 32)
	if err != nil {
		log.Error("err:", err)
		return
	}

	expireDate, isExist := req.Form["expire-date"]
	if !isExist {
		log.Error("expire-date 字段不存在:", req.Form)
		return
	}
	expireDateFloat, err := strconv.ParseFloat(expireDate[0], 64)
	if err != nil {
		log.Error("err:", err)
		return
	}
	dateTime := time.Unix(int64(expireDateFloat), 0).Format("2006-01-02 15:04:05")

	conn, isExist := Handle.ConnInfo[gwid[0]]
	if !isExist {
		log.Error("该网关不在线:", gwid)
		return
	}

	//开门控制,转发到对应的网关
	Handle.DevSettingPassword(conn, deviceid[0], keyvalue[0], dateTime, requestid[0], int(keytypeFloat))
}


//@cmt 清除节点的卡号密码信息
func httpServerFuncResetDev(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Error("err:", err)
		return
	}

	log.Debug("value:", req.Form)
	io.WriteString(w, "dev-reset recv ok")  //DeviceServer-->WechatAPI

	gwid, isExist := req.Form["gwid"]
	if !isExist {
		log.Error("gwid 字段不存在:", req.Form)
		return
	}

	requestid, isExist := req.Form["requestid"]
	if !isExist {
		log.Error("requestid 字段不存在:", req.Form)
		return
	}

	deviceid, isExist := req.Form["deviceid"]
	if !isExist {
		log.Error("deviceid 字段不存在:", req.Form)
		return
	}

	conn, isExist := Handle.ConnInfo[gwid[0]]
	if !isExist {
		log.Error("该网关不在线:", gwid)
		return
	}

	//*清除节点卡号密码信息*命令,转发到对应的网关
	Handle.DevResetEkeyInfo(conn, deviceid[0], requestid[0] )

}


//@cmt 收到从WechatAPI发来的 *设备常开常闭* 命令
func httpServerFuncNoncDev(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Error("err:", err)
		return
	}

	log.Debug("value:", req.Form)
	io.WriteString(w, "dev-nonc recv ok")  //DeviceServer-->WechatAPI

	gwid, isExist := req.Form["gwid"]
	if !isExist {
		log.Error("gwid 字段不存在:", req.Form)
		return
	}

	requestid, isExist := req.Form["requestid"]
	if !isExist {
		log.Error("requestid 字段不存在:", req.Form)
		return
	}
	//@cmt actiontype 和devtype 这两个字段现在还没有用到..
	actionType, isExist :=req.Form["actiontype"]
	if !isExist {
		log.Error("actiontype 字段不存在:", req.Form)
		return
	}
	actionTypeFloat, err := strconv.ParseFloat(actionType[0], 32)
	if err != nil {
		log.Error("err:", err)
		return
	}


	deviceid, isExist := req.Form["deviceid"]
	if !isExist {
		log.Error("deviceid 字段不存在:", req.Form)
		return
	}
	conn, isExist := Handle.ConnInfo[gwid[0]]
	if !isExist {
		log.Error("该网关不在线:", gwid)
		return
	}

	//*设备常开常闭*命令,转发到对应的网关
	Handle.DevNoncSet(conn, deviceid[0], requestid[0], int(actionTypeFloat) )

}


//@cmt 设置节点 “测试模式”
func httpServerFuncSetTestMode(w http.ResponseWriter, req *http.Request){
	err := req.ParseForm()
	if err != nil {
		log.Error("err:", err)
		return
	}

	log.Debug("value:", req.Form)
	io.WriteString(w, "dev set-mode, recv ok")  //DeviceServer-->WechatAPI

	// requestid, isExist := req.Form["requestid"]
	// if !isExist {
	// 	log.Error("requestid 字段不存在:", req.Form)
	// 	return
	// }

	workMode, isExist :=req.Form["work_mode"]
	if !isExist {
		log.Error("work_mode 字段不存在:", req.Form)
		return
	}
	workModeFloat, err := strconv.ParseFloat(workMode[0], 32) //work_mode
	if err != nil {
		log.Error("err:", err)
		return
	}
	txRate, isExist := req.Form["tx_rate"]
	if !isExist {
		log.Error("tx_rate 字段不存在:", req.Form)
		return
	}	
	txRateFloat,err := strconv.ParseFloat(txRate[0], 32)  //tx_rate

	txWait, isExist := req.Form["tx_wait"]
	if !isExist {
		log.Error("tx_wait 字段不存在:", req.Form)
		return
	}
	txWaitFloat,err := strconv.ParseFloat(txWait[0], 32)  //tx_wait

	deviceid, isExist := req.Form["deviceid"]
	if !isExist {
		log.Error("deviceid 字段不存在:", req.Form)
		return
	}

	gwid, isExist := req.Form["gwid"]
	if !isExist {
		log.Error("gwid 字段不存在:", req.Form)
		return
	}
	conn, isExist := Handle.ConnInfo[gwid[0]]
	if !isExist {
		log.Error("该网关不在线:", gwid)
		return
	}


	//*设置节点模式*命令,转发到对应的网关
	Handle.DevSetMode(conn, gwid[0], deviceid[0], int(workModeFloat), int(txRateFloat), int(txWaitFloat) )
}