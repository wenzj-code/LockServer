package main
/*
@cmt 20180915
@描述： 用于从 work_mode_config.txt中 **批量** 导入配置 并 发送 *设置正常模式*命令 给相应节点设备
@compile: go build setDevWorkMode.go
@用法步骤:  1.部署了DeviceServer, WechatAPI
           2.配置好了  work_mode_condig.txt
           3.运行 ./setDevWorkMode [gw_name]   (e.g. ./setDevWorkMode 1AAA01000151)
@work_mode_config.txt文件的格式：
          [device_mac]
          
*/
import (
    "fmt"
    "os"
    "io"
    "io/ioutil"
    "bufio"
    "net/http"
    "strings"
    //"strconv"
    //"time"
    //"../common"
    log "../../github.com/Sirupsen/logrus"
)

func main() {
    fi, err := os.Open("./work_mode_config.txt")
    if err != nil {
        fmt.Printf("Error: %s\n", err)
        return
    }
    defer fi.Close()

    br := bufio.NewReader(fi)
    for {
        a, _, c := br.ReadLine()
        paramSlices:=strings.Split(string(a), " ")
        if c == io.EOF {
            break
		}
        devicemac:= paramSlices[0]  //device mac  
        //work_mode, err:= strconv.ParseInt(paramSlices[1], 10 ,32)
		//tx_rate, err:= strconv.ParseInt(paramSlices[2], 10 ,32)
		//tx_wait, err:= strconv.ParseInt(paramSlices[3], 10 ,32)
        // if err!=nil{
        //     return 
        // }
        fmt.Println(devicemac)
        
        gatewayID:= os.Args[1]   //@cmt 命令行参数
        // //@cmt 用Redis获取该网关连接到哪台服务器，并且或者所在连接的服务器地址
        // dataBuf, isExist, err := common.RedisServerListOpt.Get(gatewayID)
        // if err != nil {
        //     log.Error("err:", err)
        //     return 
        // }
        // if !isExist {
        //     log.Error("err:", err)
        //     return 
        // }
        // serverIP := string(dataBuf) //get http server IP
        
        //通过http发送给DeviceServer....
        httpServerIP := fmt.Sprintf("http://172.18.247.53:8990/set-work-mode?gwid=%s&deviceid=%s&requestid=1234567890123",
                                    gatewayID, devicemac )
        log.Debug("httpServerIP:", httpServerIP)
        resp, err := http.Get(httpServerIP)
        if err != nil {
            log.Error("err:", err)
            return
        }
        defer resp.Body.Close()

        _, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Error("err:", err)
            return
        }

        //time.Sleep(time.Second)
    }

}