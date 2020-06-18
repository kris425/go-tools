package main

import (
	"fmt"
	"github.com/kris425/go-tools/network"
	"github.com/kris425/go-tools/network/tcp"
	"log"
)

type Handle struct {
}

func (Handle) OnRecv(conn network.IConn, data []byte) {
	log.Println("OnRecv data:", string(data))
}

func (Handle) OnNewConn(conn network.IConn) {
	log.Println("OnNewConn")
}

func (Handle) OnClosed(conn network.IConn) {
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
