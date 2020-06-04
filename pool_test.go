package pool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"log"
)

const testOpenDelay = time.Millisecond * 100

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
	log.SetFlags(log.Lmicroseconds)
	m, x := 50, 50
	var c uint32
	pool := NewPool()
	var wg sync.WaitGroup
	for i := 0; i < m; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var n int32 = 0
			for i := 0; i < x; i++ {
				wg.Add(2)
				go func(n int32) {
					defer wg.Done()
					pool.getConnection(n)
					atomic.AddUint32(&c, 1)
				}(n)
				go func(n int32) {
					defer wg.Done()
					pool.onNewRemoteConnection(n, NewConn(n, testOpenDelay))
					atomic.AddUint32(&c, 1)
				}(n)
				n++
				if n >= 9 {
					n = 0
				}
				if atomic.LoadUint32(&c) > uint32(m*x/100*90) {
					pool.shutdown()
				}
			}
		}()
	}
	wg.Wait()
}

func TestNewRemote(t *testing.T) {
	log.SetFlags(log.Lmicroseconds)
	pool := NewPool()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pool.getConnection(1)
		wg.Done()
	}()
	time.Sleep(time.Millisecond * 50)
	pool.onNewRemoteConnection(1, NewConn(1, testOpenDelay))
	wg.Wait()
}
