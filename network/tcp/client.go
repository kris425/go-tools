package tcp

import (
	"context"
	"github.com/kris425/go-tools/network"
	"log"
	"net"
)

func NewClient(name string, addr string, options ...network.OptionFunc) *Client {
	opt := network.NewConnOption(nil, network.DefaultPackageLengthMax, &StringCodec{})
	for _, fn := range options {
		fn(opt)
	}
	return &Client{
		Name:       name,
		RemoteAddr: addr,
		ConnOption: opt,
	}
}

type Client struct {
	Name string

	RemoteAddr string

	conn network.IConn

	context context.Context
	cancel  context.CancelFunc
	// 连接属性
	*network.ConnOption
}

func (c *Client) Connect() error {
	c.context, c.cancel = context.WithCancel(context.Background())
	addr, err := net.ResolveTCPAddr("tcp4", c.RemoteAddr)
	if err != nil {
		log.Println("Connect error:", err)
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Println("Connect error:", err)
		return err
	}
	c.conn = NewConn(c.context, conn, 0, c.ConnOption)
	go c.conn.Start()
	return nil
}

func (c *Client) Send(data []byte) error {
	return c.conn.Send(data)
}

func (c *Client) Stop() {
	c.cancel()
}
