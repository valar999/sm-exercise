package pool

import (
	"sync"
)

type Pool interface {
	getConnection(addr int32) Connection
	onNewRemoteConnection(remotePeer int32, conn Connection)
	shutdown()
}

type Conn struct {
	sync.Mutex
	conn Connection
}

type pool struct {
	sync.Mutex
	cache map[int32]*Conn
}

func NewPool() Pool {
	return &pool{
		cache: make(map[int32]*Conn),
	}
}

func (pool *pool) getConnection(addr int32) Connection {
	pool.Lock()
	c, ok := pool.cache[addr]
	if ok {
		pool.Unlock()
	} else {
		c = &Conn{conn: &conn{addr}}
		pool.cache[addr] = c

		c.Lock()
		pool.Unlock()
		c.conn.open()
		c.Unlock()
	}
	return c.conn
}

func (pool *pool) onNewRemoteConnection(remotePeer int32, c Connection) {
	pool.Lock()
	_, ok := pool.cache[remotePeer]
	if !ok {
		pool.cache[remotePeer] = &Conn{conn: c}
	}
	pool.Unlock()
}

func (pool *pool) shutdown() {
	pool.Lock()
	defer pool.Unlock()
	for _, c := range pool.cache {
		c.Lock()
		defer c.Unlock()
	}
	for _, c := range pool.cache {
		c.conn.close()
	}
}
