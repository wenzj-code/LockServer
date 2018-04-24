package TcpServer

import (
	"gotcp"

	log "github.com/Sirupsen/logrus"
)

type DeviceServer struct {
}

func (opt *DeviceServer) HandleMsg(conn *gotcp.Conn, data []byte) error {
	log.Debug("msg:", string(data))
	return nil
}

func (opt *DeviceServer) Close() {

}
