package net

import (
	"math"
	"testing"
	"time"
)

type TestPoolAdapter struct {
	id int
}

func (tpa *TestPoolAdapter) Close() {}

func (tpa *TestPoolAdapter) Ok() bool {
	return true
}

func TestConnPool(t *testing.T) {
	c := NewConnPool(1, 300, time.Second*60)
	var (
		m              int
		counterSuccess int
		counterGetFail int
		threadNum      = 1000
	)

	c.New(func() (PoolAdapter, error) {
		m++
		return &TestPoolAdapter{id: m}, nil
	})

	for n := 0; n < threadNum; n++ {
		go func(o int) {
			for i := 0; i < 1; i++ {
				begin := time.Now()
				x, err := c.Get()
				end := time.Now().Sub(begin)

				if err != nil {
					counterGetFail++
					t.Fatal(err, end.Seconds(), "s")
					return
				}

				time.Sleep(time.Millisecond * 100)
				c.Put(x)
				counterSuccess++
			}
		}(n)
	}

	time.Sleep(time.Second * 10)
	n := math.Trunc((float64(counterSuccess)/float64(threadNum))*1e2+0.5) * 1e-2
	t.Log("Len", c.len(), "MaxUD", m, "Success", counterSuccess, "Failed", counterGetFail, "Sum", n*100, "%", c.counter)
}
