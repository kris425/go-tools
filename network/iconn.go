package network

import "net"

type IConn interface {
	Start()

	Stop()

	GetRemoteAddr() net.Addr

	Send(data []byte) error
}
