package leet_code

import (
	"fmt"
)

type MyCircularQueue struct {
	head int
	tail int
	data []int
}


func Constructor(k int) MyCircularQueue {
	return MyCircularQueue{
		head: 0,
		tail: 0,
		data: make([]int, k + 1),
	}
}

func (this *MyCircularQueue) EnQueue(value int) bool {
	// 判断是否还有可用空间
	if this.IsFull() {
		return false
	}

	this.tail = (this.tail + 1) % len(this.data)
	this.data[this.tail] = value

	return true
}


func (this *MyCircularQueue) DeQueue() bool {
	if this.IsEmpty() {
		return false
	}

	this.head = (this.head + 1) % len(this.data)
	return true
}


func (this *MyCircularQueue) Front() int {
	if this.IsEmpty() {
		return -1
	}
	return this.data[(this.head + 1) % len(this.data)]
}


func (this *MyCircularQueue) Rear() int {
	if this.IsEmpty() {
		return -1
	}
	return this.data[this.tail]
}


func (this *MyCircularQueue) IsEmpty() bool {
	return this.head == this.tail
}


func (this *MyCircularQueue) IsFull() bool {
	return (this.tail + 1) % len(this.data) == this.head
}


/**
 * Your MyCircularQueue object will be instantiated and called as such:
 * obj := Constructor(k);
 * param_1 := obj.EnQueue(value);
 * param_2 := obj.DeQueue();
 * param_3 := obj.Front();
 * param_4 := obj.Rear();
 * param_5 := obj.IsEmpty();
 * param_6 := obj.IsFull();
 */

func main() {
	queue := Constructor(3)
	fmt.Println(queue.EnQueue(1))
	fmt.Println(queue.EnQueue(2))
	fmt.Println(queue.EnQueue(3))
	fmt.Println(queue.EnQueue(4))

	fmt.Println(queue.Rear())
	fmt.Println(queue.IsFull())
	fmt.Println(queue.DeQueue())
	fmt.Println(queue.EnQueue(4))
	fmt.Println(queue.Rear())
}
