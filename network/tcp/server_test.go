package tcp

import (
	"fmt"
	"github.com/kris425/go-tools/network"
	"testing"
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

func TestNewServer(t *testing.T) {
	s := NewServer("test", ":6500", network.SetHook(&Handle{}))
	defer s.Stop()
	s.Serve()
}
