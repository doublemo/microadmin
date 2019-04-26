package net

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
)

const (
	// SocketStatusStoped socket 服务停止状态
	SocketStatusStoped = iota

	// SocketStatusRunning socket 服务运行状态
	SocketStatusRunning
)

// Socket TCP服务实现
type Socket struct {
	// done 服务运行完成信号
	done chan error

	// exit 退出信号
	exit chan struct{}

	// listen 监听
	listen *net.TCPListener

	// callback 连接回调
	callback func(net.Conn, chan struct{})

	status int

	logger log.Logger

	wg sync.WaitGroup

	mux sync.Mutex
}

func (s *Socket) Serve(addr string, readBufferSize, writeBufferSize int) error {
	s.done = make(chan error)
	s.exit = make(chan struct{})
	defer func() {
		close(s.done)
		s.Log("socket", "close", "addr", addr)
	}()

	s.mux.Lock()
	if s.status != SocketStatusStoped {
		s.mux.Unlock()
		return errors.New("ErrorServerNonStoped")
	}
	s.status = SocketStatusRunning
	s.mux.Unlock()

	go s.serve(addr, readBufferSize, writeBufferSize)
	err := <-s.done
	close(s.exit)
	s.listen.Close()

	//waiting ...
	s.wg.Wait()

	s.mux.Lock()
	s.status = SocketStatusStoped
	s.mux.Unlock()

	return err
}

func (s *Socket) serve(addr string, readBufferSize, writeBufferSize int) {
	if err := s.listenTo(addr); err != nil {
		s.done <- err
		return
	}

	if s.callback == nil {
		s.done <- errors.New("ErrorCallBackIsNil")
		return
	}

	s.Log("socket", "on", "addr", addr)
	for {
		conn, err := s.listen.AcceptTCP()
		if err != nil {
			s.done <- err
			return
		}

		conn.SetReadBuffer(readBufferSize)
		conn.SetWriteBuffer(writeBufferSize)
		go s.client(conn)
	}
}

func (s *Socket) client(conn net.Conn) {
	defer func() {
		conn.Close()
		s.wg.Done()
	}()

	s.wg.Add(1)
	s.callback(conn, s.exit)
}

func (s *Socket) listenTo(addr string) (err error) {
	var resolveAddr *net.TCPAddr
	{
		resolveAddr, err = net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return
		}
	}

	s.listen, err = net.ListenTCP("tcp", resolveAddr)
	return
}

func (s *Socket) SetLogger(logger log.Logger) {
	s.logger = logger
}

func (s *Socket) CallBack(f func(net.Conn, chan struct{})) {
	s.callback = f
}

func (s *Socket) Log(args ...interface{}) {
	if s.logger == nil {
		return
	}

	s.logger.Log(args...)
}

func (s *Socket) Shutdown() {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.status != SocketStatusRunning {
		return
	}

	s.listen.SetDeadline(time.Now().AddDate(-1, 0, 0))
}