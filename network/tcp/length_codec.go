package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/kris425/go-tools/network/iface"
	"io"
	"log"
)

var (
	ErrPackageInvalid = errors.New("package length invalid")
	ErrPackageTooLong = errors.New("package length too long")
)

func NewLengthCodec(l int) *LengthCodec {
	return &LengthCodec{
		headLength: l,
	}
}

type LengthCodec struct {
	headLength int
}

// 发送数据编码
func (c *LengthCodec) Encode(data []byte) ([]byte, error) {
	length := int32(len(data))
	pkg := new(bytes.Buffer)
	var err error
	err = binary.Write(pkg, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg, binary.BigEndian, data)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// 接收数据解码
func (c *LengthCodec) Decode(conn iface.IConn) ([]byte, error) {
	tcpConn, ok := conn.(*Conn)
	if !ok {
		return nil, errors.New("invalid type")
	}
	headData := make([]byte, c.headLength)
	if _, err := io.ReadFull(tcpConn.Conn, headData); err != nil {
		log.Println("[Decode] read error:", err.Error())
		return nil, err
	}
	headBuf := bytes.NewReader(headData)
	var err error
	var length int32
	err = binary.Read(headBuf, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, ErrPackageInvalid
	}
	if length > int32(tcpConn.PackageMax) {
		return nil, ErrPackageTooLong
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(tcpConn.Conn, data); err != nil {
		return nil, err
	}
	return data, nil
}
