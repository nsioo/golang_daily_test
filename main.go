package main

import "do_some_fxxking_test/leet_code"

// 测试用例:[2,7,11,15] 9 测试结果:[0,0] 期望结果:[0,1] stdout:

// 新建一个 List 用来存每一位上之和的结果
// 一个数字表示进位
// 遍历长的 List
// 解答失败: 测试用例:[9,9,9,9,9,9,9] [9,9,9,9] 测试结果:[8,9,9,9,0,0,0] 期望结果:[8,9,9,9,0,0,0,1] stdout:
func main()  {
	l1 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	l2 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	l3 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	l4 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	l5 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	l6 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	l7 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	l1.Next = &l2
	l2.Next = &l3
	l3.Next = &l4
	l4.Next = &l5
	l5.Next = &l6
	l6.Next = &l7

	p1 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	p2 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	p3 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	p4 := leet_code.ListNode{
		Val:  9,
		Next: nil,
	}
	p1.Next = &p2
	p2.Next = &p3
	p3.Next = &p4

	leet_code.AddTwoNumbers(&l1, &p1)
}
