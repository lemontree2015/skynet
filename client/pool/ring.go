package pool

// FIFO规则的Ring

type Ring struct {
	count int // 元素个数
	i     int
	data  []Resource
}

func NewRing() *Ring {
	return &Ring{}
}

func (rb *Ring) Size() int {
	return rb.count
}

func (rb *Ring) Empty() bool {
	return rb.count == 0
}

// 注意:
// 如果没有Enqueue(或者Empty()为true), 直接调用该函数会panic
// 正确的调用方式是先检测Empty(), Peek()
func (rb *Ring) Peek() Resource {
	return rb.data[rb.i]
}

func (rb *Ring) Enqueue(x Resource) {
	if rb.count >= len(rb.data) {
		rb.grow(2*rb.count + 1)
	}
	rb.data[(rb.i+rb.count)%len(rb.data)] = x
	rb.count++
}

// 注意:
// 如果没有Enqueue(或者Empty()为true), 直接调用该函数会panic
// 正确的调用方式是先检测Empty(), 再Dequeue()
func (rb *Ring) Dequeue() (x Resource) {
	x = rb.Peek()
	rb.count, rb.i = rb.count-1, (rb.i+1)%len(rb.data)
	return
}

func (rb *Ring) grow(newSize int) {
	newData := make([]Resource, newSize)

	n := copy(newData, rb.data[rb.i:])
	copy(newData[n:], rb.data[:rb.count-n])

	rb.i = 0
	rb.data = newData
}
