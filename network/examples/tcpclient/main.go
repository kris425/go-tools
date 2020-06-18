package main

import (
	"fmt"
	"github.com/kris425/go-tools/network"
	"github.com/kris425/go-tools/network/iface"
	"github.com/kris425/go-tools/network/tcp"
	"log"
)

type Handle struct {
}

func (Handle) OnRecv(conn iface.IConn, data []byte) {
	log.Println("OnRecv data:", string(data))
}

func (Handle) OnNewConn(conn iface.IConn) {
	log.Println("OnNewConn")
}

func (Handle) OnClosed(conn iface.IConn) {
	log.Println("Onclosed")
}

func main() {
	client := tcp.NewClient("test", "127.0.0.1:6500",
		network.SetCodec(tcp.NewLengthCodec(4)),
		network.SetHook(&Handle{}))
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Send([]byte("hello"))
	select {}
}
