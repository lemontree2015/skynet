package pool

import (
	"testing"
)

type testRing struct {
	isClose bool
}

func (tr *testRing) Close() {
	tr.isClose = true
}

func (tr *testRing) IsClosed() bool {
	return tr.isClose
}

func TestRingBasic(t *testing.T) {
	tr1 := &testRing{
		isClose: false,
	}
	tr2 := &testRing{
		isClose: false,
	}
	tr3 := &testRing{
		isClose: false,
	}
	tr4 := &testRing{
		isClose: false,
	}
	tr5 := &testRing{
		isClose: false,
	}

	r := NewRing()
	if r.Empty() != true {
		t.Fatal("r.Empty() Error")
	}

	if r.Size() != 0 {
		t.Fatal("r.Size() Error")
	}

	// Enqueue tr1
	r.Enqueue(tr1)

	if r.Empty() != false {
		t.Fatal("r.Empty() Error")
	}
	if r.Size() != 1 {
		t.Fatal("r.Size() Error")
	}
	if r.Peek() != tr1 {
		t.Fatal("r.Peek() Error")
	}

	// Dequeue tr1
	if r.Dequeue() != tr1 {
		t.Fatal("r.Dequeue() Error")
	}
	if r.Empty() != true {
		t.Fatal("r.Empty() Error")
	}

	if r.Size() != 0 {
		t.Fatal("r.Size() Error")
	}

	// Enqueue tr2, tr3, tr4, tr5
	r.Enqueue(tr2)
	r.Enqueue(tr3)
	r.Enqueue(tr4)
	r.Enqueue(tr5)

	if r.Empty() != false {
		t.Fatal("r.Empty() Error")
	}
	if r.Size() != 4 {
		t.Fatal("r.Size() Error")
	}
	if r.Peek() != tr2 {
		t.Fatal("r.Peek() Error")
	}

	// Dequeue tr2, tr3
	if r.Dequeue() != tr2 {
		t.Fatal("r.Dequeue() Error")
	}
	if r.Dequeue() != tr3 {
		t.Fatal("r.Dequeue() Error")
	}
	if r.Empty() != false {
		t.Fatal("r.Empty() Error")
	}

	if r.Size() != 2 {
		t.Fatal("r.Size() Error")
	}

	// Dequeue tr4, tr5
	if r.Dequeue() != tr4 {
		t.Fatal("r.Dequeue() Error")
	}
	if r.Dequeue() != tr5 {
		t.Fatal("r.Dequeue() Error")
	}
	if r.Empty() != true {
		t.Fatal("r.Empty() Error")
	}

	if r.Size() != 0 {
		t.Fatal("r.Size() Error")
	}
}
