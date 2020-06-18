package main

import (
	"fmt"
	"github.com/kris425/go-tools/network"
	"github.com/kris425/go-tools/network/ws"
)

type Handle struct {
}

func (Handle) OnRecv(conn network.IConn, data []byte) {
	fmt.Println("OnRecv data:", string(data))
	conn.Send(data)
}

func (Handle) OnNewConn(conn network.IConn) {
	fmt.Println("OnNewConn")
}

func (Handle) OnClosed(conn network.IConn) {
	fmt.Println("Onclosed")
}

func main() {
	ws := ws.NewServer(":8881", nil, network.SetHook(&Handle{}), network.SetPackageMax(256))
	ws.Serve()
}
