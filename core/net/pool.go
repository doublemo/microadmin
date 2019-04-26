package net

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrNewFn 创建连接函数为nil
	ErrNewFn = errors.New("ErrNewFn")

	// ErrGetTimeout 获取超时
	ErrGetTimeout = errors.New("ErrGetTimeout")
)

// PoolNew 池创建函数类型
type PoolNew func() (PoolAdapter, error)

// Pool 连接池接口
type Pool interface {
	// New 连接创建方法
	New(PoolNew)

	// Get 获取
	Get() (PoolAdapter, error)

	// Put 回存
	Put(PoolAdapter)
}

// idleConn 带有超时的连接器
type idleConn struct {
	adapter PoolAdapter
	t       time.Time
}

// ConnPool 网络连接池
type ConnPool struct {
	// minConnect 最小连接数保持
	minConnect int

	// maxConnect 最大连接数
	maxConnect int

	// idleTimeout 连接最大空闲时间，超过该事件则将失效
	idleTimeout time.Duration

	// instances 实例存储通道
	instances chan *idleConn

	// newFn 创建连接对象方法
	newFn PoolNew

	// counter 数量
	counter int32

	mutex sync.RWMutex
}

func (p *ConnPool) New(fn PoolNew) {
	p.newFn = fn
}

func (p *ConnPool) Get() (PoolAdapter, error) {
	instances := p.insts()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case v := <-instances:
			if v == nil {
				continue
			}

			if timeout := p.idleTimeout; timeout > 0 {
				if v.t.Add(timeout).Before(time.Now()) && p.len() > p.minConnect {
					v.adapter.Close()
					p.mutex.Lock()
					p.counter--
					p.mutex.Unlock()
					continue
				}
			}

			if !v.adapter.Ok() {
				p.mutex.Lock()
				p.counter--
				p.mutex.Unlock()
				continue
			}

			return v.adapter, nil

		case <-ticker.C:
			return nil, ErrGetTimeout

		default:
			if p.newFn == nil {
				return nil, ErrNewFn
			}

			p.mutex.Lock()
			if p.counter >= int32(p.maxConnect) {
				p.mutex.Unlock()
				continue
			}

			p.counter++
			p.mutex.Unlock()

			adapter, err := p.newFn()
			if err != nil {
				p.mutex.Lock()
				p.counter--
				p.mutex.Unlock()
				return nil, err
			}

			return adapter, nil
		}
	}
}

func (p *ConnPool) Put(x PoolAdapter) {
	if x == nil {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	select {
	case p.instances <- &idleConn{adapter: x, t: time.Now()}:
	default:
		x.Close()
	}
}

func (p *ConnPool) insts() (instances chan *idleConn) {
	p.mutex.RLock()
	instances = p.instances
	p.mutex.RUnlock()
	return
}

func (p *ConnPool) len() int {
	instances := p.insts()
	return len(instances)
}

func NewConnPool(min, max int, idleTimeout time.Duration) *ConnPool {
	if min < 1 {
		min = 1
	}

	if min > max {
		max = min
	}

	return &ConnPool{
		minConnect:  min,
		maxConnect:  max,
		idleTimeout: idleTimeout,
		instances:   make(chan *idleConn, max),
	}
}
