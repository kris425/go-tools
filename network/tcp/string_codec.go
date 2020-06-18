package tcp

import "github.com/kris425/go-tools/network/iface"

func NewStringCodec() *StringCodec {
	return &StringCodec{}
}

type StringCodec struct {
}

func (StringCodec) Encode(data []byte) ([]byte, error) {
	return data, nil
}

func (StringCodec) Decode(conn iface.IConn) ([]byte, error) {
	c := conn.(*Conn)
	data := make([]byte, c.PackageMax)
	var n int
	var err error
	if n, err = c.Conn.Read(data); err != nil {
		return nil, err
	}
	return data[:n], nil
}
