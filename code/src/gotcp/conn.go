package gotcp

import (

	//"log"

	"bufio"
	"net"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

//SendChanSize ...
const SendChanSize int = 10
const ReceiveChanSize int = 10

type ConnCallbackInterface interface {
	HandleMsg(*Conn, []byte) error
	Close()
}

type Conn struct {
	Srv            *Server
	conn           *net.TCPConn
	Lock           sync.RWMutex
	closeOnce      sync.Once
	closeChan      chan struct{}
	ReceiveChan    chan []byte
	SendChan       chan []byte
	heartTimeCount int

	//标识，哪个客户端
	clientFlag string
}

func newConn(conn *net.TCPConn, srv *Server) *Conn {
	return &Conn{
		Srv:            srv,
		conn:           conn,
		closeChan:      make(chan struct{}),
		SendChan:       make(chan []byte, SendChanSize),
		ReceiveChan:    make(chan []byte, ReceiveChanSize),
		heartTimeCount: 0,
	}
}

func (c *Conn) SetClientFlag(flag string) {
	c.clientFlag = flag
}

func (c *Conn) GetRawConn() *net.TCPConn {
	return c.conn
}

func (c *Conn) GetListenAddr() string {
	return c.conn.LocalAddr().String()
}

func (c *Conn) GetRemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		c.Srv.DeleteClientSocket(c.conn)
		close(c.closeChan)
		close(c.ReceiveChan)
		close(c.SendChan)
		c.conn.Close()
	})
}

func (c *Conn) Do() {
	go c.readLoop()
	go c.handleLoop()
	go c.writeLoop()
}

func (c *Conn) updateHeartTimer() {
	c.heartTimeCount = 0
}

func (c *Conn) readLoop() {
	defer func() {
		recover()
		c.Close()
	}()
	for {
		select {
		case <-c.closeChan:
			return
		default:
			reader := bufio.NewReader(c.conn)
			data, err := reader.ReadBytes('\n')
			if err != nil {
				if !strings.Contains(err.Error(), "EOF") {
					log.Error("strFlag:", c.clientFlag, ",read pack Eror:", err, c.GetRemoteAddr())
				}
				return
			}
			c.updateHeartTimer()
			c.ReceiveChan <- data
		}
	}
}

func (c *Conn) writeLoop() {
	defer func() {
		recover()
		c.Close()
	}()
	for {
		select {
		case <-c.closeChan:
			return
		case data := <-c.SendChan:
			if _, err := c.conn.Write(data); err != nil {
				//log.Error("clientFlag:", c.clientFlag, ",conn write err: ", err, c.GetRemoteAddr())
				return
			}
		}
	}
}

func (c *Conn) handleLoop() {
	defer func() {
		recover()
		c.Close()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		case p := <-c.ReceiveChan:
			c.Srv.callback.HandleMsg(c, p)
		}
	}
}

func (c *Conn) SetKeepAlivePeriod(d time.Duration) {
	c.conn.SetKeepAlive(true)
	c.conn.SetKeepAlivePeriod(d)
}
