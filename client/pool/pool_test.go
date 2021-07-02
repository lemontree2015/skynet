package pool

import (
	"testing"
	"time"
)

type testResource struct {
	isClose bool
}

func (tr *testResource) Close() {
	tr.isClose = true
}

func (tr *testResource) IsClosed() bool {
	return tr.isClose
}

func TestResourcePoolBasic(t *testing.T) {
	// 创建一个resource pool
	pool := NewResourcePool(
		func() (Resource, error) {
			return &testResource{false}, nil
		}, 3, 10)

	if pool.NumResources() != 0 {
		t.Fatal("pool.NumResources() Error")
	}

	// 申请两个资源r1, r2
	r1, err := pool.Acquire()
	if err != nil {
		t.Fatal("pool.Acquire() Error")
	}
	r2, err := pool.Acquire()
	if err != nil {
		t.Fatal("pool.Acquire() Error")
	}

	if pool.NumResources() != 2 {
		t.Fatal("pool.NumResources() Error")
	}

	// 释放r2, r1
	pool.Release(r2)
	pool.Release(r1)

	time.Sleep(200 * time.Microsecond) // 因为relese是async的(备注: 不是一个好的testcase)

	// 申请r3(r2),r4(r1)
	r3, err := pool.Acquire()
	if err != nil {
		t.Fatal("pool.Acquire() Error")
	}

	if pool.NumResources() != 2 {
		t.Fatal("pool.NumResources() Error")
	}

	r4, err := pool.Acquire()
	if err != nil {
		t.Fatal("pool.Acquire() Error")
	}

	if pool.NumResources() != 2 {
		t.Fatal("pool.NumResources() Error")
	}

	if r3 != r2 {
		t.Fatal("pool Logic Error")
	}

	if r4 != r1 {
		t.Fatal("pool Logic Error")
	}

	// 申请r5, r6
	r5, err := pool.Acquire()
	if err != nil {
		t.Fatal("pool.Acquire() Error")
	}

	if pool.NumResources() != 3 {
		t.Fatal("pool.NumResources() Error")
	}

	r6, err := pool.Acquire()
	if err != nil {
		t.Fatal("pool.Acquire() Error")
	}

	if pool.NumResources() != 4 {
		t.Fatal("pool.NumResources() Error")
	}

	pool.Release(r3)
	pool.Release(r4)
	pool.Release(r5)
	pool.Release(r6)

	time.Sleep(200 * time.Microsecond) // 因为relese是async的(备注: 不是一个好的testcase)

	// 因为idleCapacity = 3
	if pool.NumResources() != 3 {
		t.Fatal("pool.NumResources() Error")
	}
}
