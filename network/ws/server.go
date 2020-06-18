package ws

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	. "github.com/kris425/go-tools/network"
	"github.com/kris425/go-tools/zaplog"
	"log"
	"net/http"
)

func NewServer(addr string, r *gin.Engine, options ...OptionFunc) *Server {
	if r == nil {
		r = gin.Default()
	}
	opt := NewConnOption(nil, DefaultPackageLengthMax, nil)
	for _, fn := range options {
		fn(opt)
	}
	s := &Server{
		addr:       addr,
		router:     r,
		clients:    make(map[*Conn]bool),
		register:   make(chan *Conn),
		unregister: make(chan *Conn),
		dataChan:   make(chan []byte, 1024),
		ConnOption: opt,
	}
	s.upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	s.context, s.cancel = context.WithCancel(context.Background())
	return s
}

type Server struct {
	// Listen addr ip:port
	addr string
	// gin router
	router *gin.Engine
	// websocket upgrader
	upgrader *websocket.Upgrader
	// clients
	clients map[*Conn]bool

	logger *zaplog.Logger

	context context.Context

	cancel context.CancelFunc

	// Register requests from the clients.
	register chan *Conn

	// Unregister requests from clients.
	unregister chan *Conn

	dataChan chan []byte

	*ConnOption
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

// 开启服务
func (s *Server) Serve() error {
	log.Println("websocket service serve at,", s.addr)
	go s.run()
	s.router.GET("/", func(ctx *gin.Context) {
		c, err := s.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		conn := NewConn(s, c, s.context)
		s.register <- conn
		go conn.Start()
	})
	err := http.ListenAndServe(s.addr, s.router)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() {
	s.cancel()
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.sendChan)
			}
		}
	}
}
