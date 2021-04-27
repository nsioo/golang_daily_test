package main

import (
	"container/list"
	"fmt"
	"unsafe"
)

func main() {
	// go 双向链表的基本运用
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	for i := l.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value)
	}


	// 测试电脑一个 整型 占多少个字节
	var a int = 1
	fmt.Println(unsafe.Sizeof(a))
}
