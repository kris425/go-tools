package main

import (
	"github.com/kris425/go-tools/network"
	"github.com/kris425/go-tools/network/tcp"
	"log"
)

type Handle struct {
}

func (Handle) OnRecv(conn network.IConn, data []byte) {
	log.Println("OnRecv data:", string(data))
	conn.Send(data)
}

func (Handle) OnNewConn(conn network.IConn) {
	log.Println("OnNewConn")
}

func (Handle) OnClosed(conn network.IConn) {
	log.Println("Onclosed")
}

func main() {
	s := tcp.NewServer("test", ":6500",
		network.SetCodec(tcp.NewLengthCodec(4)),
		network.SetHook(&Handle{}))
	s.Serve()
}
