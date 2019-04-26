package net

import (
	"sync"
	"errors"
	"time"
	"net"

	"github.com/go-kit/kit/log"
	kcp "github.com/xtaci/kcp-go"
)

const (
	// KCPStatusStoped KCP 服务停止状态
	KCPStatusStoped = iota

	// KCPStatusRunning KCP 服务运行状态
	KCPStatusRunning
)

type KCP struct {
	// ops 配置信息
	ops *KCPOptions

	// done 服务运行完成信号
	done chan error

	// exit 退出信号
	exit chan struct{}

	// listen 监听
	listen *kcp.Listener

	// callback 连接回调
	callback func(net.Conn, chan struct{})

	status int

	logger log.Logger

	wg sync.WaitGroup

	mux sync.Mutex
}

func (k *KCP) Serve(ops *KCPOptions) error {
	k.ops = ops
	k.done = make(chan error)
	k.exit = make(chan struct{})
	defer func() {
		close(k.done)
		k.Log("socket", "close", "addr", "")
	}()

	k.mux.Lock()
	if k.status != KCPStatusStoped {
		k.mux.Unlock()
		return errors.New("ErrorServerNonStoped")
	}
	k.status = KCPStatusRunning
	k.mux.Unlock()

	go k.serve()
	err := <-k.done
	close(k.exit)
	k.listen.Close()

	//waiting ...
	k.wg.Wait()

	k.mux.Lock()
	k.status = KCPStatusStoped
	k.mux.Unlock()

	return err
}

func (k *KCP) serve() {
	if err := k.listenTo(); err != nil {
		k.done <- err
		return
	}

	if k.callback == nil {
		k.done <- errors.New("ErrorCallBackIsNil")
		return
	}

	k.Log("kcp", "on", "addr", k.listen.Addr())
	for{
		conn, err := k.listen.AcceptKCP()
		if err != nil {
			k.done <- err
			return
		}

		conn.SetStreamMode(true)
		conn.SetWriteDelay(false)
		conn.SetNoDelay(k.ops.Delay())
		conn.SetMtu(k.ops.MTU)
		conn.SetWindowSize(k.ops.SndWnd, k.ops.RcvWnd)
		conn.SetACKNoDelay(k.ops.AckNodelay)
		
		go k.client(conn)
	}
}

func (k *KCP) listenTo() (err error) {
	block, err := k.ops.BlockCrypt()
	if err != nil {
		return err
	}

	k.listen, err = kcp.ListenWithOptions(k.ops.Addr, block, k.ops.DataShard, k.ops.ParityShard)
	if err != nil {
		return
	}

	if err := k.listen.SetDSCP(k.ops.DSCP); err != nil {
		return err
	}

	if err := k.listen.SetReadBuffer(k.ops.SockBuf); err != nil {
		return err
	}

	if err := k.listen.SetWriteBuffer(k.ops.SockBuf); err != nil {
		return err
	}
	return nil
}

func (k *KCP) client(conn net.Conn) {
	defer func() {
		conn.Close()
		k.wg.Done()
	}()

	k.wg.Add(1)
	k.callback(conn, k.exit)
}

func (k *KCP) Log(args ...interface{}) {
	if k.logger == nil {
		return
	}

	k.logger.Log(args...)
}

func (k *KCP) Shutdown() {
	k.mux.Lock()
	defer k.mux.Unlock()

	if k.status != KCPStatusRunning {
		return
	}

	k.listen.SetDeadline(time.Now().AddDate(-1, 0, 0))
}

func (k *KCP) SetLogger(logger log.Logger) {
	k.logger = logger
}

func (k *KCP) CallBack(f func(net.Conn, chan struct{})) {
	k.callback = f
}