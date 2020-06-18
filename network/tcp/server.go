package tcp

import (
	"context"
	"github.com/kris425/go-tools/network"
	"go.uber.org/atomic"
	"log"
	"net"
)

func NewServer(name string, addr string, options ...network.OptionFunc) *Server {
	opt := network.NewConnOption(nil, network.DefaultPackageLengthMax, &StringCodec{})
	for _, fn := range options {
		fn(opt)
	}
	return &Server{
		Name:       name,
		Addr:       addr,
		sessID:     atomic.NewUint64(0),
		ConnOption: opt,
	}
}

type Server struct {
	// 服务名称
	Name string

	// 监听地址
	Addr string

	// 监听器
	listener *net.TCPListener

	// session id
	sessID *atomic.Uint64

	// 钩子函数
	hook network.IHook

	// context
	context context.Context
	cancel  context.CancelFunc
	// 连接属性
	*network.ConnOption

	exitChan chan bool
}

// 开启服务并阻塞
func (s *Server) Serve() {
	log.Println("Server serve")
	defer log.Println("Serve exit")
	s.Start()
	select {
	case <-s.exitChan:
		return
	}
}

func (s *Server) Start() error {
	log.Println("Server name:", s.Name, " start at addr:", s.Addr)
	s.context, s.cancel = context.WithCancel(context.Background())
	var err error
	var addr *net.TCPAddr
	addr, err = net.ResolveTCPAddr("tcp4", s.Addr)
	if err != nil {
		log.Println("Server start error:", err.Error())
		return err
	}
	s.listener, err = net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Println("Server start error:", err.Error())
		return err
	}
	go s.handleConn()
	return nil
}

func (s *Server) Stop() {
	s.cancel()
	s.exitChan <- true
}

func (s *Server) handleConn() {
	for {
		c, err := s.listener.AcceptTCP()
		if err != nil {
			log.Println("Server handleConn error:", err.Error())
			continue
		}
		sid := s.sessID.Inc()
		conn := NewConn(s.context, c, sid, s.ConnOption)
		conn.Start()
	}
}
