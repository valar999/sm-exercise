package pool

import (
	"sync"
	"time"
)

var openDelay = time.Millisecond * 100

type Pool interface {
	getConnection(addr int32) Connection
	onNewRemoteConnection(remotePeer int32, conn Connection)
	shutdown()
}

type Conn struct {
	sync.Mutex
	ch   chan Connection
	conn Connection
}

type pool struct {
	sync.Mutex
	cache      map[int32]*Conn
	isShutdown bool
}

func NewPool() Pool {
	return &pool{
		cache: make(map[int32]*Conn),
	}
}

func (pool *pool) getConnection(addr int32) Connection {
	// Lock all Pool
	pool.Lock()
	if pool.isShutdown {
		pool.Unlock()
		return nil
	}
	// get cache item
	c, ok := pool.cache[addr]
	if ok {
		// if item exists, we can release global lock
		pool.Unlock()
		// we need the Lock here due to wait for connection appears 
		// in pool.cache[addr].conn
		c.Lock()
		c.Unlock()
	} else {
		// no cache item, we acquire Pool lock, so we are alone here
		// create new item with channel to receive connection from
		// open() or from onNewRemoteConnection
		c = &Conn{
			ch: make(chan Connection, 2),
		}
		// store new item
		pool.cache[addr] = c

		// lock this item (addr) and release global lock
		c.Lock()
		pool.Unlock()

		openConn := NewConn(addr, openDelay)
		go func(ch chan Connection) {
			// run open() in goroutine put connection to channel
			openConn.open()
			ch <- openConn
		}(c.ch)
		select {
		case conn := <-c.ch:
			// got a connection from open() or from onNew..
			// lock Pool to store connection
			pool.Lock()
			pool.cache[addr].conn = conn
			pool.Unlock()
		}
		c.Unlock()
	}
	return c.conn
}

func (pool *pool) onNewRemoteConnection(remotePeer int32, conn Connection) {
	pool.Lock()
	defer pool.Unlock()
	if pool.isShutdown {
		return
	}
	c, ok := pool.cache[remotePeer]
	if ok {
		// non-blocking send "conn" to channel
		// doesn't matter success or not
		// we need only one remote connection to accept
		// if we are in the middle of opening connection
		select {
		case c.ch <- conn:
		default:
		}
	} else {
		// just store received connection, lock acquired so it's safe
		pool.cache[remotePeer] = &Conn{conn: conn}
	}
}

func (pool *pool) shutdown() {
	pool.Lock()
	defer pool.Unlock()
	// lock all items to wait open() or other procedures
	for _, c := range pool.cache {
		c.Lock()
		// release locks on return
		defer c.Unlock()
	}
	// close all connections
	for _, c := range pool.cache {
		c.conn.close()
	}
	pool.isShutdown = true
}
