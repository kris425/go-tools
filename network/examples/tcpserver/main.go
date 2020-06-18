package main

import (
	"github.com/kris425/go-tools/network"
	"github.com/kris425/go-tools/network/iface"
	"github.com/kris425/go-tools/network/tcp"
	"log"
)

type Handle struct {
}

func (Handle) OnRecv(conn iface.IConn, data []byte) {
	log.Println("OnRecv data:", string(data))
	conn.Send(data)
}

func (Handle) OnNewConn(conn iface.IConn) {
	log.Println("OnNewConn")
}

func (Handle) OnClosed(conn iface.IConn) {
	log.Println("Onclosed")
}

func main() {
	s := tcp.NewServer("test", ":6500",
		network.SetCodec(tcp.NewLengthCodec(4)),
		network.SetHook(&Handle{}))
	s.Serve()
}
