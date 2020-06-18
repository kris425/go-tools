package tcp

import (
	"context"
	"errors"
	"github.com/kris425/go-tools/network"
	"go.uber.org/atomic"
	"io"
	"log"
	"net"
)

func NewConn(ctx context.Context, c *net.TCPConn, sessID uint64, opt *network.ConnOption) *Conn {
	return &Conn{
		context:      ctx,
		Conn:         c,
		SessionID:    sessID,
		isClosed:     atomic.NewBool(false),
		exitChan:     make(chan bool, 1),
		sendBuffChan: make(chan []byte, 32),
		ConnOption:   opt,
	}
}

// TCP连接对象
type Conn struct {
	context context.Context
	// TCP连接
	Conn *net.TCPConn

	// 连接sessionid
	SessionID uint64

	// 是否关闭
	isClosed *atomic.Bool

	// 通知退出的channel
	exitChan chan bool

	// 缓存消息channel
	sendBuffChan chan []byte

	// 连接熟悉
	*network.ConnOption
}

// 启动Conn连接
func (c *Conn) Start() {
	c.Conn.SetNoDelay(true)
	go c.startRead()
	go c.startWrite()
	if c.Hook != nil {
		c.Hook.OnNewConn(c)
	}
}

// 关闭Conn连接
func (c *Conn) Stop() {
	//如果当前链接已经关闭
	if c.isClosed.Load() == true {
		return
	}
	log.Println("Conn Stop()...ConnID = ", c.SessionID)
	c.isClosed.Store(true)
	if c.Hook != nil {
		c.Hook.OnClosed(c)
	}
	c.Conn.Close()
	c.exitChan <- true
	close(c.sendBuffChan)
	close(c.exitChan)
}

// 获取对端地址信息
func (c *Conn) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Conn) Send(data []byte) error {
	if c.isClosed.Load() {
		return errors.New("conn is closed")
	}
	var err error
	data, err = c.Codec.Encode(data)
	if err != nil {
		return err
	}
	c.sendBuffChan <- data
	return nil
}

func (c *Conn) startWrite() {
	log.Println("[startWrite] Writer goroutine is running")
	defer func() {
		log.Println(c.GetRemoteAddr().String(), "[startWriterChan] conn Writer exit!")
		c.Stop()
	}()
	for {
		select {
		case <-c.context.Done():
			log.Println("context done")
			return
		case data := <-c.sendBuffChan:
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("[startWrite] Write data error = [", err, "] Writer goroutine exit")
				return
			}
		case <-c.exitChan:
			return
		}
	}
}

func (c *Conn) startRead() {
	log.Println("[startRead] Reader Goroutine is running")
	defer c.Stop()
	for {
		var (
			data []byte
			err  error
		)
		data, err = c.Codec.Decode(c)
		if err != nil {
			if err != io.EOF {
				log.Println("startRead error:", err.Error())
			}
			break
		}

		if c.Hook != nil {
			c.Hook.OnRecv(c, data)
		}
	}
	log.Println(c.GetRemoteAddr().String(), "[startReaderChan] conn Reader exit!")
}
