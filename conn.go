package pool

import (
	"sync/atomic"
	"time"
)

type Connection interface {
	open()
	close()
}

type connMock struct {
	addr      int32
	n         uint32
	openDelay time.Duration
}

var counter uint32 = 0

func NewConn(addr int32, openDelay time.Duration) *connMock {
	return &connMock{
		addr:      addr,
		openDelay: openDelay,
		n:         atomic.AddUint32(&counter, 1),
	}
}

func (c *connMock) open() {
	time.Sleep(c.openDelay)
	// log.Println("open connection", c.addr)
}

func (c *connMock) close() {
	// log.Println("close connection", c.addr)
}
