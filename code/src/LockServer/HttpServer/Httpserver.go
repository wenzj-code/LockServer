package HttpServer

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func httpCallBack(w http.ResponseWriter, req *http.Request) {

}

//InitHTTPServer 初始化HTTP请求
func InitHTTPServer(httpAddr string) error {
	http.HandleFunc("/lock-control", httpCallBack)
	err := http.ListenAndServe(httpAddr, nil)
	if err != nil {
		log.Error("err:", err)
		return err
	}

	return nil
}
