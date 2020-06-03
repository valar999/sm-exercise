package pool

import (
	"sync"
)

type Pool interface {
	getConnection(addr int32) Connection
	onNewRemoteConnection(remotePeer int32, conn Connection)
	shutdown()
}

type poolConn struct {
	sync.Mutex
	conn Connection
}

type connPool struct {
	sync.Mutex
	cache map[int32]*poolConn
}

func NewPool() Pool {
	return &connPool{
		cache: make(map[int32]*poolConn),
	}
}

func (pool *connPool) getConnection(addr int32) Connection {
	pool.Lock()
	c, ok := pool.cache[addr]
	if ok {
		pool.Unlock()
	} else {
		c = &poolConn{conn: &Conn{addr}}
		pool.cache[addr] = c

		c.Lock()
		pool.Unlock()
		c.conn.open()
		c.Unlock()
	}
	return c.conn
}

func (pool *connPool) onNewRemoteConnection(remotePeer int32, c Connection) {
	pool.Lock()
	_, ok := pool.cache[remotePeer]
	if !ok {
		pool.cache[remotePeer] = &poolConn{conn: c}
	}
	pool.Unlock()
}

func (pool *connPool) shutdown() {
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
