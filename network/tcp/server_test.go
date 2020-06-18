package tcp

import (
	"fmt"
	"github.com/kris425/go-tools/network"
	"github.com/kris425/go-tools/network/iface"
	"testing"
)

type Handle struct {
}

func (Handle) OnRecv(conn iface.IConn, data []byte) {
	fmt.Println("OnRecv data:", string(data))
	conn.Send(data)
}

func (Handle) OnNewConn(conn iface.IConn) {
	fmt.Println("OnNewConn")
}

func (Handle) OnClosed(conn iface.IConn) {
	fmt.Println("Onclosed")
}

func TestNewServer(t *testing.T) {
	s := NewServer("test", ":6500", network.SetHook(&Handle{}))
	defer s.Stop()
	s.Serve()
}
