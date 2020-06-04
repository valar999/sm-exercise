package pool

import "time"

type Connection interface {
	open()
	close()
}

type conn struct {
	addr      int32
	openDelay time.Duration
}

func (c *conn) open() {
	time.Sleep(c.openDelay)
	// log.Println("open connection", c.addr)
}

func (c *conn) close() {
	// log.Println("close connection", c.addr)
}
