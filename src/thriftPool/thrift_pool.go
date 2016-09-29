package thriftPool

import (
	"container/list"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	CHECKINTERVAL = 60 //清除超时连接间隔
)

type ThriftDial func(ip, port string, connTimeout time.Duration) (*IdleClient, error)
type ThriftClientClose func(c *IdleClient) error

type ThriftPool struct {
	Dial  ThriftDial
	Close ThriftClientClose

	lock        *sync.Mutex
	idle        list.List
	idleTimeout time.Duration
	connTimeout time.Duration
	maxConn     uint32
	count       uint32
	ip          string
	port        string
	closed      bool
}

type IdleClient struct {
	Socket *thrift.TSocket
	Client interface{}
}

type idleConn struct {
	c *IdleClient
	t time.Time
}

var nowFunc = time.Now

//error
var (
	ErrOverMax          = errors.New("连接超过设置的最大连接数")
	ErrInvalidConn      = errors.New("Client回收时变成nil")
	ErrPoolClosed       = errors.New("连接池已经被关闭")
	ErrSocketDisconnect = errors.New("客户端socket连接已断开")
)

func NewThriftPool(ip, port string,
	maxConn, connTimeout, idleTimeout uint32,
	dial ThriftDial, closeFunc ThriftClientClose) *ThriftPool {

	thriftPool := &ThriftPool{
		Dial:        dial,
		Close:       closeFunc,
		ip:          ip,
		port:        port,
		lock:        new(sync.Mutex),
		maxConn:     maxConn,
		idleTimeout: time.Duration(idleTimeout) * time.Second,
		connTimeout: time.Duration(connTimeout) * time.Second,
		closed:      false,
		count:       0,
	}

	go thriftPool.ClearConn()

	return thriftPool
}

func (p *ThriftPool) Get() (*IdleClient, error) {
	p.lock.Lock()
	if p.closed {
		p.lock.Unlock()
		return nil, ErrPoolClosed
	}

	if p.idle.Len() == 0 && p.count >= p.maxConn {
		p.lock.Unlock()
		return nil, ErrOverMax
	}

	p.lock.Unlock()
	p.CheckTimeout()

	p.lock.Lock()

	if p.idle.Len() == 0 {
		dial := p.Dial
		p.count += 1
		p.lock.Unlock()
		client, err := dial(p.ip, p.port, p.connTimeout)
		if err != nil {
			p.lock.Lock()
			p.count -= 1
			p.lock.Unlock()
			return nil, err
		}
		if !client.Check() {
			return nil, ErrSocketDisconnect
		}
		return client, nil
	} else {
		ele := p.idle.Front()
		idlec := ele.Value.(*idleConn)
		p.idle.Remove(ele)
		p.lock.Unlock()

		if !idlec.c.Check() {
			return nil, ErrSocketDisconnect
		}
		return idlec.c, nil
	}
}

func (p *ThriftPool) Put(client *IdleClient) error {
	if client == nil {
		return ErrInvalidConn
	}

	p.lock.Lock()
	if p.closed {
		err := p.Close(client)
		client = nil
		p.lock.Unlock()
		return err
	}

	if p.count > p.maxConn {
		p.count -= 1
		err := p.Close(client)
		p.lock.Unlock()
		client = nil
		return err
	}

	if !client.Check() {
		p.count -= 1
		err := p.Close(client)
		p.lock.Unlock()
		client = nil
		return err
	}

	p.idle.PushBack(&idleConn{
		c: client,
		t: nowFunc(),
	})
	p.lock.Unlock()
	return nil
}

func (p *ThriftPool) CheckTimeout() {
	p.lock.Lock()
	for p.idle.Len() != 0 {
		ele := p.idle.Back()
		if ele == nil {
			break
		}
		v := ele.Value.(*idleConn)
		if v.t.Add(p.idleTimeout).After(nowFunc()) {
			break
		}

		//timeout && clear
		p.idle.Remove(ele)
		p.lock.Unlock()
		p.Close(v.c) //close client connection
		p.lock.Lock()
		p.count -= 1
	}
	p.lock.Unlock()

	return
}

func (c *IdleClient) SetConnTimeout(connTimeout uint32) {
	c.Socket.SetTimeout(time.Duration(connTimeout) * time.Second)
}

func (c *IdleClient) LocalAddr() net.Addr {
	return c.Socket.Conn().LocalAddr()
}

func (c *IdleClient) RemoteAddr() net.Addr {
	return c.Socket.Conn().RemoteAddr()
}

func (c *IdleClient) Check() bool {
	if c.Socket == nil || c.Client == nil {
		return false
	}
	return c.Socket.IsOpen()
}

func (p *ThriftPool) GetIdleCount() uint32 {
	return uint32(p.idle.Len())
}

func (p *ThriftPool) GetConnCount() uint32 {
	return p.count
}

func (p *ThriftPool) ClearConn() {
	for {
		p.CheckTimeout()
		time.Sleep(CHECKINTERVAL * time.Second)
	}
}

func (p *ThriftPool) Release() {
	p.lock.Lock()
	idle := p.idle
	p.idle.Init()
	p.closed = true
	p.count -= uint32(idle.Len())
	p.lock.Unlock()

	for iter := idle.Front(); iter != nil; iter = iter.Next() {
		p.Close(iter.Value.(*idleConn).c)
	}
}

func (p *ThriftPool) Recover() {
	p.lock.Lock()
	if p.closed == true {
		p.closed = false
	}
	p.lock.Unlock()
}

type MapPool struct {
	Dial  ThriftDial
	Close ThriftClientClose

	lock *sync.Mutex

	idleTimeout uint32
	connTimeout uint32
	maxConn     uint32

	pools map[string]*ThriftPool
}

func NewMapPool(maxConn, connTimeout, idleTimeout uint32,
	dial ThriftDial, closeFunc ThriftClientClose) *MapPool {

	return &MapPool{
		Dial:        dial,
		Close:       closeFunc,
		maxConn:     maxConn,
		idleTimeout: idleTimeout,
		connTimeout: connTimeout,
		pools:       make(map[string]*ThriftPool),
		lock:        new(sync.Mutex),
	}
}

func (mp *MapPool) getServerPool(ip, port string) (*ThriftPool, error) {
	addr := fmt.Sprintf("%s:%s", ip, port)
	mp.lock.Lock()
	serverPool, ok := mp.pools[addr]
	if !ok {
		mp.lock.Unlock()
		err := errors.New(fmt.Sprintf("Addr:%s thrift pool not exist", addr))
		return nil, err
	}
	mp.lock.Unlock()
	return serverPool, nil
}

func (mp *MapPool) Get(ip, port string) *ThriftPool {
	serverPool, err := mp.getServerPool(ip, port)
	if err != nil {
		addr := fmt.Sprintf("%s:%s", ip, port)
		serverPool = NewThriftPool(ip,
			port,
			mp.maxConn,
			mp.connTimeout,
			mp.idleTimeout,
			mp.Dial,
			mp.Close,
		)
		mp.lock.Lock()
		mp.pools[addr] = serverPool
		mp.lock.Unlock()
	}
	return serverPool
}

func (mp *MapPool) Release(ip, port string) error {
	serverPool, err := mp.getServerPool(ip, port)
	if err != nil {
		return err
	}

	mp.lock.Lock()
	delete(mp.pools, fmt.Sprintf("%s:%s", ip, port))
	mp.lock.Unlock()

	serverPool.Release()

	return nil
}

func (mp *MapPool) ReleaseAll() {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	for _, serverPool := range mp.pools {
		serverPool.Release()
	}
}
