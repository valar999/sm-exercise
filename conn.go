package pool

import (
	"sync/atomic"
	"time"
)

type Connection interface {
	open()
	close()
}

type conn struct {
	addr      int32
	n         uint32
	openDelay time.Duration
}

var counter uint32 = 0

func NewConn(addr int32, openDelay time.Duration) *conn {
	return &conn{
		addr:      addr,
		openDelay: openDelay,
		n:         atomic.AddUint32(&counter, 1),
	}
}

func (c *conn) open() {
	time.Sleep(c.openDelay)
	// log.Println("open connection", c.addr)
}

func (c *conn) close() {
	// log.Println("close connection", c.addr)
}
