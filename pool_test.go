package pool

import (
	"testing"
)

func TestCache(t *testing.T) {
	pool := NewPool()
	conn1 := pool.getConnection(1)
	conn2 := pool.getConnection(1)
	if conn1 != conn2 {
		t.Error("return new connection")
	}
	pool.shutdown()
}

func TestSimultaneous(t *testing.T) {
	pool := NewPool()
	for i := 0; i < 1000; i++ {
		go func() {
			var n int32 = 0
			for i := 0; i < 1000; i++ {
				pool.getConnection(n)
				pool.onNewRemoteConnection(n, &Conn{n})
				n++
				if n >= 9 {
					n = 0
				}
			}
		}()
	}
	conn1 := pool.getConnection(1)
	conn2 := pool.getConnection(1)
	if conn1 != conn2 {
		t.Error("return new connection")
	}
	pool.shutdown()
}
