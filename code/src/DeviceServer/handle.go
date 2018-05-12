package main

import (
	"gotcp"
)

type CallBack struct {
}

func (cb *CallBack) HandleMsg(conn *gotcp.Conn, msg []byte) error {

	return nil
}

func (cb *CallBack) Close() {
}

func baseSendMsg(conn *gotcp.Conn, msg []byte) {

	conn.SendChan <- msg
}
