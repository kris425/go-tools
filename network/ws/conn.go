package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
	"log"
	"net"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 65535
)

func NewConn(s *Server, c *websocket.Conn, ctx context.Context) *Conn {
	return &Conn{
		server:   s,
		conn:     c,
		sendChan: make(chan []byte, 256),
		context:  ctx,
	}
}

type Conn struct {
	server *Server
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	sendChan chan []byte

	context context.Context

	// 是否关闭
	isClosed *atomic.Bool
}

func (c *Conn) startRead() {
	defer func() {
		c.server.unregister <- c
		log.Println("Conn closed")
		c.Stop()
	}()
	c.conn.SetReadLimit(int64(c.server.ConnOption.PackageMax))
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if c.server.Hook != nil {
			c.server.Hook.OnRecv(c, message)
		}
	}
}

func (c *Conn) startWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Stop()
	}()
	for {
		select {
		case <-c.context.Done():
			return
		case message, ok := <-c.sendChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.sendChan)
			for i := 0; i < n; i++ {
				w.Write(<-c.sendChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Conn) Start() {
	go c.startRead()
	go c.startWrite()
	if c.server.Hook != nil {
		c.server.Hook.OnNewConn(c)
	}
}

func (c *Conn) Stop() {
	if c.isClosed.Load() {
		return
	}
	c.conn.Close()
	if c.server.Hook != nil {
		c.server.Hook.OnClosed(c)
	}
}

func (c *Conn) GetRemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) Send(data []byte) error {
	c.sendChan <- data
	return nil
}
